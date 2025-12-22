package entity

import "time"

// UserOAuthProvider entity represents an OAuth provider linked to a user account
// Users can have multiple OAuth providers (Google, Facebook, GitHub, etc.)
type UserOAuthProvider struct {
	ID             int    `json:"id" gorm:"primarykey"`
	UserID         int    `json:"user_id" gorm:"not null;index:idx_user_provider;index:idx_user_id"`
	Provider       string `json:"provider" gorm:"type:varchar(50);not null;index:idx_user_provider;index:idx_provider_user_id"`
	ProviderUserID string `json:"provider_user_id" gorm:"type:varchar(255);not null;index:idx_provider_user_id"`

	// OAuth tokens (encrypted in database)
	AccessToken  *string    `json:"-" gorm:"type:text"`
	RefreshToken *string    `json:"-" gorm:"type:text"`
	TokenExpiry  *time.Time `json:"token_expiry,omitempty"`

	// Provider-specific profile data
	Email         *string `json:"email,omitempty" gorm:"type:varchar(255)"`
	Name          *string `json:"name,omitempty" gorm:"type:varchar(255)"`
	FirstName     *string `json:"first_name,omitempty" gorm:"type:varchar(100)"`
	LastName      *string `json:"last_name,omitempty" gorm:"type:varchar(100)"`
	AvatarURL     *string `json:"avatar_url,omitempty" gorm:"type:text"`
	EmailVerified bool    `json:"email_verified" gorm:"default:false"`
	Locale        *string `json:"locale,omitempty" gorm:"type:varchar(10)"`

	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	User User `json:"-" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// TableName overrides the table name
func (UserOAuthProvider) TableName() string {
	return "user_oauth_providers"
}
