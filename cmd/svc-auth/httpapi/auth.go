package httpapi

import (
	"context"
	"errors"
	"net/http"

	"github.com/0x46656C6978/go-project-boilerplate/cmd/svc-auth/config"
	"github.com/0x46656C6978/go-project-boilerplate/cmd/svc-auth/service"
	v1 "github.com/0x46656C6978/go-project-boilerplate/pkg/rpc/api/auth/v1"
)

const (
	ErrInternalServerError = "internal server error"
)

// AuthHttpApi is a struct that implements the AuthServiceServer interface
// contains all methods that will be used to handle authentication
type AuthHttpApi struct {
	v1.UnimplementedAuthServiceServer
	s   service.UserServiceInterface
	cfg *config.Config
}

// NewAuthServiceServer returns a new instance of AuthHttpApi struct that implements the AuthServiceServer interface
func NewAuthServiceServer(cfg *config.Config, s service.UserServiceInterface) v1.AuthServiceServer {
	return &AuthHttpApi{
		s:   s,
		cfg: cfg,
	}
}

// Ping is a method that handles the ping request
func (u *AuthHttpApi) Ping(ctx context.Context, req *v1.Auth_PingRequest) (*v1.Auth_PingResponse, error) {
	return &v1.Auth_PingResponse{
		Message: "pong",
	}, nil
}

// Login is a method that handles the login request
// Returns OAuth2-compliant TokenResponse with access_token and refresh_token
func (u *AuthHttpApi) Login(ctx context.Context, req *v1.Auth_LoginRequest) (*v1.Auth_LoginResponse, error) {
	tokenPair, err := u.s.Login(ctx, req.GetData().GetEmail(), req.GetData().GetPassword())
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return nil, NewError(http.StatusUnauthorized, "invalid email or password")
		}
		return nil, NewError(http.StatusInternalServerError, ErrInternalServerError)
	}

	return &v1.Auth_LoginResponse{
		Data: &v1.Auth_TokenResponse{
			AccessToken:  tokenPair.AccessToken,
			RefreshToken: tokenPair.RefreshToken,
			TokenType:    tokenPair.TokenType,
			ExpiresIn:    tokenPair.ExpiresIn,
			Scope:        tokenPair.Scope,
		},
	}, nil
}

// Regiter is a method that handles the register request
func (u *AuthHttpApi) Regiter(ctx context.Context, req *v1.Auth_RegisterRequest) (*v1.Auth_RegisterResponse, error) {
	err := u.s.Register(ctx, req.GetData().GetEmail(), req.GetData().GetPassword())
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			return nil, NewError(http.StatusConflict, "user with this email already exists")
		}
		if errors.Is(err, service.ErrBadRequest) {
			return nil, NewError(http.StatusBadRequest, "invalid request data")
		}
		return nil, NewError(http.StatusInternalServerError, ErrInternalServerError)
	}

	return &v1.Auth_RegisterResponse{}, nil
}

// RefreshToken is a method that handles the refresh token request
func (u *AuthHttpApi) RefreshToken(ctx context.Context, req *v1.Auth_RefreshTokenRequest) (*v1.Auth_RefreshTokenResponse, error) {
	tokenPair, err := u.s.RefreshAccessToken(ctx, req.GetRefreshToken())
	if err != nil {
		if errors.Is(err, service.ErrRefreshTokenNotFound) ||
			errors.Is(err, service.ErrRefreshTokenExpired) ||
			errors.Is(err, service.ErrRefreshTokenRevoked) {
			return nil, NewError(http.StatusUnauthorized, "invalid or expired refresh token")
		}
		return nil, NewError(http.StatusInternalServerError, ErrInternalServerError)
	}

	return &v1.Auth_RefreshTokenResponse{
		Data: &v1.Auth_TokenResponse{
			AccessToken:  tokenPair.AccessToken,
			RefreshToken: tokenPair.RefreshToken,
			TokenType:    tokenPair.TokenType,
			ExpiresIn:    tokenPair.ExpiresIn,
			Scope:        tokenPair.Scope,
		},
	}, nil
}

// RevokeToken is a method that handles the revoke token request
func (u *AuthHttpApi) RevokeToken(ctx context.Context, req *v1.Auth_RevokeTokenRequest) (*v1.Auth_RevokeTokenResponse, error) {
	err := u.s.RevokeToken(ctx, req.GetToken())
	if err != nil {
		// OAuth2 spec requires 200 OK even if token is invalid
		// Just log the error but return success
		return &v1.Auth_RevokeTokenResponse{}, nil
	}

	return &v1.Auth_RevokeTokenResponse{}, nil
}

