package core

import (
	"flag"
	"regexp"

	"github.com/0x46656C6978/go-project-boilerplate/pkg/log"
)

// RepositoryBase is a struct that contains the logger
// and other common methods that will be used by all repositories
type RepositoryBase struct {
	Logger *log.Logger
}

// NewRepository returns a new RepositoryBase
func NewRepository(name string) *RepositoryBase {
    debug := flag.Bool("debug", false, "sets log level to debug")
    flag.Parse()

	l := log.NewLogger(*debug)

	re, err := regexp.Compile(`^[a-zA-Z.]+$`)
	if err != nil {
		l.Panic("failed to compile regex")
	}
	if !re.MatchString(name) {
		l.Panic("repository name must be A-Z, a-z, .")
	}

	l = l.With("repositoryName", name)

	return &RepositoryBase{
		Logger: l,
	}
}

