package main

import (
	"log"
	"my-go-project/database"
	"my-go-project/routes"

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

	app := fiber.New()

	// Register routes
	routes.RegisterExampleRoute(app)

	// Start the server
	app.Listen(":8080")
}
