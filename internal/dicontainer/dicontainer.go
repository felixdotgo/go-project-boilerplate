package dicontainer

import (
	"github.com/0x46656C6978/go-project-boilerplate/internal/config"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

func ProvideCore(c *dig.Container) {
	c.Provide(func() *config.Config {
		return config.New()
	})
	c.Provide(func() (*zap.Logger, error) {
		return zap.NewProduction()
	})
}
