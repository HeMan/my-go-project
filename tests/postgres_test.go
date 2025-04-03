package tests

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"testing"

	"my-go-project/models"
	"my-go-project/routes"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/gavv/httpexpect/v2"
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

type fiberTransport struct {
	app *fiber.App
}

func (ft *fiberTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Use Fiber's app.Test method to handle the request
	resp, err := ft.app.Test(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func TestExampleRouteFunctional(t *testing.T) {
	server := httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: &fiberTransport{app: app}, // Use custom transport
		},
		Reporter: httpexpect.NewRequireReporter(t),
	})
	server.GET("/example").
		Expect().
		Status(200).
		Body().IsEqual("Hello, this is an example route!")
	t.Log("TestExampleRouteFunctional passed")
}

func TestTodoRouteFunctional(t *testing.T) {
	server := httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: &fiberTransport{app: app}, // Use custom transport
		},
		Reporter: httpexpect.NewRequireReporter(t),
	})

	// Define the expected todos using the Todo struct
	expectedTodos := []models.Todo{
		{ID: 1, Subject: "Buy groceries", Completed: false},
		{ID: 2, Subject: "Read a book", Completed: true},
		{ID: 3, Subject: "Write some code", Completed: false},
	}

	server.GET("/todo").
		Expect().
		Status(200).
		JSON().Array().
		IsEqual(expectedTodos)
	t.Log("TestTodoRouteFunctional passed")
}
