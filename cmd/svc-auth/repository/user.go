package repository

import (
	"context"

	"github.com/0x46656C6978/go-project-boilerplate/cmd/svc-auth/entity"
)

// UserRepoInterface is an interface define all methods that will be used to handle user
type UserRepoInterface interface {
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByID(ctx context.Context, id int) (*entity.User, error)
	Save(ctx context.Context, user *entity.User) error
}
