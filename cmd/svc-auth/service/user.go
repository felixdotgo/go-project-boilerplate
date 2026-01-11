package service

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/0x46656C6978/go-project-boilerplate/cmd/svc-auth/config"
	"github.com/0x46656C6978/go-project-boilerplate/cmd/svc-auth/entity"
	"github.com/0x46656C6978/go-project-boilerplate/cmd/svc-auth/oauth"
	"github.com/0x46656C6978/go-project-boilerplate/cmd/svc-auth/repository"
	"github.com/0x46656C6978/go-project-boilerplate/cmd/svc-auth/repository/model"
	"github.com/0x46656C6978/go-project-boilerplate/pkg/core"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
)

// TokenPair represents access and refresh tokens
type TokenPair struct {
	AccessToken  string
	RefreshToken string
	TokenType    string
	ExpiresIn    int64
	Scope        string
}

// OAuthUserInfo contains user data from OAuth provider
type OAuthUserInfo struct {
	ProviderUserID string
	Email          string
	EmailVerified  bool
	Name           string
	FirstName      string
	LastName       string
	AvatarURL      string
	Locale         string
	AccessToken    string
	RefreshToken   string
	TokenExpiry    time.Time
}

// UserServiceInterface is an interface define all methods that will be used to handle user
type UserServiceInterface interface {
	Login(ctx context.Context, email, password string) (*TokenPair, error)
	Register(ctx context.Context, email, password string) error

	RefreshAccessToken(ctx context.Context, refreshToken string) (*TokenPair, error)
	RevokeToken(ctx context.Context, refreshToken string) error

	GetOAuthAuthorizationURL(provider, state, redirectURI string) (string, error)
	OAuthCallback(ctx context.Context, provider, code, redirectURI string) (*TokenPair, bool, error)
	LinkOAuthProvider(ctx context.Context, userID uint, provider, code, redirectURI string) error
	UnlinkOAuthProvider(ctx context.Context, userID uint, provider string) error
	GetUserLinkedProviders(ctx context.Context, userID uint) ([]string, error)
}

// UserService is a struct that implements UserServiceInterface
type UserService struct {
	*core.ServiceBase
	r   repository.UserRepoInterface
	cfg *config.Config
}

// NewUserService creates a new UserService
func NewUserService(cfg *config.Config, userRepo repository.UserRepoInterface) UserServiceInterface {
	return &UserService{
		ServiceBase: core.NewService("user"),
		r:           userRepo,
		cfg:         cfg,
	}
}

// Login authenticates user and returns token pair
func (u *UserService) Login(ctx context.Context, email, password string) (*TokenPair, error) {
	user, err := u.r.FindByEmail(ctx, email)
	if err != nil {
		return nil, u.repoToServiceError(err)
	}

	// Verify password
	err = u.VerifyCredentials(ctx, user, email, password)
	if err != nil {
		return nil, err
	}

	return u.generateTokenPair(ctx, user)
}

// Register creates a new user with email and password
func (u *UserService) Register(ctx context.Context, email, password string) error {
	user, err := u.r.FindByEmail(ctx, email)
	if err == nil && user != nil {
		return ErrUserAlreadyExists
	}

	newUser := &entity.User{
		Email: email,
	}
	if err := newUser.SetPassword(password); err != nil {
		return ErrBadRequest
	}

	return u.r.Save(ctx, newUser)
}

// VerifyCredentials verify user credentials
func (u *UserService) VerifyCredentials(ctx context.Context, user *entity.User, email, password string) error {
	// comparing user email
	if user.Email != email {
		return errors.New("invalid email")
	}
	// comparing user password
	if !user.IsValidPassword(password) {
		return errors.New("invalid password")
	}
	return nil
}

// repoToServiceError convert repository error to service error
func (u *UserService) repoToServiceError(err error) error {
	if errors.Is(err, repository.ErrNotFound) {
		return ErrUserNotFound
	}
	return err
}

// GenerateTokenPair creates both access and refresh tokens
func (u *UserService) generateTokenPair(ctx context.Context, user *entity.User) (*TokenPair, error) {
	// Generate JWT access token (you'll need to implement this based on your JWT setup)
	accessToken, err := u.generateJWTToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token
	refreshToken, err := u.generateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Hash and store refresh token
	tokenHash := u.hashToken(refreshToken)
	refreshTokenEntity := &model.RefreshToken{
		UserID:    user.ID,
		Token:     tokenHash,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour), // 30 days
	}

	if err := u.r.CreateRefreshToken(ctx, refreshTokenEntity); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600, // 1 hour
		Scope:        "openid profile email",
	}, nil
}

