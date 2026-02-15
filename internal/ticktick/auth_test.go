package ticktick_test

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/tackeyy/ticky/internal/ticktick"
)

// setupAuthEnv sets TICKTICK_CLIENT_ID and TICKTICK_CLIENT_SECRET for testing.
func setupAuthEnv(t *testing.T) {
	t.Helper()
	t.Setenv("TICKTICK_CLIENT_ID", "test-client-id")
	t.Setenv("TICKTICK_CLIENT_SECRET", "test-client-secret")
}

// --- exchangeToken ---

func TestExchangeToken_Success(t *testing.T) {
	// Arrange
	wantToken := ticktick.OAuthToken{
		AccessToken:  "access-token-123",
		TokenType:    "bearer",
		ExpiresIn:    3600,
		Scope:        "tasks:read tasks:write",
		RefreshToken: "refresh-token-456",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify method
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}

		// Verify Content-Type
		if ct := r.Header.Get("Content-Type"); ct != "application/x-www-form-urlencoded" {
			t.Errorf("Content-Type = %q, want %q", ct, "application/x-www-form-urlencoded")
		}

		// Verify Basic Auth
		username, password, ok := r.BasicAuth()
		if !ok {
			t.Error("Basic Auth not set")
		}
		if username != "test-client-id" {
			t.Errorf("Basic Auth username = %q, want %q", username, "test-client-id")
		}
		if password != "test-client-secret" {
			t.Errorf("Basic Auth password = %q, want %q", password, "test-client-secret")
		}

		// Verify form body
		body, _ := io.ReadAll(r.Body)
		bodyStr := string(body)
		if !strings.Contains(bodyStr, "code=test-code") {
			t.Errorf("body missing code=test-code, got %q", bodyStr)
		}
		if !strings.Contains(bodyStr, "grant_type=authorization_code") {
			t.Errorf("body missing grant_type=authorization_code, got %q", bodyStr)
		}
		if !strings.Contains(bodyStr, "redirect_uri=") {
			t.Errorf("body missing redirect_uri, got %q", bodyStr)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(wantToken)
	}))
	defer server.Close()

	restore := ticktick.SetTokenURL(server.URL)
	defer restore()

	// Act
	got, err := ticktick.ExchangeToken("test-client-id", "test-client-secret", "test-code")

	// Assert
	if err != nil {
		t.Fatalf("ExchangeToken() returned unexpected error: %v", err)
	}
	if got.AccessToken != wantToken.AccessToken {
		t.Errorf("AccessToken = %q, want %q", got.AccessToken, wantToken.AccessToken)
	}
	if got.TokenType != wantToken.TokenType {
		t.Errorf("TokenType = %q, want %q", got.TokenType, wantToken.TokenType)
	}
	if got.RefreshToken != wantToken.RefreshToken {
		t.Errorf("RefreshToken = %q, want %q", got.RefreshToken, wantToken.RefreshToken)
	}
	if got.Scope != wantToken.Scope {
		t.Errorf("Scope = %q, want %q", got.Scope, wantToken.Scope)
	}
}

func TestExchangeToken_HTTPError(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("invalid credentials"))
	}))
	defer server.Close()

	restore := ticktick.SetTokenURL(server.URL)
	defer restore()

	// Act
	_, err := ticktick.ExchangeToken("bad-id", "bad-secret", "test-code")

	// Assert
	if err == nil {
		t.Fatal("ExchangeToken() expected error for 401 response, got nil")
	}
	if !strings.Contains(err.Error(), "status 401") {
		t.Errorf("error = %q, want to contain 'status 401'", err.Error())
	}
}

func TestExchangeToken_InvalidJSON(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{invalid json"))
	}))
	defer server.Close()

	restore := ticktick.SetTokenURL(server.URL)
	defer restore()

	// Act
	_, err := ticktick.ExchangeToken("test-client-id", "test-client-secret", "test-code")

	// Assert
	if err == nil {
		t.Fatal("ExchangeToken() expected error for invalid JSON, got nil")
	}
	if !strings.Contains(err.Error(), "parse") {
		t.Errorf("error = %q, want to contain 'parse'", err.Error())
	}
}

