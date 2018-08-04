package main

import (
	"go-api/pkg/app"
	"os"
	"fmt"
)

func initDB() *app.DB {
	dbPassword := os.Getenv("DB_PASSWORD")
	dbUser := os.Getenv("DB_USER")

	if dbPassword == "" || dbUser == "" {
		fmt.Println("Invalid username or password")
		os.Exit(1)
	}

	db, _ := app.InitDB(app.NewDBConfig(dbUser, dbPassword))

	return db
}
