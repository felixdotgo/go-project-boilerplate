package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User entity represents a real user in our application (domain layer)
type User struct {
	// Core identity
	ID       int     `json:"id"`
	Email    string  `json:"email"`
	Password *string `json:"-"` // Nullable for OAuth-only users

	// User profile (aggregated from OAuth providers or self-entered)
	Name          *string `json:"name,omitempty"`
	FirstName     *string `json:"first_name,omitempty"`
	LastName      *string `json:"last_name,omitempty"`
	AvatarURL     *string `json:"avatar_url,omitempty"`
	EmailVerified bool    `json:"email_verified"`
	Locale        *string `json:"locale,omitempty"`

	// Timestamps
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// SetPassword convert plaintext password into encrypted password
func (u *User) SetPassword(password string) error {
	if u == nil {
		return nil
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	hashedPassword := string(bytes)
	u.Password = &hashedPassword
	return nil
}

// IsValidPassword check given password is correct with what was set in Password field
func (u *User) IsValidPassword(password string) bool {
	if u == nil || u.Password == nil {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(*u.Password), []byte(password))
	return err == nil
}

// HasPassword checks if user has a password set (for email/password authentication)
func (u *User) HasPassword() bool {
	return u != nil && u.Password != nil && *u.Password != ""
}

// GetID return ID value
func (u *User) GetID() int {
	if u == nil {
		return 0
	}
	return u.ID
}

// GetEmail return Email value
func (u *User) GetEmail() string {
	if u == nil {
		return ""
	}
	return u.Email
}

// GetPassword return Password value
func (u *User) GetPassword() string {
	if u == nil || u.Password == nil {
		return ""
	}
	return *u.Password
}

// GetCreatedAt return CreatedAt value
func (u *User) GetCreatedAt() time.Time {
	if u == nil {
		return time.Time{}
	}
	return u.CreatedAt
}

// GetUpdatedAt return UpdatedAt value
func (u *User) GetUpdatedAt() time.Time {
	if u == nil {
		return time.Time{}
	}
	return u.UpdatedAt
}


