package core

import (
	"flag"

	"github.com/0x46656C6978/go-project-boilerplate/pkg/log"
)

// ServiceBase is a struct that contains the logger
// and other common methods that will be used by all services
type ServiceBase struct {
	Logger *log.Logger
}

// NewService returns a new ServiceBase
func NewService(name string) *ServiceBase {
    debug := flag.Bool("debug", false, "sets log level to debug")
    flag.Parse()

	l := log.NewLogger(*debug).With("serviceName", name)

	return &ServiceBase{
		Logger: l,
	}
}

