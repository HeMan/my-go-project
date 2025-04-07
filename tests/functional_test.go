package tests

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"my-go-project/database"
	"my-go-project/models"
	"my-go-project/routes"
	"my-go-project/utils"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	pg "gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gavv/httpexpect/v2"
	"github.com/gofiber/fiber/v2"
)

var (
	postgresContainer testcontainers.Container
	db                *gorm.DB
	app               *fiber.App
)

func TestMain(m *testing.M) {
	if os.Getenv("RUN_TESTCONTAINER") == "" {
		fmt.Println("Skipping tests as RUN_TESTCONTAINER is not set")
		return
	}
	ctx := context.Background()
	// Setup PostgreSQL container
	postgresContainer, db = setupPostgresContainer(ctx, nil)
	defer postgresContainer.Terminate(ctx)
	sqlDb, _ := db.DB()
	defer sqlDb.Close()

	// Run migrations
	for _, model := range models.GetRegisteredModels() {
		if err := db.AutoMigrate(model); err != nil {
			fmt.Printf("Failed to migrate model %T: %s\n", model, err)
			os.Exit(m.Run())
		}
	}

	// Populate the database with test data
	database.PopulateDatabase(db)

	// Setup Fiber app and register routes
	app = fiber.New()
	routes.RegisterExampleRoute(app)
	routes.RegisterTodoRoutes(app, db)

	// Run tests
	code := m.Run()
	os.Exit(code)
}

func setupPostgresContainer(ctx context.Context, t *testing.T) (testcontainers.Container, *gorm.DB) {
	dbName := "users"
	dbUser := "user"
	dbPassword := "password"

	postgresContainer, err := postgres.Run(ctx,
		"postgres:17",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		fmt.Printf("failed to start container: %s", err)
		os.Exit(1)
	}

	// Connect to the PostgreSQL database
	dsn, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		if t != nil {
			t.Fatalf("Failed to get connection string: %s", err)
		} else {
			fmt.Printf("Failed to get connection string: %s\n", err)
			os.Exit(1)
		}
	}

	db, err := gorm.Open(pg.Open(dsn), &gorm.Config{})
	if err != nil {
		if t != nil {
			t.Fatalf("Failed to connect to PostgreSQL: %s", err)
		} else {
			fmt.Printf("Failed to connect to PostgreSQL: %s\n", err)
			os.Exit(1)
		}
	}
	sqlDb, _ := db.DB()
	// Verify the connection
	err = sqlDb.Ping()
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
		{Subject: "Buy groceries", Completed: false},
		{Subject: "Read a book", Completed: true},
		{Subject: "Write some code", Completed: false},
		{Subject: "Due tomorrow", Completed: false, DueDate: utils.ParseDate("2023-10-01")},
		{Subject: "Some notes", Completed: false,
			Notes: []models.Note{
				{Note: "Note 1"},
				{Note: "Note 2"}},
		},
	}
	var todos []models.Todo
	result := server.GET("/todos").
		Expect().
		Status(200).
		JSON().Array()

	result.Length().IsEqual(len(expectedTodos))
	result.Decode(&todos)
	for index, todo := range todos {
		assert.Equal(t, todo.Subject, expectedTodos[index].Subject)
		assert.Equal(t, todo.Completed, expectedTodos[index].Completed)
		if expectedTodos[index].DueDate != nil {
			assert.Equal(t, *todo.DueDate, *expectedTodos[index].DueDate)
		}
		if expectedTodos[index].Notes != nil {
			assert.Len(t, todo.Notes, len(expectedTodos[index].Notes))
			for i, note := range todo.Notes {
				assert.Equal(t, note.Note, expectedTodos[index].Notes[i].Note)
			}
		}
	}

	t.Log("TestTodoRouteFunctional passed")
}
