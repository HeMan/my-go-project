package tests

import (
	"net/http"
	"testing"

	"my-go-project/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestExampleRoute(t *testing.T) {
	app := fiber.New()
	routes.RegisterExampleRoute(app)

	req, _ := http.NewRequest("GET", "/example", nil) // Updated to use http.NewRequest

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	body := make([]byte, resp.ContentLength)
	resp.Body.Read(body)
	assert.Equal(t, "Hello, this is an example route!", string(body))
}
