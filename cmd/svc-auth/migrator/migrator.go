package migrator

type MigratorInterface interface {
	Up() error
	Down() error
	Create(name, ext string) error
}
