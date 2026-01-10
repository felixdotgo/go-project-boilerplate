package model

import (
	"time"
)

// RefreshToken entity represents a refresh token for OAuth2 token refresh flow (domain layer)
type RefreshToken struct {
	ID        int        `json:"id"`
	UserID    int        `json:"user_id"`
	Token     string     `json:"-"` // Hashed token
	ExpiresAt time.Time  `json:"expires_at"`
	CreatedAt time.Time  `json:"created_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
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
