package tests

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"os"
	"testing"

	"my-go-project/routes"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/gofiber/fiber/v2"
)

var (
	postgresContainer testcontainers.Container
	db                *sql.DB
	app               *fiber.App
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	// Setup PostgreSQL container
	postgresContainer, db = setupPostgresContainer(ctx, nil)
	defer postgresContainer.Terminate(ctx)
	defer db.Close()

	// Setup Fiber app and register routes
	app = fiber.New()
	routes.RegisterExampleRoute(app)
	routes.RegisterTodoRoutes(app)

	// Run tests
	code := m.Run()
	os.Exit(code)
}

func setupPostgresContainer(ctx context.Context, t *testing.T) (testcontainers.Container, *sql.DB) {
	// Create a PostgreSQL container
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpassword",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}
	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		if t != nil {
			t.Fatalf("Failed to start PostgreSQL container: %s", err)
		} else {
			fmt.Printf("Failed to start PostgreSQL container: %s\n", err)
			os.Exit(1)
		}
	}

	// Get the container's host and port
	host, err := postgresContainer.Host(ctx)
	if err != nil {
		if t != nil {
			t.Fatalf("Failed to get container host: %s", err)
		} else {
			fmt.Printf("Failed to get container host: %s\n", err)
			os.Exit(1)
		}
	}
	port, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		if t != nil {
			t.Fatalf("Failed to get container port: %s", err)
		} else {
			fmt.Printf("Failed to get container port: %s\n", err)
			os.Exit(1)
		}
	}

	// Connect to the PostgreSQL database
	dsn := fmt.Sprintf("postgres://testuser:testpassword@%s:%s/testdb?sslmode=disable", host, port.Port())
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		if t != nil {
			t.Fatalf("Failed to connect to PostgreSQL: %s", err)
		} else {
			fmt.Printf("Failed to connect to PostgreSQL: %s\n", err)
			os.Exit(1)
		}
	}

	// Verify the connection
	err = db.Ping()
	if err != nil {
		if t != nil {
			t.Fatalf("Failed to ping PostgreSQL: %s", err)
		} else {
			fmt.Printf("Failed to ping PostgreSQL: %s\n", err)
			os.Exit(1)
		}
	}

	return postgresContainer, db
}

func TestExampleRouteFunctional(t *testing.T) {
	http_req := httptest.NewRequest("GET", "/example", nil)
	resp, err := app.Test(http_req)
	if err != nil {
		t.Fatalf("Failed to test Fiber app: %s", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Expected status code 200, got %d", resp.StatusCode)
	}
	respBody, _ := io.ReadAll(resp.Body)
	assert.Equal(t, string(respBody), "Hello, this is an example route!")
}

func TestTodoRouteFunctional(t *testing.T) {
	http_req := httptest.NewRequest("GET", "/todo", nil)
	resp, err := app.Test(http_req)
	if err != nil {
		t.Fatalf("Failed to test Fiber app: %s", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Expected status code 200, got %d", resp.StatusCode)
	}
	respBody, _ := io.ReadAll(resp.Body)

	// Define the expected and actual structs
	type Todo struct {
		ID   int    `json:"id"`
		Task string `json:"task"`
	}
	var actualTodos []Todo
	err = json.Unmarshal(respBody, &actualTodos)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %s", err)
	}

	expectedTodos := []Todo{
		{ID: 1, Task: "Buy groceries"},
		{ID: 2, Task: "Clean the house"},
	}

	assert.Equal(t, expectedTodos, actualTodos, "Response body does not match expected todos")
}
