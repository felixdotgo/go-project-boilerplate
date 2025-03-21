package core

import (
	"flag"

	"github.com/0x46656C6978/go-project-boilerplate/pkg/log"
)

type ServiceBase struct {
	Logger *log.Logger
}

func NewService(name string) *ServiceBase {
    debug := flag.Bool("debug", false, "sets log level to debug")
    flag.Parse()

	l := log.NewLogger(*debug).With("serviceName", name)

	return &ServiceBase{
		Logger: l,
	}
}

