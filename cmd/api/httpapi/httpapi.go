package httpapi

import (
	"github.com/0x46656C6978/go-project-boilerplate/cmd/api/config"
	"time"

	"github.com/0x46656C6978/go-project-boilerplate/pkg/conv"
	"github.com/fvbock/endless"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type HttpApiInterface interface {
	RegisterApiEndpoints(engine *gin.Engine)
}

type Server struct {
	engine *gin.Engine
	db     *gorm.DB
}

func (s *Server) Register(api HttpApiInterface) {
	api.RegisterApiEndpoints(s.engine)
}

func (s *Server) Run(port int) {
	err := endless.ListenAndServe(":"+conv.ToString(port), s.engine)
	panic(err)
}

func New(cfg *config.Config, db *gorm.DB, logger *zap.Logger) *Server {
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	engine.Use(ginzap.RecoveryWithZap(logger, true))

	return &Server{
		engine: engine,
		db:     db,
	}
}
