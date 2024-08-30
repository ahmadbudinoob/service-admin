package main

import (
	"log"

	"saranasistemsolusindo.com/gusen-admin/internal/config"
	"saranasistemsolusindo.com/gusen-admin/internal/router"
)

func main() {
	// Initialize the application configuration
	config.Init()

	// Initialize the router
	e := router.NewRouter()

	// Start the server
	log.Println("Server is running on port 8080")
	log.Fatal(e.Start(":8080"))
}
