package repository

import (
	"context"

	"github.com/0x46656C6978/go-project-boilerplate/cmd/svc-auth/entity"
)

// UserRepoInterface is an interface define all methods that will be used to handle user
type UserRepoInterface interface {
	// User management
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByID(ctx context.Context, id int) (*entity.User, error)
	Save(ctx context.Context, user *entity.User) error
	Update(ctx context.Context, user *entity.User) error

	// OAuth provider management
	FindUserByProviderAndProviderUserID(ctx context.Context, provider, providerUserID string) (*entity.User, error)
	CreateOAuthProvider(ctx context.Context, oauthProvider *entity.UserOAuthProvider) error
	UpdateOAuthProvider(ctx context.Context, oauthProvider *entity.UserOAuthProvider) error
	FindOAuthProvider(ctx context.Context, userID int, provider string) (*entity.UserOAuthProvider, error)
	FindUserOAuthProviders(ctx context.Context, userID int) ([]entity.UserOAuthProvider, error)
	DeleteOAuthProvider(ctx context.Context, userID int, provider string) error

	// Refresh token management
	CreateRefreshToken(ctx context.Context, token *entity.RefreshToken) error
	FindRefreshTokenByToken(ctx context.Context, tokenHash string) (*entity.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, tokenID int) error
	RevokeAllUserRefreshTokens(ctx context.Context, userID int) error
	DeleteExpiredRefreshTokens(ctx context.Context) error
}
