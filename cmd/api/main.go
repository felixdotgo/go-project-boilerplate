package main

import (
	"fmt"

	"github.com/0x46656C6978/go-project-boilerplate/cmd/api/config"
	"github.com/0x46656C6978/go-project-boilerplate/cmd/api/httpapi"
	"github.com/0x46656C6978/go-project-boilerplate/cmd/api/repository"
	"github.com/0x46656C6978/go-project-boilerplate/cmd/api/service"
	"gorm.io/driver/postgres"
	gormlogger "gorm.io/gorm/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}
	zlogger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Etc/GMT",
			cfg.DB.Host,
			cfg.DB.User,
			cfg.DB.Password,
			cfg.DB.DBName,
			cfg.DB.Port,
		),
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})

	if err != nil {
		panic(err)
	}

	authRepo := repository.NewUserRepo(db)
	authSvc := service.NewUserService(authRepo)
	authApi := httpapi.NewAuthHttpApi(cfg, authSvc)

	server := httpapi.New(cfg, db, zlogger)
	server.Register(authApi)
	server.Run(cfg.GetPort())
}