// RefreshAccessToken validates refresh token and issues new access token
func (u *UserService) RefreshAccessToken(ctx context.Context, refreshToken string) (*TokenPair, error) {
	tokenHash := u.hashToken(refreshToken)

	// Retrieve token from database
	storedToken, err := u.r.FindRefreshTokenByToken(ctx, tokenHash)
	if err != nil {
		return nil, ErrRefreshTokenInvalid
	}

	// Validate token
	if !storedToken.IsValid() {
		return nil, ErrRefreshTokenInvalid
	}

	// Fetch the user
	user, err := u.r.FindByID(ctx, storedToken.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Generate new token pair
	return u.generateTokenPair(ctx, user)
}

// RevokeToken revokes a refresh token
func (u *UserService) RevokeToken(ctx context.Context, refreshToken string) error {
	tokenHash := u.hashToken(refreshToken)

	storedToken, err := u.r.FindRefreshTokenByToken(ctx, tokenHash)
	if err != nil {
		return ErrRefreshTokenInvalid
	}

	return u.r.RevokeRefreshToken(ctx, storedToken.ID)
}

// GetOAuthAuthorizationURL generates the OAuth authorization URL
func (u *UserService) GetOAuthAuthorizationURL(provider, state, redirectURI string) (string, error) {
	if provider == "" {
		return "", ErrOAuthProviderEmpty
	}
	if state == "" {
		return "", ErrOAuthStateEmpty
	}

	oauthConfig := u.getOAuthConfig(provider)
	if oauthConfig == nil {
		return "", ErrOAuthProviderNotConfigured
	}

	oauthProvider, err := oauth.NewProvider(provider, oauthConfig)
	if err != nil {
		return "", ErrOAuthProviderNotConfigured
	}

	return oauthProvider.GetAuthURL(state), nil
}

func (u *UserService) OAuthCallback(ctx context.Context, provider, code, redirectURI string) (*TokenPair, bool, error) {
	if provider == "" {
		return nil, false, ErrOAuthProviderEmpty
	}
	if code == "" {
		return nil, false, ErrOAuthCodeEmpty
	}

	user, isNewUser, err := u.findOrCreateOAuthUser(ctx, provider, code, redirectURI)
	if err != nil {
		return nil, false, err
	}

	tokenPair, err := u.generateTokenPair(ctx, user)
	if err != nil {
		return nil, false, err
	}

	return tokenPair, isNewUser, nil
}

// FindOrCreateOAuthUser handles OAuth user authentication via code exchange
// Returns (user, isNewUser, error)
func (u *UserService) findOrCreateOAuthUser(ctx context.Context, provider, code, redirectURI string) (*entity.User, bool, error) {
	// Get OAuth config for provider
	oauthConfig := u.getOAuthConfig(provider)
	if oauthConfig == nil {
		return nil, false, ErrOAuthProviderNotConfigured
	}

	// Get OAuth provider
	oauthProvider, err := oauth.NewProvider(provider, oauthConfig)
	if err != nil {
		return nil, false, ErrOAuthProviderNotConfigured
	}

	// Exchange code for token
	token, err := oauthProvider.ExchangeCode(ctx, code)
	if err != nil {
		return nil, false, ErrOAuthExchangeFailed
	}

	// Get user info from provider
	userInfo, err := oauthProvider.GetUserInfo(ctx, token)
	if err != nil {
		return nil, false, ErrOAuthUserInfoFailed
	}

	// Try to find existing user by provider and provider user ID
	user, err := u.r.FindUserByProviderAndProviderUserID(ctx, provider, userInfo.ProviderUserID)
	if err == nil {
		// User found - update OAuth provider info
		oauthProviderEntity, err := u.r.FindOAuthProvider(ctx, user.ID, provider)
		if err != nil {
			return nil, false, fmt.Errorf("failed to get oauth provider: %w", err)
		}

		// Update OAuth tokens and profile (encrypted)
		encryptedAccess := u.encrypt(token.AccessToken)
		encryptedRefresh := ""
		if token.RefreshToken != "" {
			encryptedRefresh = u.encrypt(token.RefreshToken)
		}
		oauthProviderEntity.AccessToken = &encryptedAccess
		if encryptedRefresh != "" {
			oauthProviderEntity.RefreshToken = &encryptedRefresh
		}
		oauthProviderEntity.TokenExpiry = &token.Expiry
		oauthProviderEntity.Name = &userInfo.Name
		oauthProviderEntity.AvatarURL = &userInfo.AvatarURL
		oauthProviderEntity.EmailVerified = userInfo.EmailVerified

		if err := u.r.UpdateOAuthProvider(ctx, oauthProviderEntity); err != nil {
			return nil, false, err
		}

		// Update user's aggregated profile if needed
		if user.Name == nil || user.AvatarURL == nil {
			user.Name = &userInfo.Name
			user.AvatarURL = &userInfo.AvatarURL
			user.EmailVerified = userInfo.EmailVerified
			if err := u.r.Update(ctx, user); err != nil {
				return nil, false, err
			}
		}

		return user, false, nil
	}

	// Check if email already exists (potential account linking scenario)
	existingUser, err := u.r.FindByEmail(ctx, userInfo.Email)
	if err == nil {
		// User exists with this email - check if they can link this provider
		// Check if provider is already linked
		_, err := u.r.FindOAuthProvider(ctx, existingUser.ID, provider)
		if err == nil {
			return nil, false, ErrOAuthProviderAlreadyLinked
		}

		// Auto-link if email is verified
		if userInfo.EmailVerified {
			if err := u.linkOAuthProviderInternal(ctx, existingUser.ID, provider, userInfo, token); err != nil {
				return nil, false, fmt.Errorf("failed to auto-link provider: %w", err)
			}
			return existingUser, false, nil
		}

		// Email not verified - return error requiring manual linking
		return nil, false, fmt.Errorf("email already registered - please login and link provider manually")
	}

	// Create new user with OAuth provider
	newUser := &entity.User{
		Email:         userInfo.Email,
		Name:          &userInfo.Name,
		FirstName:     &userInfo.FirstName,
		LastName:      &userInfo.LastName,
		AvatarURL:     &userInfo.AvatarURL,
		EmailVerified: userInfo.EmailVerified,
		Locale:        &userInfo.Locale,
	}

	if err := u.r.Save(ctx, newUser); err != nil {
		return nil, false, err
	}

	// Link the OAuth provider to the new user
	if err := u.linkOAuthProviderInternal(ctx, newUser.ID, provider, userInfo, token); err != nil {
		return nil, false, fmt.Errorf("failed to link oauth provider: %w", err)
	}

	return newUser, true, nil
}

// LinkOAuthProvider links a new OAuth provider to an existing user via code exchange
func (u *UserService) LinkOAuthProvider(ctx context.Context, userID uint, provider, code, redirectURI string) error {
	if provider == "" {
		return ErrOAuthProviderEmpty
	}
	if code == "" {
		return ErrOAuthCodeEmpty
	}

	// Get OAuth config for provider
	oauthConfig := u.getOAuthConfig(provider)
	if oauthConfig == nil {
		return ErrOAuthProviderNotConfigured
	}

	// Get OAuth provider
	oauthProvider, err := oauth.NewProvider(provider, oauthConfig)
	if err != nil {
		return ErrOAuthProviderNotConfigured
	}

	// Exchange code for token
	token, err := oauthProvider.ExchangeCode(ctx, code)
	if err != nil {
		return ErrOAuthExchangeFailed
	}

	// Get user info from provider
	userInfo, err := oauthProvider.GetUserInfo(ctx, token)
	if err != nil {
		return ErrOAuthUserInfoFailed
	}

	// Check if provider is already linked to another user
	existingUser, err := u.r.FindUserByProviderAndProviderUserID(ctx, provider, userInfo.ProviderUserID)
	if err == nil && uint(existingUser.ID) != userID {
		return ErrOAuthProviderAlreadyLinked
	}

	// Check if provider is already linked to this user
	_, err = u.r.FindOAuthProvider(ctx, int(userID), provider)
	if err == nil {
		return ErrOAuthProviderAlreadyLinked
	}

	return u.linkOAuthProviderInternal(ctx, int(userID), provider, userInfo, token)
}

// linkOAuthProviderInternal is an internal helper to link OAuth provider
func (u *UserService) linkOAuthProviderInternal(ctx context.Context, userID int, provider string, userInfo *oauth.UserInfo, token *oauth2.Token) error {
	encryptedAccess := u.encrypt(token.AccessToken)
	var encryptedRefresh *string
	if token.RefreshToken != "" {
		encrypted := u.encrypt(token.RefreshToken)
		encryptedRefresh = &encrypted
	}

	var tokenExpiry *time.Time
	if !token.Expiry.IsZero() {
		tokenExpiry = &token.Expiry
	}

	return u.r.CreateOAuthProvider(ctx, userID, provider, &encryptedAccess, encryptedRefresh, userInfo, tokenExpiry)
}

// UnlinkOAuthProvider removes an OAuth provider from a user account
func (u *UserService) UnlinkOAuthProvider(ctx context.Context, userID uint, provider string) error {
	if provider == "" {
		return ErrOAuthProviderEmpty
	}

	// Get user to check authentication methods
	user, err := u.r.FindByID(ctx, int(userID))
	if err != nil {
		return ErrUserNotFound
	}

	// Check if provider is linked
	_, err = u.r.FindOAuthProvider(ctx, int(userID), provider)
	if err != nil {
		return ErrOAuthProviderNotLinked
	}

	// Count linked providers
	providers, err := u.r.FindUserOAuthProviders(ctx, int(userID))
	if err != nil {
		return err
	}

	// Check if this is the last authentication method
	hasPassword := user.HasPassword()
	if !hasPassword && len(providers) <= 1 {
		return ErrCannotUnlinkLastAuthMethod
	}

	return u.r.DeleteOAuthProvider(ctx, int(userID), provider)
}

// GetUserLinkedProviders returns list of linked OAuth providers
func (u *UserService) GetUserLinkedProviders(ctx context.Context, userID uint) ([]string, error) {
	providers, err := u.r.FindUserOAuthProviders(ctx, int(userID))
	if err != nil {
		return nil, err
	}

	providerNames := make([]string, len(providers))
	for i, p := range providers {
		providerNames[i] = p.Provider
	}
	return providerNames, nil
}

// Helper methods

// getOAuthConfig converts app config to OAuth config for a specific provider
func (u *UserService) getOAuthConfig(provider string) *oauth.Config {
	switch provider {
	case "google":
		if !u.cfg.OAuth.Google.IsEnabled() {
			return nil
		}
		redirectURL := u.cfg.OAuth.Google.RedirectURL
		if redirectURL == "" {
			redirectURL = u.cfg.OAuth.RedirectURL
		}
		return &oauth.Config{
			ClientID:     u.cfg.OAuth.Google.ClientID,
			ClientSecret: u.cfg.OAuth.Google.ClientSecret,
			RedirectURL:  redirectURL,
			Scopes:       u.cfg.OAuth.Google.GetScopes(),
		}
	case "facebook":
		if !u.cfg.OAuth.Facebook.IsEnabled() {
			return nil
		}
		redirectURL := u.cfg.OAuth.Facebook.RedirectURL
		if redirectURL == "" {
			redirectURL = u.cfg.OAuth.RedirectURL
		}
		return &oauth.Config{
			ClientID:     u.cfg.OAuth.Facebook.ClientID,
			ClientSecret: u.cfg.OAuth.Facebook.ClientSecret,
			RedirectURL:  redirectURL,
			Scopes:       u.cfg.OAuth.Facebook.GetScopes(),
		}
	case "github":
		if !u.cfg.OAuth.GitHub.IsEnabled() {
			return nil
		}
		redirectURL := u.cfg.OAuth.GitHub.RedirectURL
		if redirectURL == "" {
			redirectURL = u.cfg.OAuth.RedirectURL
		}
		return &oauth.Config{
			ClientID:     u.cfg.OAuth.GitHub.ClientID,
			ClientSecret: u.cfg.OAuth.GitHub.ClientSecret,
			RedirectURL:  redirectURL,
			Scopes:       u.cfg.OAuth.GitHub.GetScopes(),
		}
	default:
		return nil
	}
}

func (u *UserService) generateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (u *UserService) hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func (u *UserService) encrypt(plaintext string) string {
	if plaintext == "" {
		return ""
	}

	key := []byte(u.cfg.JWT.EncryptionKey)
	if len(key) != 32 {
		// Encryption key must be 32 bytes for AES-256
		// In production, this should fail gracefully or use a default secure key
		return plaintext
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return plaintext
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return plaintext
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return plaintext
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext)
}

// generateJWTToken generates a JWT access token for the user
func (u *UserService) generateJWTToken(user *entity.User) (string, error) {
	now := time.Now()

	// Parse access token TTL from config (e.g., "1h")
	ttlDuration, err := time.ParseDuration(u.cfg.JWT.AccessTokenTTL)
	if err != nil {
		// Fallback to 1 hour if parsing fails
		ttlDuration = 1 * time.Hour
	}

	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(now.Add(ttlDuration)),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		Issuer:    u.cfg.JWT.Issuer,
		Subject:   strconv.Itoa(user.ID),
		ID:        strconv.Itoa(user.ID),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(u.cfg.JWT.Secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT token: %w", err)
	}

	return signedToken, nil
}
