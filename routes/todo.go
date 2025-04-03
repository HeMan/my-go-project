package routes

import (
	"my-go-project/models"

	"github.com/gofiber/fiber/v2"
)

func RegisterTodoRoutes(app *fiber.App) {
	app.Get("/todo", func(c *fiber.Ctx) error {
		todos := []models.Todo{
			{ID: 1, Subject: "Buy groceries", Completed: false},
			{ID: 2, Subject: "Read a book", Completed: true},
			{ID: 3, Subject: "Write some code", Completed: false},
		}
		return c.JSON(todos)
	})
}
