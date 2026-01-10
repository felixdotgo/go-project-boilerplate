package repository

import (
	"context"
	"errors"
	"time"

	"github.com/0x46656C6978/go-project-boilerplate/cmd/svc-auth/entity"
	"github.com/0x46656C6978/go-project-boilerplate/cmd/svc-auth/oauth"
	"github.com/0x46656C6978/go-project-boilerplate/cmd/svc-auth/repository/dto"
	"github.com/0x46656C6978/go-project-boilerplate/cmd/svc-auth/repository/model"
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
	var userModel model.User
	err := r.db.First(&userModel, "email = ?", email).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return dto.ToUserEntity(&userModel), nil
}

// FindByID returns a user by given id, return error if any
func (r *UserRepo) FindByID(ctx context.Context, id int) (*entity.User, error) {
	var userModel model.User
	err := r.db.First(&userModel, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return dto.ToUserEntity(&userModel), nil
}

// Save creates or updates a user
func (r *UserRepo) Save(ctx context.Context, user *entity.User) error {
	userModel := dto.FromUserEntity(user)
	err := r.db.Save(userModel).Error
	if err != nil {
		return err
	}
	// Update entity with generated ID if it was an insert
	user.ID = userModel.ID
	return nil
}

// Update updates an existing user
func (r *UserRepo) Update(ctx context.Context, user *entity.User) error {
	userModel := dto.FromUserEntity(user)
	return r.db.WithContext(ctx).Save(userModel).Error
}

// FindUserByProviderAndProviderUserID finds user by OAuth provider
func (r *UserRepo) FindUserByProviderAndProviderUserID(ctx context.Context, provider, providerUserID string) (*entity.User, error) {
	var oauthProviderModel model.UserOAuthProvider
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("provider = ? AND provider_user_id = ?", provider, providerUserID).
		First(&oauthProviderModel).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return dto.ToUserEntity(oauthProviderModel.User), nil
}

// CreateOAuthProvider links a new OAuth provider to a user
func (r *UserRepo) CreateOAuthProvider(ctx context.Context, userID int, provider string, encryptedAccessToken, encryptedRefreshToken *string, userInfo *oauth.UserInfo, tokenExpiry *time.Time) error {
	if userInfo == nil {
		return ErrBadRequest
	}

	oauthProvider := &model.UserOAuthProvider{
		UserID:         userID,
		Provider:       provider,
		ProviderUserID: userInfo.ProviderUserID,
		AccessToken:    encryptedAccessToken,
		RefreshToken:   encryptedRefreshToken,
		TokenExpiry:    tokenExpiry,
		Email:          &userInfo.Email,
		Name:           &userInfo.Name,
		FirstName:      &userInfo.FirstName,
		LastName:       &userInfo.LastName,
		AvatarURL:      &userInfo.AvatarURL,
		EmailVerified:  userInfo.EmailVerified,
		Locale:         &userInfo.Locale,
	}

	err := r.db.WithContext(ctx).Create(oauthProvider).Error
	if err != nil {
		return err
	}
	return nil
}

// UpdateOAuthProvider updates an existing OAuth provider link
func (r *UserRepo) UpdateOAuthProvider(ctx context.Context, oauthProvider *model.UserOAuthProvider) error {
	return r.db.WithContext(ctx).Save(oauthProvider).Error
}

// FindOAuthProvider retrieves a specific OAuth provider for a user
func (r *UserRepo) FindOAuthProvider(ctx context.Context, userID int, provider string) (*model.UserOAuthProvider, error) {
	var oauthProviderModel model.UserOAuthProvider
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND provider = ?", userID, provider).
		First(&oauthProviderModel).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOAuthProviderNotFound
		}
		return nil, err
	}
	return &oauthProviderModel, nil
}

// FindUserOAuthProviders retrieves all OAuth providers for a user
func (r *UserRepo) FindUserOAuthProviders(ctx context.Context, userID int) ([]model.UserOAuthProvider, error) {
	var providerModels []model.UserOAuthProvider
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&providerModels).Error

	if err != nil {
		return nil, err
	}

	// Convert models to entities
	providers := make([]model.UserOAuthProvider, 0, len(providerModels))
	for _, model := range providerModels {
		providers = append(providers, model)
	}
	return providers, nil
}

// DeleteOAuthProvider unlinks an OAuth provider from a user
func (r *UserRepo) DeleteOAuthProvider(ctx context.Context, userID int, provider string) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND provider = ?", userID, provider).
		Delete(&model.UserOAuthProvider{}).Error
}

// CreateRefreshToken stores a new refresh token
func (r *UserRepo) CreateRefreshToken(ctx context.Context, token *model.RefreshToken) error {
	err := r.db.WithContext(ctx).Create(token).Error
	if err != nil {
		return err
	}
	return nil
}

// FindRefreshTokenByToken retrieves refresh token by hash
func (r *UserRepo) FindRefreshTokenByToken(ctx context.Context, tokenHash string) (*model.RefreshToken, error) {
	var tokenModel model.RefreshToken
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("token = ?", tokenHash).
		First(&tokenModel).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRefreshTokenNotFound
		}
		return nil, err
	}
	return &tokenModel, nil
}

// RevokeRefreshToken marks a token as revoked
func (r *UserRepo) RevokeRefreshToken(ctx context.Context, tokenID int) error {
	now := r.db.NowFunc()
	return r.db.WithContext(ctx).
		Model(&model.RefreshToken{}).
		Where("id = ?", tokenID).
		Update("revoked_at", now).Error
}

// RevokeAllUserRefreshTokens revokes all tokens for a user
func (r *UserRepo) RevokeAllUserRefreshTokens(ctx context.Context, userID int) error {
	now := r.db.NowFunc()
	return r.db.WithContext(ctx).
		Model(&model.RefreshToken{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", now).Error
}

// DeleteExpiredRefreshTokens cleanup job for expired/revoked tokens
func (r *UserRepo) DeleteExpiredRefreshTokens(ctx context.Context) error {
	now := r.db.NowFunc()
	return r.db.WithContext(ctx).
		Where("expires_at < ? OR revoked_at IS NOT NULL", now).
		Delete(&model.RefreshToken{}).Error
}
