package main

import (
	"go-api/pkg/app"
	"os"
	"fmt"
)

func main() {
	db := initDB()
	port := os.Getenv("WEB_HTTP_PORT")

	// API init.
	api := &app.App{
		Config: &app.Config{Host: fmt.Sprintf(":%s", port)},
		Route:  registerRoutes(db),
	}

	// Run the API
	api.Run()
}
