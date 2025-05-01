package main

import (
	"fmt"

	"github.com/0x46656C6978/go-project-boilerplate/cmd/svc-auth/config"
	"github.com/0x46656C6978/go-project-boilerplate/cmd/svc-auth/migrator"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/spf13/cobra"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Etc/GMT",
			cfg.DB.Host,
			cfg.DB.User,
			cfg.DB.Password,
			cfg.DB.DBName,
			cfg.DB.Port,
		),
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic(err)
	}

	m, err := migrator.New(db)
	if err != nil {
		panic(err)
	}

	cmd := &cobra.Command{}
	cmd.AddCommand(createCmd(m))
	cmd.AddCommand(upCmd(m))
	cmd.AddCommand(downCmd(m))
	err = cmd.Execute()
	if err != nil {
		return
	}
}

func createCmd(m migrator.MigratorInterface) *cobra.Command {
	migrationName := ""
	migrationExtension := ""
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create new migration file",
		Run: func(cmd *cobra.Command, args []string) {
			err := m.Create(migrationName, migrationExtension)
			if err != nil {
				panic(err)
			}
		},
	}
	cmd.Flags().StringVarP(&migrationName, "name", "n", "", "migration name")
	cmd.Flags().StringVarP(&migrationExtension, "ext", "e", "sql", "file extension")
	err := cmd.MarkFlagRequired("name")
	if err != nil {
		return nil
	}
	return cmd
}

func upCmd(m migrator.MigratorInterface) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "up",
		Short: "Execute migrations",
		Run: func(cmd *cobra.Command, args []string) {
			err := m.Up()
			if err != nil {
				panic(err)
			}
		},
	}
	return cmd
}

func downCmd(m migrator.MigratorInterface) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "down",
		Short: "Revert all migrations",
		Run: func(cmd *cobra.Command, args []string) {
			err := m.Down()
			if err != nil {
				panic(err)
			}
		},
	}
	return cmd
}
