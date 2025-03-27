package migrator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/gorm"
)

// Migrator is a wrapper around golang-migrate
// contains all methods that will be used to handle database migrations
type Migrator struct {
	m *migrate.Migrate
}

const MIGRATION_DIR = "file://migrations/sql"

// New creates a new migrator
func New(db *gorm.DB) (*Migrator, error) {
	conn, err := db.DB()
	if err != nil {
		return nil, err
	}

	driver, err := postgres.WithInstance(conn, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(MIGRATION_DIR, "postgres", driver)
	if err != nil {
		return nil, err
	}

	return &Migrator{
		m: m,
	}, nil
}

// Up runs all migrations
func (m *Migrator) Up() error {
	return m.m.Up()
}

// Down reverts all migrations
func (m *Migrator) Down() error {
	return m.m.Down()
}

// Create creates a new migration file
func (m *Migrator) Create(name, ext string) error {
	dir := strings.Replace(MIGRATION_DIR, "file://", "./", -1)
	dir = filepath.Clean(dir)
	ext = "." + strings.TrimPrefix(ext, ".")
	version := fmt.Sprintf("%d", time.Now().UnixMilli())
	versionGlob := filepath.Join(dir, version+"_*"+ext)
	matches, err := filepath.Glob(versionGlob)
	if err != nil {
		return err
	}

	if len(matches) > 0 {
		return fmt.Errorf("duplicate migration version: %s", version)
	}

	if err = os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	for _, direction := range []string{"up", "down"} {
		basename := fmt.Sprintf("%s_%s.%s%s", version, name, direction, ext)
		filename := filepath.Join(dir, basename)

		if _, err = os.Create(filename); err != nil {
			return err
		}
	}

	return nil
}
