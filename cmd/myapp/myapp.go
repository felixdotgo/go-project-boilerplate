package main

import (
	"fmt"
	"github.com/0x46656C6978/go-project-boilerplate/internal/repository"
	"github.com/0x46656C6978/go-project-boilerplate/internal/service"
	"gorm.io/driver/postgres"
	gormlogger "gorm.io/gorm/logger"

	"github.com/0x46656C6978/go-project-boilerplate/internal/config"
	"github.com/0x46656C6978/go-project-boilerplate/internal/httpapi"
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
			"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Etc/GMT",
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

	authRepo := repository.NewUserRepo(db)
	authSvc := service.NewUserService(authRepo)
	authApi := httpapi.NewAuthHttpApi(cfg, authSvc)

	server := httpapi.New(cfg, db, zlogger)
	server.Register(authApi)
	server.Run(cfg.GetPort())
}
