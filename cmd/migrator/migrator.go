package main

import (
	"fmt"

	"github.com/0x46656C6978/go-project-boilerplate/internal/dicontainer"
	"github.com/0x46656C6978/go-project-boilerplate/internal/migrator"
	"github.com/spf13/cobra"
	"go.uber.org/dig"
)

func main() {
	c := dig.New()
	dicontainer.ProvideCore(c)
	dicontainer.ProvideDB(c)

	cmd := &cobra.Command{}
	cmd.AddCommand(createCmd(c))
	cmd.AddCommand(upCmd(c))
	cmd.AddCommand(downCmd(c))
	cmd.Execute()
}

func createCmd(c *dig.Container) *cobra.Command {
	migrationName := ""
	migrationExtension := ""
	createCmd := &cobra.Command{
		Use:     "create",
		Short:   "Create new migration file",
		Run: func(cmd *cobra.Command, args []string) {
			err := c.Invoke(func(m *migrator.Migrator) error {
				return m.Create(migrationName, migrationExtension)
			})
			if err != nil {
				fmt.Print(err.Error())
			}
		},
	}
	createCmd.Flags().StringVarP(&migrationName, "name", "n", "", "migration name")
	createCmd.Flags().StringVarP(&migrationExtension, "ext", "e", "sql", "file extension")
	createCmd.MarkFlagRequired("name")
	return createCmd
}

func upCmd(c *dig.Container) *cobra.Command {
	upCmd := &cobra.Command{
		Use:     "up",
		Short:   "Execute migrations",
		Run: func(cmd *cobra.Command, args []string) {
			err := c.Invoke(func(m *migrator.Migrator) error {
				return m.Up()
			})
			if err != nil {
				fmt.Print(err.Error())
			}
		},
	}
	return upCmd
}

func downCmd(c *dig.Container) *cobra.Command {
	downCmd := &cobra.Command{
		Use:     "down",
		Short:   "Revert all migrations",
		Run: func(cmd *cobra.Command, args []string) {
			err := c.Invoke(func(m *migrator.Migrator) error {
				return m.Down()
			})
			if err != nil {
				fmt.Print(err.Error())
			}
		},
	}
	return downCmd
}
