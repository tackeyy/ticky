package ticktick

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)


// setupTestHome overrides the HOME environment variable to isolate
// file system operations within t.TempDir().
func setupTestHome(t *testing.T) {
	t.Helper()
	t.Setenv("HOME", t.TempDir())
}

func TestSaveAndLoadToken_RoundTrip(t *testing.T) {
	setupTestHome(t)

	// Arrange
	want := &OAuthToken{
		AccessToken:  "test-access-token",
		TokenType:    "bearer",
		ExpiresIn:    3600,
		Scope:        "tasks:read tasks:write",
		RefreshToken: "test-refresh-token",
		ExpiresAt:    1700000000,
	}

	// Act
	if err := SaveToken(want); err != nil {
		t.Fatalf("SaveToken() returned unexpected error: %v", err)
	}
	got, err := LoadToken()
	if err != nil {
		t.Fatalf("LoadToken() returned unexpected error: %v", err)
	}

	// Assert
	if got.AccessToken != want.AccessToken {
		t.Errorf("AccessToken = %q, want %q", got.AccessToken, want.AccessToken)
	}
	if got.TokenType != want.TokenType {
		t.Errorf("TokenType = %q, want %q", got.TokenType, want.TokenType)
	}
	if got.ExpiresIn != want.ExpiresIn {
		t.Errorf("ExpiresIn = %d, want %d", got.ExpiresIn, want.ExpiresIn)
	}
	if got.Scope != want.Scope {
		t.Errorf("Scope = %q, want %q", got.Scope, want.Scope)
	}
	if got.RefreshToken != want.RefreshToken {
		t.Errorf("RefreshToken = %q, want %q", got.RefreshToken, want.RefreshToken)
	}
	if got.ExpiresAt != want.ExpiresAt {
		t.Errorf("ExpiresAt = %d, want %d", got.ExpiresAt, want.ExpiresAt)
	}
}

func TestSaveToken_CreatesDirectory(t *testing.T) {
	setupTestHome(t)

	// Arrange — directory does not exist yet
	token := &OAuthToken{AccessToken: "dir-test"}

	// Act
	err := SaveToken(token)

	// Assert
	if err != nil {
		t.Fatalf("SaveToken() returned unexpected error: %v", err)
	}
	path := TokenPath()
	if _, statErr := os.Stat(path); os.IsNotExist(statErr) {
		t.Errorf("token file was not created at %s", path)
	}
}

func TestSaveToken_FilePermissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("file permissions not applicable on Windows")
	}

	setupTestHome(t)

	// Arrange
	token := &OAuthToken{AccessToken: "perm-test"}

	// Act
	if err := SaveToken(token); err != nil {
		t.Fatalf("SaveToken() returned unexpected error: %v", err)
	}

	// Assert
	info, err := os.Stat(TokenPath())
	if err != nil {
		t.Fatalf("os.Stat() returned unexpected error: %v", err)
	}
	perm := info.Mode().Perm()
	if perm != 0600 {
		t.Errorf("file permission = %o, want 0600", perm)
	}
}

func TestLoadToken_FileNotFound(t *testing.T) {
	setupTestHome(t)

	// Act
	_, err := LoadToken()

	// Assert
	if err == nil {
		t.Fatal("LoadToken() expected error for missing file, got nil")
	}
}

func TestLoadToken_InvalidJSON(t *testing.T) {
	setupTestHome(t)

	// Arrange — write invalid JSON to the token path
	path := TokenPath()
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}
	if err := os.WriteFile(path, []byte("{invalid json}"), 0600); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	// Act
	_, err := LoadToken()

	// Assert
	if err == nil {
		t.Fatal("LoadToken() expected error for invalid JSON, got nil")
	}
}

func TestDeleteToken_Success(t *testing.T) {
	setupTestHome(t)

	// Arrange — save a token first
	token := &OAuthToken{AccessToken: "delete-test"}
	if err := SaveToken(token); err != nil {
		t.Fatalf("SaveToken() returned unexpected error: %v", err)
	}

	// Act
	if err := DeleteToken(); err != nil {
		t.Fatalf("DeleteToken() returned unexpected error: %v", err)
	}

	// Assert — LoadToken should now fail
	_, err := LoadToken()
	if err == nil {
		t.Fatal("LoadToken() expected error after DeleteToken, got nil")
	}
}

func TestDeleteToken_NonExistent(t *testing.T) {
	setupTestHome(t)

	// Act — delete when no token file exists
	err := DeleteToken()

	// Assert
	if err != nil {
		t.Fatalf("DeleteToken() expected no error for non-existent file, got: %v", err)
	}
}

func TestSaveAndLoadInboxID_RoundTrip(t *testing.T) {
	setupTestHome(t)

	// Arrange
	want := "inbox-project-id-123"

	// Act
	if err := SaveInboxID(want); err != nil {
		t.Fatalf("SaveInboxID() returned unexpected error: %v", err)
	}
	got, err := LoadInboxID()
	if err != nil {
		t.Fatalf("LoadInboxID() returned unexpected error: %v", err)
	}

	// Assert
	if got != want {
		t.Errorf("LoadInboxID() = %q, want %q", got, want)
	}
}

func TestLoadInboxID_FileNotFound(t *testing.T) {
	setupTestHome(t)

	// Act
	_, err := LoadInboxID()

	// Assert
	if err == nil {
		t.Fatal("LoadInboxID() expected error for missing file, got nil")
	}
}

func TestDeleteToken_AlsoRemovesConfig(t *testing.T) {
	setupTestHome(t)

	// Arrange — save both token and inbox config
	token := &OAuthToken{AccessToken: "config-cleanup-test"}
	if err := SaveToken(token); err != nil {
		t.Fatalf("SaveToken() returned unexpected error: %v", err)
	}
	if err := SaveInboxID("inbox-to-delete"); err != nil {
		t.Fatalf("SaveInboxID() returned unexpected error: %v", err)
	}

	// Verify config.json exists before deletion
	configPath := filepath.Join(configDir(), "config.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("config.json should exist before DeleteToken")
	}

	// Act
	if err := DeleteToken(); err != nil {
		t.Fatalf("DeleteToken() returned unexpected error: %v", err)
	}

	// Assert — config.json should also be removed
	if _, err := os.Stat(configPath); !os.IsNotExist(err) {
		t.Errorf("config.json should be removed after DeleteToken, but still exists")
	}
}
