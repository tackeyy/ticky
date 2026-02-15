package ticktick

import "net/http"

// NewTestClient creates a Client with a custom httpClient for testing.
func NewTestClient(httpClient *http.Client, accessToken string) *Client {
	return &Client{
		httpClient:  httpClient,
		accessToken: accessToken,
	}
}

// SetBaseURL overrides the package-level baseURL for testing.
// Returns a restore function that resets the original value.
func SetBaseURL(url string) func() {
	orig := baseURL
	baseURL = url
	return func() { baseURL = orig }
}

// SetTokenURL overrides the package-level tokenURL for testing.
// Returns a restore function that resets the original value.
func SetTokenURL(url string) func() {
	orig := tokenURL
	tokenURL = url
	return func() { tokenURL = orig }
}

// ExchangeToken exposes exchangeToken for testing.
func ExchangeToken(clientID, clientSecret, code string) (*OAuthToken, error) {
	return exchangeToken(clientID, clientSecret, code)
}

// GenerateState exposes generateState for testing.
func GenerateState() (string, error) {
	return generateState()
}
