package main

import "app"

func main() {
	// API init.
	api := &app.App{
		Config: &app.Config{Host: ":8080"},
		Route:  registerRoutes(),
	}

	// Run the API
	api.Run()
}
