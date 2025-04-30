package service

import (
	"context"
	"errors"

	"github.com/0x46656C6978/go-project-boilerplate/cmd/svc-auth/entity"
	"github.com/0x46656C6978/go-project-boilerplate/cmd/svc-auth/repository"
	"github.com/0x46656C6978/go-project-boilerplate/pkg/core"
)

// UserServiceInterface is an interface define all methods that will be used to handle user
type UserServiceInterface interface {
	Create(ctx context.Context, user *entity.User) error
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByID(ctx context.Context, id int) (*entity.User, error)
	Save(ctx context.Context, user *entity.User) error
	VerifyCredentials(ctx context.Context, user *entity.User, email, password string) error
}

// UserService is a struct that implements UserServiceInterface
type UserService struct {
	*core.ServiceBase
	r repository.UserRepoInterface
}

// NewUserService creates a new UserService
func NewUserService(userRepo repository.UserRepoInterface) UserServiceInterface {
	return &UserService{
		core.NewService("user"),
		userRepo,
	}
}

// Create creates a new user
func (u *UserService) Create(ctx context.Context, user *entity.User) error {
	return u.r.Save(ctx, user)
}

// FindByEmail returns a user by given email, return error if any
func (u *UserService) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	return u.r.FindByEmail(ctx, email)
}

// FindByID returns a user by given id, return error if any
func (u *UserService) FindByID(ctx context.Context, id int) (*entity.User, error) {
	return u.r.FindByID(ctx, id)
}

// Save creates or updates a user
func (u *UserService) Save(ctx context.Context, user *entity.User) error {
	return u.r.Save(ctx, user)
}

// VerifyCredentials verify user credentials
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
