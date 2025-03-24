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
	"github.com/0x46656C6978/go-project-boilerplate/pkg/log"
	authv1 "github.com/0x46656C6978/go-project-boilerplate/rpc/api/auth/v1"
	"github.com/fvbock/endless"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"gorm.io/driver/postgres"
	gormlogger "gorm.io/gorm/logger"

	"gorm.io/gorm"
)

func main() {
	// Load config
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	logDebug := true
	if cfg.GetEnvMode() == config.ENV_PRODUCTION {
		logDebug = false
	}

	ctx := context.Background()
	logger := log.NewLogger(logDebug)

	db, err := getDBConnection(cfg)
	if err != nil {
		panic(err)
	}

	authRepo := repository.NewUserRepo(db)
	authSvc := service.NewUserService(authRepo)
	authApi := httpapi.NewAuthServiceServer(cfg, authSvc)

	mux := runtime.NewServeMux()
	runtime.WithMiddlewares(logMiddleware(logger))(mux)

	// APIs register
	authv1.RegisterAuthServiceHandlerServer(ctx, mux, authApi)

	engine := h2c.NewHandler(mux, &http2.Server{})
	err = endless.ListenAndServe(":"+conv.ToString(cfg.Port), engine)
	if err != nil {
		panic(err)
	}
}

func logMiddleware(l *log.Logger) (func(h runtime.HandlerFunc) runtime.HandlerFunc) {
	return func(h runtime.HandlerFunc) runtime.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request, params map[string]string) {
			l.With(
				"protocol", r.Proto,
				"method", r.Method,
				"host", r.Host,
				"path", r.URL.Path,
				"uri", r.RequestURI,
				"remote_addr", r.RemoteAddr,
			).Info("")
			h(w, r, params)
		}
	}
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
