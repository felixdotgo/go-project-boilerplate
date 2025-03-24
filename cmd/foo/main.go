package main

import (
	"net/http"

	"github.com/0x46656C6978/go-project-boilerplate/cmd/api/config"
	"github.com/0x46656C6978/go-project-boilerplate/pkg/conv"
	"github.com/0x46656C6978/go-project-boilerplate/pkg/log"
	"github.com/fvbock/endless"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
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

	logger := log.NewLogger(logDebug)

	mux := runtime.NewServeMux()
	runtime.WithMiddlewares(logMiddleware(logger))(mux)

	mux.HandlePath(http.MethodGet, "/", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	})

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
