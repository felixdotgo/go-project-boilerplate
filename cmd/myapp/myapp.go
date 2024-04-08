package main

import (
	"fmt"

	"github.com/0x46656C6978/go-project-boilerplate/internal/config"
	"github.com/0x46656C6978/go-project-boilerplate/internal/dicontainer"
	"github.com/0x46656C6978/go-project-boilerplate/internal/httpapi"
	"go.uber.org/dig"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func main() {
	c := dig.New()
	dicontainer.ProvideCore(c)
	dicontainer.ProvideDB(c)
	dicontainer.ProvideRepositories(c)
	dicontainer.ProvideServices(c)
	dicontainer.ProvideHttpApis(c)

	err := c.Invoke(func(
		cfg *config.Config,
		db *gorm.DB,
		logger *zap.Logger,
		authApi *httpapi.AuthHttpApi,
	){
		server := httpapi.New(cfg, db, logger)
		server.Register(authApi)
		server.Run(cfg.Port)
	})
	if err != nil {
		fmt.Println(err.Error())
	}
}
