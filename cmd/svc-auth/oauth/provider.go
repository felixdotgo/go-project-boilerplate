package oauth

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
)

// Provider defines OAuth provider interface
type Provider interface {
	// GetAuthURL returns authorization URL for user redirection
	GetAuthURL(state string) string

	// ExchangeCode exchanges authorization code for tokens
	ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error)

	// GetUserInfo retrieves user information using access token
	GetUserInfo(ctx context.Context, token *oauth2.Token) (*UserInfo, error)

	// RefreshToken refreshes an expired access token
	RefreshToken(ctx context.Context, refreshToken string) (*oauth2.Token, error)
}

// UserInfo standardized user information from provider
type UserInfo struct {
	ProviderUserID string // Unique ID from provider
	Email          string
	EmailVerified  bool
	Name           string
	FirstName      string
	LastName       string
	AvatarURL      string
	Locale         string
}

// Config holds OAuth provider configuration
type Config struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
}

// NewProvider creates a provider instance
func NewProvider(providerName string, config *Config) (Provider, error) {
	switch providerName {
	case "google":
		return NewGoogleProvider(config), nil
	case "facebook":
		return NewFacebookProvider(config), nil
	case "github":
		return NewGitHubProvider(config), nil
	default:
		return nil, fmt.Errorf("unsupported provider: %s", providerName)
	}
}
