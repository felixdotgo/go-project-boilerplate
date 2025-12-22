package entity

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User entity represents a real user in our application
type User struct {
	// Core identity
	ID       int            `json:"id" gorm:"primarykey"`
	Email    string         `json:"email" gorm:"uniqueIndex;not null"`
	Password *string        `json:"-" gorm:"type:varchar(255)"` // Nullable for OAuth-only users

	// User profile (aggregated from OAuth providers or self-entered)
	Name          *string `json:"name,omitempty" gorm:"type:varchar(255)"`
	FirstName     *string `json:"first_name,omitempty" gorm:"type:varchar(100)"`
	LastName      *string `json:"last_name,omitempty" gorm:"type:varchar(100)"`
	AvatarURL     *string `json:"avatar_url,omitempty" gorm:"type:text"`
	EmailVerified bool    `json:"email_verified" gorm:"default:false"`
	Locale        *string `json:"locale,omitempty" gorm:"type:varchar(10)"`

	// Timestamps
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	OAuthProviders []UserOAuthProvider `json:"-" gorm:"foreignKey:UserID"`
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

// GetDeletedAt return DeletedAt value
func (u *User) GetDeletedAt() gorm.DeletedAt {
	if u == nil {
		return gorm.DeletedAt{}
	}
	return u.DeletedAt
}
