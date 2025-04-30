package repository

import (
	"context"

	"github.com/0x46656C6978/go-project-boilerplate/cmd/svc-auth/entity"

	"gorm.io/gorm"
)

// UserRepoInterface is an interface define all methods that will be used to handle user
type UserRepoInterface interface {
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByID(ctx context.Context, id int) (*entity.User, error)
	Save(ctx context.Context, user *entity.User) error
}

// UserRepo is a struct that implements UserRepoInterface
type UserRepo struct {
	db *gorm.DB
}

// NewUserRepo creates a new UserRepo
func NewUserRepo(db *gorm.DB) UserRepoInterface {
	return &UserRepo{
		db: db,
	}
}

// FindByEmail returns a user by given email, return error if any
func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user *entity.User
	err := r.db.First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

// FindByID returns a user by given id, return error if any
func (r *UserRepo) FindByID(ctx context.Context, id int) (*entity.User, error) {
	var user *entity.User
	err := r.db.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Save creates or updates a user
func (r *UserRepo) Save(ctx context.Context, user *entity.User) error {
	return r.db.Save(user).Error
}
