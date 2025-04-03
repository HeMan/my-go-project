package tests

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestPostgresContainer(t *testing.T) {
	ctx := context.Background()

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
		t.Fatalf("Failed to start PostgreSQL container: %s", err)
	}
	defer postgresContainer.Terminate(ctx)

	// Get the container's host and port
	host, err := postgresContainer.Host(ctx)
	if err != nil {
		t.Fatalf("Failed to get container host: %s", err)
	}
	port, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("Failed to get container port: %s", err)
	}

	// Connect to the PostgreSQL database
	dsn := fmt.Sprintf("postgres://testuser:testpassword@%s:%s/testdb?sslmode=disable", host, port.Port())
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("Failed to connect to PostgreSQL: %s", err)
	}
	defer db.Close()

	// Verify the connection
	err = db.Ping()
	if err != nil {
		t.Fatalf("Failed to ping PostgreSQL: %s", err)
	}
}
