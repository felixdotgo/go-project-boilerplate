package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

// GitHubProvider implements Provider for GitHub OAuth
type GitHubProvider struct {
	config *oauth2.Config
}

// NewGitHubProvider creates a new GitHub OAuth provider
func NewGitHubProvider(cfg *Config) *GitHubProvider {
	scopes := cfg.Scopes
	if len(scopes) == 0 {
		scopes = []string{"user:email", "read:user"}
	}

	return &GitHubProvider{
		config: &oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURL:  cfg.RedirectURL,
			Scopes:       scopes,
			Endpoint:     github.Endpoint,
		},
	}
}

// GetAuthURL returns authorization URL for user redirection
func (gh *GitHubProvider) GetAuthURL(state string) string {
	return gh.config.AuthCodeURL(state)
}

// ExchangeCode exchanges authorization code for tokens
func (gh *GitHubProvider) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	return gh.config.Exchange(ctx, code)
}

// GetUserInfo retrieves user information using access token
func (gh *GitHubProvider) GetUserInfo(ctx context.Context, token *oauth2.Token) (*UserInfo, error) {
	client := gh.config.Client(ctx, token)

	// Get user profile
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("github API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var ghUser struct {
		ID        int64  `json:"id"`
		Login     string `json:"login"`
		Email     string `json:"email"`
		Name      string `json:"name"`
		AvatarURL string `json:"avatar_url"`
		Location  string `json:"location"`
	}

	if err := json.Unmarshal(body, &ghUser); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %w", err)
	}

	// If email is not in profile, fetch from emails endpoint
	email := ghUser.Email
	emailVerified := false
	if email == "" {
		var err error
		email, emailVerified, err = gh.getPrimaryEmail(ctx, client)
		if err != nil {
			return nil, fmt.Errorf("failed to get email: %w", err)
		}
	} else {
		emailVerified = true // Email in profile is verified
	}

	// Parse name into first/last
	firstName, lastName := parseName(ghUser.Name)

	return &UserInfo{
		ProviderUserID: fmt.Sprintf("%d", ghUser.ID),
		Email:          email,
		EmailVerified:  emailVerified,
		Name:           ghUser.Name,
		FirstName:      firstName,
		LastName:       lastName,
		AvatarURL:      ghUser.AvatarURL,
		Locale:         "", // GitHub doesn't provide locale
	}, nil
}

// getPrimaryEmail retrieves primary verified email from GitHub
func (gh *GitHubProvider) getPrimaryEmail(ctx context.Context, client *http.Client) (string, bool, error) {
	resp, err := client.Get("https://api.github.com/user/emails")
	if err != nil {
		return "", false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", false, fmt.Errorf("failed to get emails, status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", false, err
	}

	var emails []struct {
		Email    string `json:"email"`
		Verified bool   `json:"verified"`
		Primary  bool   `json:"primary"`
	}

	if err := json.Unmarshal(body, &emails); err != nil {
		return "", false, err
	}

	// Find primary verified email
	for _, e := range emails {
		if e.Primary && e.Verified {
			return e.Email, true, nil
		}
	}

	// Fallback to any verified email
	for _, e := range emails {
		if e.Verified {
			return e.Email, true, nil
		}
	}

	// Last resort: any email
	if len(emails) > 0 {
		return emails[0].Email, emails[0].Verified, nil
	}

	return "", false, fmt.Errorf("no email found")
}

// RefreshToken refreshes an expired access token
func (gh *GitHubProvider) RefreshToken(ctx context.Context, refreshToken string) (*oauth2.Token, error) {
	// GitHub tokens don't expire by default
	return nil, fmt.Errorf("github tokens don't support refresh")
}

// parseName splits a full name into first and last name
func parseName(fullName string) (string, string) {
	if fullName == "" {
		return "", ""
	}

	// Simple split on first space
	for i, r := range fullName {
		if r == ' ' {
			return fullName[:i], fullName[i+1:]
		}
	}

	// No space found, return full name as first name
	return fullName, ""
}
