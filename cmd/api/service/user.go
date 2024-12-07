package service

import (
	"context"
	"errors"
	"github.com/0x46656C6978/go-project-boilerplate/cmd/api/entity"
	"github.com/0x46656C6978/go-project-boilerplate/cmd/api/repository"
)

type UserServiceInterface interface {
	Create(ctx context.Context, user *entity.User) error
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByID(ctx context.Context, id int) (*entity.User, error)
	Save(ctx context.Context, user *entity.User) error
	VerifyCredentials(ctx context.Context, user *entity.User, email, password string) error
}

type UserService struct {
	r repository.UserRepoInterface
}

func NewUserService(r repository.UserRepoInterface) UserServiceInterface {
	return &UserService{r: r}
}

func (u *UserService) Create(ctx context.Context, user *entity.User) error {
	return u.r.Save(ctx, user)
}

func (u *UserService) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	return u.r.FindByEmail(ctx, email)
}

func (u *UserService) FindByID(ctx context.Context, id int) (*entity.User, error) {
	return u.r.FindByID(ctx, id)
}

func (u *UserService) Save(ctx context.Context, user *entity.User) error {
	return u.r.Save(ctx, user)
}

func (u *UserService) VerifyCredentials(ctx context.Context, user *entity.User, email, password string) error {
	// comparing user email
	if user.Email != email {
		return errors.New("invalid email")
	}
	// comparing user password
	if !user.IsValidPassword(password) {
		return errors.New("invalid password")
	}
	return nil
}
