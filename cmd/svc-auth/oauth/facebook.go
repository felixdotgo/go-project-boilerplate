package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

// FacebookProvider implements Provider for Facebook OAuth
type FacebookProvider struct {
	config *oauth2.Config
}

// NewFacebookProvider creates a new Facebook OAuth provider
func NewFacebookProvider(cfg *Config) *FacebookProvider {
	scopes := cfg.Scopes
	if len(scopes) == 0 {
		scopes = []string{"email", "public_profile"}
	}

	return &FacebookProvider{
		config: &oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURL:  cfg.RedirectURL,
			Scopes:       scopes,
			Endpoint:     facebook.Endpoint,
		},
	}
}

// GetAuthURL returns authorization URL for user redirection
func (f *FacebookProvider) GetAuthURL(state string) string {
	return f.config.AuthCodeURL(state)
}

// ExchangeCode exchanges authorization code for tokens
func (f *FacebookProvider) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	return f.config.Exchange(ctx, code)
}

// GetUserInfo retrieves user information using access token
func (f *FacebookProvider) GetUserInfo(ctx context.Context, token *oauth2.Token) (*UserInfo, error) {
	client := f.config.Client(ctx, token)
	url := "https://graph.facebook.com/me?fields=id,email,name,first_name,last_name,picture,locale"
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("facebook API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var fbUser struct {
		ID        string `json:"id"`
		Email     string `json:"email"`
		Name      string `json:"name"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Picture   struct {
			Data struct {
				URL string `json:"url"`
			} `json:"data"`
		} `json:"picture"`
		Locale string `json:"locale"`
	}

	if err := json.Unmarshal(body, &fbUser); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %w", err)
	}

	return &UserInfo{
		ProviderUserID: fbUser.ID,
		Email:          fbUser.Email,
		EmailVerified:  true, // Facebook verifies emails
		Name:           fbUser.Name,
		FirstName:      fbUser.FirstName,
		LastName:       fbUser.LastName,
		AvatarURL:      fbUser.Picture.Data.URL,
		Locale:         fbUser.Locale,
	}, nil
}

// RefreshToken refreshes an expired access token
func (f *FacebookProvider) RefreshToken(ctx context.Context, refreshToken string) (*oauth2.Token, error) {
	// Facebook uses long-lived tokens instead of refresh tokens
	// This would need to be implemented using token exchange endpoint
	return nil, fmt.Errorf("facebook token refresh not yet implemented")
}
