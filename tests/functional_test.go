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
	client            *http.Client
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
	client = &http.Client{
		Transport: &fiberTransport{app: app}, // Use custom transport
	}

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

func TestAllTodosRouteFunctional(t *testing.T) {
	server := httpexpect.WithConfig(httpexpect.Config{
		Client:   client,
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

func TestSingleTodoRouteFunctional(t *testing.T) {
	server := httpexpect.WithConfig(httpexpect.Config{
		Client:   client,
		Reporter: httpexpect.NewRequireReporter(t),
	})

	// Define the expected todo using the Todo struct
	expectedTodo := models.Todo{
		Subject:   "Buy groceries",
		Completed: false,
	}

	var todo models.Todo
	result := server.GET("/todos/1").
		Expect().
		Status(200).
		JSON().Object()

	result.Decode(&todo)
	assert.Equal(t, todo.Subject, expectedTodo.Subject)
	assert.Equal(t, todo.Completed, expectedTodo.Completed)

	t.Log("TestSingleTodoRouteFunctional passed")
}

func TestDeleteTodoRouteFunctional(t *testing.T) {
	server := httpexpect.WithConfig(httpexpect.Config{
		Client:   client,
		Reporter: httpexpect.NewRequireReporter(t),
	})

	// Ensure the todo exists before deletion
	server.GET("/todos/3").
		Expect().
		Status(200)

	// Perform the delete operation
	server.DELETE("/todos/3").
		Expect().
		Status(204)

	// Verify the todo no longer exists
	server.GET("/todos/3").
		Expect().
		Status(404)

	t.Log("TestDeleteTodoRouteFunctional passed")
}
func TestCreateTodoRouteFunctional(t *testing.T) {
	server := httpexpect.WithConfig(httpexpect.Config{
		Client:   client,
		Reporter: httpexpect.NewRequireReporter(t),
	})

	// Define the new todo to be created
	newTodo := models.Todo{
		Subject:   "New Task",
		Completed: false,
	}

	// Perform the POST operation
	response := server.POST("/todos").
		WithJSON(newTodo).
		Expect().
		Status(201).
		JSON().Object()

	// Verify the response contains the created todo
	response.Value("subject").IsEqual(newTodo.Subject)
	response.Value("completed").IsEqual(newTodo.Completed)

	// Verify the todo exists in the database
	var createdTodo models.Todo
	response.Value("ID").Number().Gt(0) // Ensure ID is valid
	todoID := int(response.Value("ID").Raw().(float64))

	server.GET(fmt.Sprintf("/todos/%d", todoID)).
		Expect().
		Status(200).
		JSON().Object().
		Decode(&createdTodo)

	assert.Equal(t, createdTodo.Subject, newTodo.Subject)
	assert.Equal(t, createdTodo.Completed, newTodo.Completed)

	t.Log("TestCreateTodoRouteFunctional passed")
}

func TestAddNoteToTodoFunctional(t *testing.T) {
	server := httpexpect.WithConfig(httpexpect.Config{
		Client:   client,
		Reporter: httpexpect.NewRequireReporter(t),
	})

	// Define the new note to be added
	newNote := models.Note{
		Note: "This is a new note",
	}

	// Perform the POST operation to add the note
	response := server.POST("/todos/2/notes").
		WithJSON(newNote).
		Expect().
		Status(201).
		JSON().Object()

	// Verify the response contains the created note
	response.Value("note").IsEqual(newNote.Note)

	// Verify the note exists in the database for todo/2
	var updatedTodo models.Todo
	server.GET("/todos/2").
		Expect().
		Status(200).
		JSON().Object().
		Decode(&updatedTodo)

	assert.NotNil(t, updatedTodo.Notes)
	assert.GreaterOrEqual(t, len(updatedTodo.Notes), 1)
	assert.Equal(t, updatedTodo.Notes[len(updatedTodo.Notes)-1].Note, newNote.Note)

	t.Log("TestAddNoteToTodoFunctional passed")
}

func TestRemoveNoteFromTodoFunctional(t *testing.T) {
	server := httpexpect.WithConfig(httpexpect.Config{
		Client:   client,
		Reporter: httpexpect.NewRequireReporter(t),
	})

	// Ensure the note exists before deletion
	var todoBefore models.Todo
	server.GET("/todos/2").
		Expect().
		Status(200).
		JSON().Object().
		Decode(&todoBefore)

	assert.NotNil(t, todoBefore.Notes)
	assert.GreaterOrEqual(t, len(todoBefore.Notes), 1)

	noteID := todoBefore.Notes[0].ID

	// Perform the DELETE operation to remove the note
	server.DELETE(fmt.Sprintf("/todos/2/notes/%d", noteID)).
		Expect().
		Status(204)

	// Verify the note no longer exists in the database for todo/2
	var todoAfter models.Todo
	server.GET("/todos/2").
		Expect().
		Status(200).
		JSON().Object().
		Decode(&todoAfter)

	assert.NotNil(t, todoAfter.Notes)
	for _, note := range todoAfter.Notes {
		assert.NotEqual(t, note.ID, noteID)
	}

	t.Log("TestRemoveNoteFromTodoFunctional passed")
}

func TestAddTodoWithDueDateFunctional(t *testing.T) {
	server := httpexpect.WithConfig(httpexpect.Config{
		Client:   client,
		Reporter: httpexpect.NewRequireReporter(t),
	})

	// Define the new todo with a due date
	newTodo := models.Todo{
		Subject:   "Task with due date",
		Completed: false,
		DueDate:   utils.ParseDate("2023-12-31"),
	}

	// Perform the POST operation
	response := server.POST("/todos").
		WithJSON(newTodo).
		Expect().
		Status(201).
		JSON().Object()

	// Verify the response contains the created todo
	response.Value("subject").IsEqual(newTodo.Subject)
	response.Value("completed").IsEqual(newTodo.Completed)
	response.Value("due_date").IsEqual("2023-12-31T00:00:00Z")

	// Verify the todo exists in the database
	var createdTodo models.Todo
	response.Value("ID").Number().Gt(0) // Ensure ID is valid
	todoID := int(response.Value("ID").Raw().(float64))

	server.GET(fmt.Sprintf("/todos/%d", todoID)).
		Expect().
		Status(200).
		JSON().Object().
		Decode(&createdTodo)

	assert.Equal(t, createdTodo.Subject, newTodo.Subject)
	assert.Equal(t, createdTodo.Completed, newTodo.Completed)
	assert.NotNil(t, createdTodo.DueDate)
	assert.Equal(t, *createdTodo.DueDate, *newTodo.DueDate)

	t.Log("TestAddTodoWithDueDateFunctional passed")
}

func TestMarkTodoAsCompletedFunctional(t *testing.T) {
	server := httpexpect.WithConfig(httpexpect.Config{
		Client:   client,
		Reporter: httpexpect.NewRequireReporter(t),
	})

	// Ensure the todo exists and is not completed
	var todoBefore models.Todo
	server.GET("/todos/4").
		Expect().
		Status(200).
		JSON().Object().
		Decode(&todoBefore)

	assert.False(t, todoBefore.Completed)

	// Perform the PATCH operation to mark the todo as completed
	server.PATCH("/todos/4").
		WithJSON(map[string]interface{}{"completed": true}).
		Expect().
		Status(200)

	// Verify the todo is now marked as completed
	var todoAfter models.Todo
	server.GET("/todos/4").
		Expect().
		Status(200).
		JSON().Object().
		Decode(&todoAfter)

	assert.True(t, todoAfter.Completed)

	t.Log("TestMarkTodoAsCompletedFunctional passed")
}
