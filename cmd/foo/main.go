package main

import (
	"net/http"

	"github.com/0x46656C6978/go-project-boilerplate/cmd/api/config"
	"github.com/0x46656C6978/go-project-boilerplate/pkg/conv"
	"github.com/0x46656C6978/go-project-boilerplate/pkg/core"
	"github.com/0x46656C6978/go-project-boilerplate/pkg/log"
)

func main() {
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

	// Run handler server
	server := core.NewHandlerServer(logger)

	server.GetMux().HandlePath(http.MethodGet, "/", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	})

	server.Run(conv.ToString(cfg.Port))
}
