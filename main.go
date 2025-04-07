package main

import (
	"log"
	"my-go-project/database"
	"my-go-project/routes"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	// Initialize the database
	database.Init()

	// Check if the "populate" argument is present
	if len(os.Args) > 1 && os.Args[1] == "populate" {
		database.PopulateDatabase(database.DB)
	}

	app := fiber.New()

	// Register routes
	routes.RegisterExampleRoute(app)
	routes.RegisterTodoRoutes(app, database.DB)
	routes.RegisterSwaggerRoute(app)

	// Debug: Print all registered routes
	for _, route := range app.Stack() {
		for _, r := range route {
			log.Printf("Route registered: %s %s", r.Method, r.Path)
		}
	}

	// Start the server
	app.Listen(":8080")
}
