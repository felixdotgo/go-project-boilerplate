package dicontainer

import (
	"github.com/0x46656C6978/go-project-boilerplate/internal/repository"
	"github.com/0x46656C6978/go-project-boilerplate/internal/service"
	"go.uber.org/dig"
)

func ProvideServices(c *dig.Container) {
	c.Provide(func(repo repository.UserRepoInterface) service.UserServiceInterface {
		return service.NewUserService(repo)
	})
}
