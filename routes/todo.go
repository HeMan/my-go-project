package routes

import (
	"log" // Added for logging
	"my-go-project/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterTodoRoutes(app *fiber.App, db *gorm.DB) {

	app.Get("/todos", func(c *fiber.Ctx) error {
		var todos []models.Todo

		// Attempt to fetch todos with their corresponding notes
		if err := db.Preload("Notes").Find(&todos).Error; err != nil {
			log.Printf("Error fetching todos with notes in transaction: %v", err) // Log the error
			return c.Status(500).JSON(fiber.Map{
				"error":   "Failed to fetch todos with notes",
				"details": err.Error(),
			})
		}
		return c.JSON(todos)
	})
}
