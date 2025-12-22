package repository

import (
	"context"
	"errors"

	"github.com/0x46656C6978/go-project-boilerplate/cmd/svc-auth/entity"
	"github.com/0x46656C6978/go-project-boilerplate/pkg/core"
	"gorm.io/gorm"
)

// UserRepo is a struct that implements UserRepoInterface
type UserRepo struct {
	*core.RepositoryBase
	db *gorm.DB
}

// NewUserRepo creates a new UserRepo
func NewUserRepo(db *gorm.DB) UserRepoInterface {
	return &UserRepo{
		core.NewRepository("user"),
		db,
	}
}

// FindByEmail returns a user by given email, return error if any
func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user *entity.User
	err := r.db.First(&user, "email = ?", email).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return user, nil
}

// FindByID returns a user by given id, return error if any
func (r *UserRepo) FindByID(ctx context.Context, id int) (*entity.User, error) {
	var user *entity.User
	err := r.db.First(&user, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return user, nil
}

// Save creates or updates a user
func (r *UserRepo) Save(ctx context.Context, user *entity.User) error {
	return r.db.Save(user).Error
}

// Update updates an existing user
func (r *UserRepo) Update(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// FindUserByProviderAndProviderUserID finds user by OAuth provider
func (r *UserRepo) FindUserByProviderAndProviderUserID(ctx context.Context, provider, providerUserID string) (*entity.User, error) {
	var oauthProvider entity.UserOAuthProvider
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("provider = ? AND provider_user_id = ?", provider, providerUserID).
		First(&oauthProvider).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &oauthProvider.User, nil
}

// CreateOAuthProvider links a new OAuth provider to a user
func (r *UserRepo) CreateOAuthProvider(ctx context.Context, oauthProvider *entity.UserOAuthProvider) error {
	return r.db.WithContext(ctx).Create(oauthProvider).Error
}

// UpdateOAuthProvider updates an existing OAuth provider link
func (r *UserRepo) UpdateOAuthProvider(ctx context.Context, oauthProvider *entity.UserOAuthProvider) error {
	return r.db.WithContext(ctx).Save(oauthProvider).Error
}

// FindOAuthProvider retrieves a specific OAuth provider for a user
func (r *UserRepo) FindOAuthProvider(ctx context.Context, userID int, provider string) (*entity.UserOAuthProvider, error) {
	var oauthProvider entity.UserOAuthProvider
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND provider = ?", userID, provider).
		First(&oauthProvider).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOAuthProviderNotFound
		}
		return nil, err
	}
	return &oauthProvider, nil
}

// FindUserOAuthProviders retrieves all OAuth providers for a user
func (r *UserRepo) FindUserOAuthProviders(ctx context.Context, userID int) ([]entity.UserOAuthProvider, error) {
	var providers []entity.UserOAuthProvider
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&providers).Error
	return providers, err
}

// DeleteOAuthProvider unlinks an OAuth provider from a user
func (r *UserRepo) DeleteOAuthProvider(ctx context.Context, userID int, provider string) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND provider = ?", userID, provider).
		Delete(&entity.UserOAuthProvider{}).Error
}

// CreateRefreshToken stores a new refresh token
func (r *UserRepo) CreateRefreshToken(ctx context.Context, token *entity.RefreshToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

// FindRefreshTokenByToken retrieves refresh token by hash
func (r *UserRepo) FindRefreshTokenByToken(ctx context.Context, tokenHash string) (*entity.RefreshToken, error) {
	var token entity.RefreshToken
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("token = ?", tokenHash).
		First(&token).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRefreshTokenNotFound
		}
		return nil, err
	}
	return &token, nil
}

// RevokeRefreshToken marks a token as revoked
func (r *UserRepo) RevokeRefreshToken(ctx context.Context, tokenID int) error {
	now := r.db.NowFunc()
	return r.db.WithContext(ctx).
		Model(&entity.RefreshToken{}).
		Where("id = ?", tokenID).
		Update("revoked_at", now).Error
}

// RevokeAllUserRefreshTokens revokes all tokens for a user
func (r *UserRepo) RevokeAllUserRefreshTokens(ctx context.Context, userID int) error {
	now := r.db.NowFunc()
	return r.db.WithContext(ctx).
		Model(&entity.RefreshToken{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", now).Error
}

// DeleteExpiredRefreshTokens cleanup job for expired/revoked tokens
func (r *UserRepo) DeleteExpiredRefreshTokens(ctx context.Context) error {
	now := r.db.NowFunc()
	return r.db.WithContext(ctx).
		Where("expires_at < ? OR revoked_at IS NOT NULL", now).
		Delete(&entity.RefreshToken{}).Error
}
