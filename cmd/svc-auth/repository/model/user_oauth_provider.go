package model

import "time"

// UserOAuthProvider entity represents an OAuth provider linked to a user account (domain layer)
// Users can have multiple OAuth providers (Google, Facebook, GitHub, etc.)
type UserOAuthProvider struct {
	ID             int    `json:"id"`
	UserID         int    `json:"user_id"`
	Provider       string `json:"provider"`
	ProviderUserID string `json:"provider_user_id"`

	// OAuth tokens (encrypted in database)
	AccessToken  *string    `json:"-"`
	RefreshToken *string    `json:"-"`
	TokenExpiry  *time.Time `json:"token_expiry,omitempty"`

	// Provider-specific profile data
	Email         *string `json:"email,omitempty"`
	Name          *string `json:"name,omitempty"`
	FirstName     *string `json:"first_name,omitempty"`
	LastName      *string `json:"last_name,omitempty"`
	AvatarURL     *string `json:"avatar_url,omitempty"`
	EmailVerified bool    `json:"email_verified"`
	Locale        *string `json:"locale,omitempty"`

	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Associations
	User *User `gorm:"foreignKey:UserID;references:ID"`
}

// TableName sets the insert table name for this struct type
func (UserOAuthProvider) TableName() string {
	return "user_oauth_providers"
}