// OAuthLogin is a method that handles the OAuth login request
// returns the authorization URL that the client should redirect to
func (u *AuthHttpApi) OAuthLogin(ctx context.Context, req *v1.Auth_OAuthLoginRequest) (*v1.Auth_OAuthLoginResponse, error) {
	provider := req.GetProvider()
	state := req.GetState()
	redirectUri := req.GetRedirectUri()

	authURL, err := u.s.GetOAuthAuthorizationURL(provider, state, redirectUri)
	if err != nil {
		return nil, FromUserToHTTPError(err)
	}

	return &v1.Auth_OAuthLoginResponse{
		AuthorizationUrl: authURL,
	}, nil
}

// OAuthCallback is a method that handles the OAuth callback request
// exchanges the authorization code for tokens and creates/authenticates the user
func (u *AuthHttpApi) OAuthCallback(ctx context.Context, req *v1.Auth_OAuthCallbackRequest) (*v1.Auth_OAuthCallbackResponse, error) {
	provider := req.GetProvider()
	code := req.GetCode()
	redirectUri := req.GetRedirectUri()

	// Find or create user from OAuth provider
	tokenPair, isNewUser, err := u.s.OAuthCallback(ctx, provider, code, redirectUri)
	if err != nil {
		return nil, FromUserToHTTPError(err)
	}

	return &v1.Auth_OAuthCallbackResponse{
		Data: &v1.Auth_TokenResponse{
			AccessToken:  tokenPair.AccessToken,
			RefreshToken: tokenPair.RefreshToken,
			TokenType:    tokenPair.TokenType,
			ExpiresIn:    tokenPair.ExpiresIn,
			Scope:        tokenPair.Scope,
		},
		IsNewUser: isNewUser,
	}, nil
}

// LinkProvider is a method that handles linking an OAuth provider to an existing user
func (u *AuthHttpApi) LinkProvider(ctx context.Context, req *v1.Auth_LinkProviderRequest) (*v1.Auth_LinkProviderResponse, error) {
	// Extract user ID from context (set by auth middleware)
	userID, err := u.getUserIDFromContext(ctx)
	if err != nil {
		return nil, NewError(http.StatusUnauthorized, "unauthorized")
	}

	provider := req.GetProvider()
	code := req.GetCode()
	redirectUri := req.GetRedirectUri()

	err = u.s.LinkOAuthProvider(ctx, userID, provider, code, redirectUri)
	if err != nil {
		return nil, FromUserToHTTPError(err)
	}

	return &v1.Auth_LinkProviderResponse{}, nil
}

// UnlinkProvider is a method that handles unlinking an OAuth provider from a user
func (u *AuthHttpApi) UnlinkProvider(ctx context.Context, req *v1.Auth_UnlinkProviderRequest) (*v1.Auth_UnlinkProviderResponse, error) {
	// Extract user ID from context (set by auth middleware)
	userID, err := u.getUserIDFromContext(ctx)
	if err != nil {
		return nil, NewError(http.StatusUnauthorized, "unauthorized")
	}

	provider := req.GetProvider()

	err = u.s.UnlinkOAuthProvider(ctx, userID, provider)
	if err != nil {
		return nil, FromUserToHTTPError(err)
	}

	return &v1.Auth_UnlinkProviderResponse{}, nil
}

// GetProviders is a method that handles getting all linked OAuth providers for a user
func (u *AuthHttpApi) GetProviders(ctx context.Context, req *v1.Auth_GetProvidersRequest) (*v1.Auth_GetProvidersResponse, error) {
	// Extract user ID from context (set by auth middleware)
	userID, err := u.getUserIDFromContext(ctx)
	if err != nil {
		return nil, NewError(http.StatusUnauthorized, "unauthorized")
	}

	providers, err := u.s.GetUserLinkedProviders(ctx, userID)
	if err != nil {
		return nil, FromUserToHTTPError(err)
	}

	return &v1.Auth_GetProvidersResponse{
		Providers: providers,
	}, nil
}

// getUserIDFromContext extracts the user ID from the context
// This assumes an auth middleware has set the user ID in the context
func (u *AuthHttpApi) getUserIDFromContext(ctx context.Context) (uint, error) {
	// Try to extract from JWT claims in metadata
	// This is a simplified version - in production, you'd use proper middleware
	// that validates the JWT and sets the user ID in context
	userID, ok := ctx.Value("user_id").(uint)
	if !ok {
		return 0, errors.New("user_id not found in context")
	}
	return userID, nil
}
