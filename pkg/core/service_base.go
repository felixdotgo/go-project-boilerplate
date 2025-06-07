package core

import (
	"flag"
	"regexp"

	"github.com/0x46656C6978/go-project-boilerplate/pkg/log"
)

// ServiceBase is a struct that contains the logger
// and other common methods that will be used by all services
type ServiceBase struct {
	Logger *log.Logger
}

// NewService returns a new ServiceBase
func NewService(name string) *ServiceBase {
	debug := false
	if flag.Lookup("debug") == nil {
		debug = *flag.Bool("debug", false, "sets log level to debug")
	}
	flag.Parse()

	l := log.NewLogger(debug)

	re, err := regexp.Compile(`^[a-zA-Z.]+$`)
	if err != nil {
		l.Panic("failed to compile regex")
	}
	if !re.MatchString(name) {
		l.Panic("service name must be A-Z, a-z, .")
	}

	l = l.With("serviceName", name)

	return &ServiceBase{
		Logger: l,
	}
}

