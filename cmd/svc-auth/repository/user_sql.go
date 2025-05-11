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
