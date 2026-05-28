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
	m   *migrate.Migrate
	dir string
}

// New creates a new migrator
func New(db *gorm.DB, migrationDir string) (MigratorInterface, error) {
	conn, err := db.DB()
	if err != nil {
		return nil, err
	}

	absDir, err := filepath.Abs(migrationDir)
	if err != nil {
		return nil, err
	}

	driver, err := postgres.WithInstance(conn, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	migrationURL := "file://" + filepath.ToSlash(absDir)
	m, err := migrate.NewWithDatabaseInstance(migrationURL, "postgres", driver)
	if err != nil {
		return nil, err
	}

	return &Migrator{
		m:   m,
		dir: absDir,
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
	ext = "." + strings.TrimPrefix(ext, ".")
	version := fmt.Sprintf("%d", time.Now().UnixMilli())
	versionGlob := filepath.Join(m.dir, version+"_*"+ext)
	matches, err := filepath.Glob(versionGlob)
	if err != nil {
		return err
	}

	if len(matches) > 0 {
		return fmt.Errorf("duplicate migration version: %s", version)
	}

	if err = os.MkdirAll(m.dir, os.ModePerm); err != nil {
		return err
	}

	for _, direction := range []string{"up", "down"} {
		basename := fmt.Sprintf("%s_%s.%s%s", version, name, direction, ext)
		filename := filepath.Join(m.dir, basename)

		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		if err := file.Close(); err != nil {
			return err
		}
	}

	return nil
}
