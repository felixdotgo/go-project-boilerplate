package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/0x46656C6978/go-project-boilerplate/pkg/boilerplate"
)

func main() {
	var modulePath string
	flag.StringVar(&modulePath, "module", "", "Target Go module path (defaults to origin remote)")
	flag.Parse()

	result, err := boilerplate.RenameModule(boilerplate.RenameModuleOptions{
		RootDir:      ".",
		TargetModule: modulePath,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("Current module: %s\n", result.OldModule)
	fmt.Printf("Target module: %s\n", result.NewModule)
	if len(result.ChangedFiles) == 0 {
		fmt.Println("No files changed.")
		return
	}

	fmt.Println("Updated files:")
	for _, file := range result.ChangedFiles {
		fmt.Printf("- %s\n", file)
	}
}
