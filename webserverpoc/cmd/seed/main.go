package main

import (
	"fmt"
	"os"

	gowebserver "github.com/edwardmccormick/gowebserver"
)

func main() {
	config, err := gowebserver.LoadRuntimeConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	db, err := gowebserver.ConnectToPostgresWithConfig(config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := gowebserver.MigrateSchema(db); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := gowebserver.SeedDemoData(db); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("database seeded")
}
