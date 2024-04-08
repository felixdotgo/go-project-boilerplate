package dicontainer

import (
	"github.com/0x46656C6978/go-project-boilerplate/internal/config"
	"github.com/0x46656C6978/go-project-boilerplate/internal/httpapi"
	"github.com/0x46656C6978/go-project-boilerplate/internal/service"
	"go.uber.org/dig"
)

func ProvideHttpApis(c *dig.Container) {
	c.Provide(func(cfg *config.Config, srv service.UserServiceInterface) *httpapi.AuthHttpApi {
		return httpapi.NewAuthHttpApi(cfg, srv)
	})
}