func TestExchangeToken_ExpiresAtCalculated(t *testing.T) {
	// Arrange
	expiresIn := 7200
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := ticktick.OAuthToken{
			AccessToken: "access-token",
			ExpiresIn:   expiresIn,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	restore := ticktick.SetTokenURL(server.URL)
	defer restore()

	before := time.Now().Unix()

	// Act
	got, err := ticktick.ExchangeToken("test-client-id", "test-client-secret", "test-code")

	// Assert
	after := time.Now().Unix()
	if err != nil {
		t.Fatalf("ExchangeToken() returned unexpected error: %v", err)
	}
	if got.ExpiresAt == 0 {
		t.Fatal("ExpiresAt should be set when ExpiresIn > 0, got 0")
	}
	expectedMin := before + int64(expiresIn)
	expectedMax := after + int64(expiresIn)
	if got.ExpiresAt < expectedMin || got.ExpiresAt > expectedMax {
		t.Errorf("ExpiresAt = %d, want between %d and %d", got.ExpiresAt, expectedMin, expectedMax)
	}
}

// --- RefreshAccessToken ---

func TestRefreshAccessToken_Success(t *testing.T) {
	// Arrange
	setupAuthEnv(t)

	wantToken := ticktick.OAuthToken{
		AccessToken:  "new-access-token",
		TokenType:    "bearer",
		ExpiresIn:    3600,
		Scope:        "tasks:read tasks:write",
		RefreshToken: "new-refresh-token",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify method
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}

		// Verify Content-Type
		if ct := r.Header.Get("Content-Type"); ct != "application/x-www-form-urlencoded" {
			t.Errorf("Content-Type = %q, want %q", ct, "application/x-www-form-urlencoded")
		}

		// Verify Basic Auth
		username, password, ok := r.BasicAuth()
		if !ok {
			t.Error("Basic Auth not set")
		}
		if username != "test-client-id" {
			t.Errorf("Basic Auth username = %q, want %q", username, "test-client-id")
		}
		if password != "test-client-secret" {
			t.Errorf("Basic Auth password = %q, want %q", password, "test-client-secret")
		}

		// Verify form body
		body, _ := io.ReadAll(r.Body)
		bodyStr := string(body)
		if !strings.Contains(bodyStr, "grant_type=refresh_token") {
			t.Errorf("body missing grant_type=refresh_token, got %q", bodyStr)
		}
		if !strings.Contains(bodyStr, "refresh_token=my-refresh-token") {
			t.Errorf("body missing refresh_token=my-refresh-token, got %q", bodyStr)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(wantToken)
	}))
	defer server.Close()

	restore := ticktick.SetTokenURL(server.URL)
	defer restore()

	// Act
	got, err := ticktick.RefreshAccessToken("my-refresh-token")

	// Assert
	if err != nil {
		t.Fatalf("RefreshAccessToken() returned unexpected error: %v", err)
	}
	if got.AccessToken != wantToken.AccessToken {
		t.Errorf("AccessToken = %q, want %q", got.AccessToken, wantToken.AccessToken)
	}
	if got.RefreshToken != wantToken.RefreshToken {
		t.Errorf("RefreshToken = %q, want %q", got.RefreshToken, wantToken.RefreshToken)
	}
}

func TestRefreshAccessToken_MissingEnvVars(t *testing.T) {
	// Arrange: ensure env vars are not set
	os.Unsetenv("TICKTICK_CLIENT_ID")
	os.Unsetenv("TICKTICK_CLIENT_SECRET")

	// Act
	_, err := ticktick.RefreshAccessToken("some-refresh-token")

	// Assert
	if err == nil {
		t.Fatal("RefreshAccessToken() expected error when env vars are missing, got nil")
	}
	if !strings.Contains(err.Error(), "TICKTICK_CLIENT_ID") {
		t.Errorf("error = %q, want to contain 'TICKTICK_CLIENT_ID'", err.Error())
	}
}

func TestRefreshAccessToken_HTTPError(t *testing.T) {
	// Arrange
	setupAuthEnv(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("invalid refresh token"))
	}))
	defer server.Close()

	restore := ticktick.SetTokenURL(server.URL)
	defer restore()

	// Act
	_, err := ticktick.RefreshAccessToken("expired-refresh-token")

	// Assert
	if err == nil {
		t.Fatal("RefreshAccessToken() expected error for 401 response, got nil")
	}
	if !strings.Contains(err.Error(), "status 401") {
		t.Errorf("error = %q, want to contain 'status 401'", err.Error())
	}
}

// --- generateState ---

func TestGenerateState_Uniqueness(t *testing.T) {
	// Act
	state1, err1 := ticktick.GenerateState()
	state2, err2 := ticktick.GenerateState()

	// Assert
	if err1 != nil {
		t.Fatalf("GenerateState() first call returned error: %v", err1)
	}
	if err2 != nil {
		t.Fatalf("GenerateState() second call returned error: %v", err2)
	}
	if state1 == state2 {
		t.Errorf("GenerateState() returned same value twice: %q", state1)
	}
}

func TestGenerateState_Length(t *testing.T) {
	// Act
	state, err := ticktick.GenerateState()

	// Assert
	if err != nil {
		t.Fatalf("GenerateState() returned error: %v", err)
	}
	// 16 bytes -> base64 URL encoding -> 24 characters (with padding)
	// base64.URLEncoding includes padding: ceil(16/3)*4 = 24
	if len(state) != 24 {
		t.Errorf("GenerateState() length = %d, want 24", len(state))
	}
}

func TestGenerateState_Base64(t *testing.T) {
	// Act
	state, err := ticktick.GenerateState()

	// Assert
	if err != nil {
		t.Fatalf("GenerateState() returned error: %v", err)
	}
	decoded, err := base64.URLEncoding.DecodeString(state)
	if err != nil {
		t.Errorf("GenerateState() returned invalid base64 URL string: %v", err)
	}
	if len(decoded) != 16 {
		t.Errorf("decoded bytes length = %d, want 16", len(decoded))
	}
}
