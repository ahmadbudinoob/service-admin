package main

import (
	"log"

	"saranasistemsolusindo.com/gusen-admin/internal/config"
	"saranasistemsolusindo.com/gusen-admin/internal/db"
	"saranasistemsolusindo.com/gusen-admin/internal/router"
)

func main() {
	config.Init()

	database, err := db.InitializeDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	e, err := router.InitRouter(database)
	if err != nil {
		log.Fatalf("Failed to initialize router: %v", err)
	}

	log.Println("Server is running on port 8080")
	log.Fatal(e.Start(":8080"))
}
