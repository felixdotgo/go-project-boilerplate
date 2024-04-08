package dicontainer

import (
	"github.com/0x46656C6978/go-project-boilerplate/internal/repository"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

func ProvideRepositories(c *dig.Container) {
	c.Provide(func(db *gorm.DB) repository.UserRepoInterface {
		return repository.NewUserRepo(db)
	})
}
