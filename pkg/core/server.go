package core

import (
	"context"
	"net/http"

	"github.com/0x46656C6978/go-project-boilerplate/pkg/log"
	"github.com/fvbock/endless"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

// HandlerServer is a struct that wraps the http.ServeMux
// contains all methods that will be used to bring up the server
type HandlerServer struct {
	logger *log.Logger
	mux    *runtime.ServeMux
}

// NewHandlerServer returns a new http.ServeMux
func NewHandlerServer(logger *log.Logger) *HandlerServer {
	mux := runtime.NewServeMux(
		runtime.WithRoutingErrorHandler(func(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, code int) {
			// gprc gateway don't generate OPTIONS requests, when mux checking with list handlers it will fail
			// we're adding this condition to handle OPTIONS requests
			if r.Method == http.MethodOptions {
				w.Header().Set("Access-Control-Allow-Headers", "*")
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, HEAD, OPTIONS")
				w.WriteHeader(http.StatusOK)
				return
			}
		}),
	)
	runtime.WithMiddlewares(handlerServerLogMiddleware(logger))(mux)

	return &HandlerServer{
		logger: logger,
		mux:    mux,
	}
}

func handlerServerLogMiddleware(l *log.Logger) func(h runtime.HandlerFunc) runtime.HandlerFunc {
	return func(h runtime.HandlerFunc) runtime.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request, params map[string]string) {
			l.With(
				"protocol", r.Proto,
				"method", r.Method,
				"host", r.Host,
				"path", r.URL.Path,
				"uri", r.RequestURI,
				"remoteAddr", r.RemoteAddr,
			).Info("")
			h(w, r, params)
		}
	}
}

// GetMux returns the http.ServeMux
func (s *HandlerServer) GetMux() *runtime.ServeMux {
	if s == nil {
		return nil
	}
	return s.mux
}

// Run starts the server
func (s *HandlerServer) Run(port string) {
	engine := h2c.NewHandler(s.mux, &http2.Server{})
	err := endless.ListenAndServe(":"+port, engine)
	if err != nil {
		panic(err)
	}
}

// AddMiddlewares adds middlewares to the http.ServeMux
func (s *HandlerServer) AddMiddlewares(middlewares ...runtime.Middleware) {
	runtime.WithMiddlewares(middlewares...)(s.mux)
}

// EnableCORS enables CORS for service that need it
func (s *HandlerServer) EnableCORS() {
	s.AddMiddlewares(func(h runtime.HandlerFunc) runtime.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request, params map[string]string) {
			w.Header().Set("Access-Control-Allow-Headers", "*")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, HEAD, OPTIONS")
			h(w, r, params)
		}
	})
}
