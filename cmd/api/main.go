package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/0x46656C6978/go-project-boilerplate/cmd/api/config"
	"github.com/0x46656C6978/go-project-boilerplate/cmd/api/httpapi"
	"github.com/0x46656C6978/go-project-boilerplate/cmd/api/repository"
	"github.com/0x46656C6978/go-project-boilerplate/cmd/api/service"
	"github.com/0x46656C6978/go-project-boilerplate/pkg/conv"
	authv1 "github.com/0x46656C6978/go-project-boilerplate/rpc/api/auth/v1"
	"github.com/fvbock/endless"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"gorm.io/driver/postgres"
	gormlogger "gorm.io/gorm/logger"

	"gorm.io/gorm"
)

func main() {
	ctx := context.Background()

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
	authApi := httpapi.NewAuthServiceServer(cfg, authSvc)

	logMiddleware := func(h runtime.HandlerFunc) runtime.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request, params map[string]string) {
			zlogger.Info("request", zap.String("method", r.Method), zap.String("path", r.URL.Path))
			h(w, r, params)
		}
	}

	mux := runtime.NewServeMux()
	runtime.WithMiddlewares(logMiddleware)(mux)
	authv1.RegisterAuthServiceHandlerServer(ctx, mux, authApi)

	engine := h2c.NewHandler(mux, &http2.Server{})
	err = endless.ListenAndServe(":"+conv.ToString(cfg.Port), engine)
	if err != nil {
		panic(err)
	}
}
