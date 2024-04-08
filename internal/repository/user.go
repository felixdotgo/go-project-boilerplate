package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/0x46656C6978/go-project-boilerplate/internal/entity"
)

type UserRepoInterface interface {
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByID(ctx context.Context, id int) (*entity.User, error)
	Save(ctx context.Context, user *entity.User) error
}

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepoInterface {
	return &UserRepo{
		db: db,
	}
}

func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user *entity.User
	err := r.db.First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepo) FindByID(ctx context.Context, id int) (*entity.User, error) {
	var user *entity.User
	err := r.db.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepo) Save(ctx context.Context, user *entity.User) error {
	return r.db.Save(user).Error
}
