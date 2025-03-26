package main

import (
	"context"
	"fmt"

	"github.com/0x46656C6978/go-project-boilerplate/cmd/api/config"
	"github.com/0x46656C6978/go-project-boilerplate/cmd/api/httpapi"
	"github.com/0x46656C6978/go-project-boilerplate/cmd/api/repository"
	"github.com/0x46656C6978/go-project-boilerplate/cmd/api/service"
	"github.com/0x46656C6978/go-project-boilerplate/pkg/conv"
	"github.com/0x46656C6978/go-project-boilerplate/pkg/core"
	"github.com/0x46656C6978/go-project-boilerplate/pkg/log"
	authv1 "github.com/0x46656C6978/go-project-boilerplate/rpc/api/auth/v1"
	"gorm.io/driver/postgres"
	gormlogger "gorm.io/gorm/logger"

	"gorm.io/gorm"
)

func main() {
	ctx := context.Background()

	// Load config
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	// Logger initialization
	logDebug := true
	if cfg.GetEnvMode() == config.ENV_PRODUCTION {
		logDebug = false
	}
	logger := log.NewLogger(logDebug)

	// DB connection
	db, err := getDBConnection(cfg)
	if err != nil {
		panic(err)
	}

	// Run handler server
	server := core.NewHandlerServer(logger)
	serviceServerRegistry(ctx, server, cfg, db)
	server.Run(conv.ToString(cfg.Port))
}

func getDBConnection(cfg *config.Config) (*gorm.DB, error) {
	return gorm.Open(postgres.New(postgres.Config{
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
}

func serviceServerRegistry(ctx context.Context, server *core.HandlerServer, cfg *config.Config, db *gorm.DB) {
	authRepo := repository.NewUserRepo(db)
	authSvc := service.NewUserService(authRepo)
	authApi := httpapi.NewAuthServiceServer(cfg, authSvc)

	// APIs register
	authv1.RegisterAuthServiceHandlerServer(ctx, server.GetMux(), authApi)
}
