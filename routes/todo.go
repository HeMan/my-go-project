package routes

import "github.com/gofiber/fiber/v2"

func RegisterTodoRoutes(app *fiber.App) {
	app.Get("/todo", func(c *fiber.Ctx) error {
		return c.SendString("Hello, this is an example route!")
	})
}
