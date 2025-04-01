package routes

import "github.com/gofiber/fiber/v2"

func RegisterExampleRoute(app *fiber.App) {
	app.Get("/example", func(c *fiber.Ctx) error {
		return c.SendString("Hello, this is an example route!")
	})
}
