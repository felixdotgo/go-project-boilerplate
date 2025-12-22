package entity

import (
	"time"
)

// RefreshToken entity represents a refresh token for OAuth2 token refresh flow
type RefreshToken struct {
	ID        int        `json:"id" gorm:"primarykey"`
	UserID    int        `json:"user_id" gorm:"not null;index"`
	Token     string     `json:"-" gorm:"uniqueIndex;type:varchar(255);not null"` // Hashed token
	ExpiresAt time.Time  `json:"expires_at" gorm:"not null;index"`
	CreatedAt time.Time  `json:"created_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty" gorm:"index"`

	// Relationships
	User User `json:"-" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// TableName overrides the table name
func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

// IsValid checks if token is still valid (not expired and not revoked)
func (rt *RefreshToken) IsValid() bool {
	if rt == nil {
		return false
	}
	return rt.RevokedAt == nil && time.Now().Before(rt.ExpiresAt)
}

// IsExpired checks if token is expired
func (rt *RefreshToken) IsExpired() bool {
	if rt == nil {
		return true
	}
	return time.Now().After(rt.ExpiresAt)
}

// IsRevoked checks if token is revoked
func (rt *RefreshToken) IsRevoked() bool {
	if rt == nil {
		return false
	}
	return rt.RevokedAt != nil
}
