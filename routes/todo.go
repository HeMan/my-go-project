package routes

import (
	"log" // Added for logging
	"my-go-project/models"
	"strconv"

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
	app.Get("/todos/:id", func(c *fiber.Ctx) error {
		var todo models.Todo
		id := c.Params("id")
		if err := db.Preload("Notes").First(&todo, id).Error; err != nil {
			log.Printf("Error fetching todo with ID %s: %v", id, err) // Log the error
			return c.Status(404).JSON(fiber.Map{
				"error":   "Todo not found",
				"details": err.Error(),
			})
		}
		return c.JSON(todo)
	})
	app.Delete("/todos/:id", func(c *fiber.Ctx) error {
		var todo models.Todo
		id := c.Params("id")
		if err := db.Delete(&todo, id).Error; err != nil {
			log.Printf("Error deleting todo with ID %s: %v", id, err) // Log the error
			return c.Status(404).JSON(fiber.Map{
				"error":   "Todo not found",
				"details": err.Error(),
			})
		}
		return c.SendStatus(204)
	})
	app.Post("/todos", func(c *fiber.Ctx) error {
		var todo models.Todo
		if err := c.BodyParser(&todo); err != nil {
			log.Printf("Error parsing request body: %v", err) // Log the error
			return c.Status(400).JSON(fiber.Map{
				"error":   "Invalid request body",
				"details": err.Error(),
			})
		}
		if err := db.Create(&todo).Error; err != nil {
			log.Printf("Error creating todo: %v", err) // Log the error
			return c.Status(500).JSON(fiber.Map{
				"error":   "Failed to create todo",
				"details": err.Error(),
			})
		}
		return c.Status(201).JSON(todo) // Return 201 Created on success
	})
	app.Post("/todos/:id/notes", func(c *fiber.Ctx) error {
		var note models.Note
		id := c.Params("id")

		// Parse the request body into the note struct
		if err := c.BodyParser(&note); err != nil {
			log.Printf("Error parsing request body for note: %v", err) // Log the error
			return c.Status(400).JSON(fiber.Map{
				"error":   "Invalid request body",
				"details": err.Error(),
			})
		}

		// Set the TodoID of the note to associate it with the correct todo
		todoID, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			log.Printf("Error converting TodoID to uint: %v", err) // Log the error
			return c.Status(400).JSON(fiber.Map{
				"error":   "Invalid TodoID",
				"details": err.Error(),
			})
		}
		note.TodoID = uint(todoID)

		// Save the note to the database
		if err := db.Create(&note).Error; err != nil {
			log.Printf("Error creating note for todo with ID %s: %v", id, err) // Log the error
			return c.Status(500).JSON(fiber.Map{
				"error":   "Failed to create note",
				"details": err.Error(),
			})
		}

		return c.Status(201).JSON(note) // Return 201 Created on success
	})
	app.Delete("/todos/:todoId/notes/:noteId", func(c *fiber.Ctx) error {
		todoId := c.Params("todoId")
		noteId := c.Params("noteId")

		// Delete the note with the specified ID that belongs to the given TodoID
		if err := db.Where("todo_id = ? AND id = ?", todoId, noteId).Delete(&models.Note{}).Error; err != nil {
			log.Printf("Error deleting note with ID %s for todo with ID %s: %v", noteId, todoId, err) // Log the error
			return c.Status(500).JSON(fiber.Map{
				"error":   "Failed to delete note",
				"details": err.Error(),
			})
		}

		return c.SendStatus(204) // Return 204 No Content on success
	})
	app.Patch("/todos/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var todo models.Todo

		// Find the todo by ID
		if err := db.First(&todo, id).Error; err != nil {
			log.Printf("Error fetching todo with ID %s: %v", id, err) // Log the error
			return c.Status(404).JSON(fiber.Map{
				"error":   "Todo not found",
				"details": err.Error(),
			})
		}

		// Parse the request body and update the todo
		if err := c.BodyParser(&todo); err != nil {
			log.Printf("Error parsing request body for todo with ID %s: %v", id, err) // Log the error
			return c.Status(400).JSON(fiber.Map{
				"error":   "Invalid request body",
				"details": err.Error(),
			})
		}

		// Save the updated todo to the database
		if err := db.Save(&todo).Error; err != nil {
			log.Printf("Error updating todo with ID %s: %v", id, err) // Log the error
			return c.Status(500).JSON(fiber.Map{
				"error":   "Failed to update todo",
				"details": err.Error(),
			})
		}

		return c.JSON(todo) // Return the updated todo
	})

}
