package main

import (
	"os"
	"fmt"
	"go-api/pkg/app"
	"go-api/internal/app/model"
)

func main() {
	dbPassword := os.Getenv("DB_PASSWORD")
	dbUser := os.Getenv("DB_USER")

	if dbPassword == "" || dbUser == "" {
		fmt.Println("Invalid username or password")
		os.Exit(1)
	}

	db, err := app.InitDB(app.NewDBConfig(dbUser, dbPassword))

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if migrateErrs := db.Migrate(&model.Project{}); len(migrateErrs) != 0 {
		for _, migrateErr := range migrateErrs {
			fmt.Println(migrateErr.Error())
		}

		os.Exit(1)
	}

	fmt.Println("success migrate")
	os.Exit(0)
}
