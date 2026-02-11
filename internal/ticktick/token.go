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

// DeleteToken removes the token file.
func DeleteToken() error {
	path := TokenPath()
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete token file: %w", err)
	}
	return nil
}
