package ticktick

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const tokenDir = ".config/ticky"
const tokenFile = "token.json"

// TokenPath returns the full path to the token file.
func TokenPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, tokenDir, tokenFile)
}

// SaveToken writes the token to disk with 0600 permissions.
func SaveToken(token *OAuthToken) error {
	path := TokenPath()
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("failed to create token directory: %w", err)
	}

	data, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal token: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write token file: %w", err)
	}

	return nil
}

// LoadToken reads the token from disk.
func LoadToken() (*OAuthToken, error) {
	path := TokenPath()
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read token file: %w", err)
	}

	var token OAuthToken
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, fmt.Errorf("failed to parse token file: %w", err)
	}

	return &token, nil
}

// DeleteToken removes the token file and cached config.
func DeleteToken() error {
	path := TokenPath()
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete token file: %w", err)
	}
	// Also remove cached config
	configPath := configDir()
	_ = os.Remove(filepath.Join(configPath, "config.json"))
	return nil
}

func configDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, tokenDir)
}

// SaveInboxID caches the inbox project ID.
func SaveInboxID(id string) error {
	dir := configDir()
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	data, _ := json.Marshal(map[string]string{"inbox_id": id})
	return os.WriteFile(filepath.Join(dir, "config.json"), data, 0600)
}

// LoadInboxID reads the cached inbox project ID.
func LoadInboxID() (string, error) {
	data, err := os.ReadFile(filepath.Join(configDir(), "config.json"))
	if err != nil {
		return "", err
	}
	var cfg map[string]string
	if err := json.Unmarshal(data, &cfg); err != nil {
		return "", err
	}
	return cfg["inbox_id"], nil
}
