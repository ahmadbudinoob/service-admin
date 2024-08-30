package main

import (
	"log"

	"saranasistemsolusindo.com/gusen-admin/internal/config"
	"saranasistemsolusindo.com/gusen-admin/internal/db"
	"saranasistemsolusindo.com/gusen-admin/internal/router"
)

func main() {
	// Initialize the application configuration
	config.Init()

	database, err := db.InitializeDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Initialize the router
	e, err := router.InitRouter(database)
	if err != nil {
		log.Fatalf("Failed to initialize router: %v", err)
	}

	// Start the server
	log.Println("Server is running on port 8080")
	log.Fatal(e.Start(":8080"))
}
