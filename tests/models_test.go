package tests

import (
	"testing"
	"time"

	"my-go-project/models"

	"github.com/stretchr/testify/assert"
)

func TestTodoStruct(t *testing.T) {
	createdAt, _ := time.Parse(time.RFC3339, "2023-10-01T00:00:00Z")
	dueDate, _ := time.Parse(time.RFC3339, "2025-10-10T00:00:00Z")

	todo := models.Todo{
		ID:        1,
		Subject:   "Test Todo",
		CreatedAt: createdAt,
		DueDate:   &dueDate,
		Completed: false,
		Notes: []models.Note{
			{ID: 1, Note: "First note", TodoID: 1},
			{ID: 2, Note: "Second note", TodoID: 1},
		},
	}

	assertTodo(t, todo, createdAt, dueDate)
	assertNotes(t, todo.Notes)
}

func assertTodo(t *testing.T, todo models.Todo, createdAt, dueDate time.Time) {
	assert.Equal(t, int(1), todo.ID)
	assert.Equal(t, "Test Todo", todo.Subject)
	assert.Equal(t, createdAt, todo.CreatedAt)
	assert.Equal(t, dueDate, *todo.DueDate)
	assert.Equal(t, false, todo.Completed)
}

func assertNotes(t *testing.T, notes []models.Note) {
	assert.Len(t, notes, 2)
	assert.Equal(t, 1, notes[0].ID)
	assert.Equal(t, "First note", notes[0].Note)
	assert.Equal(t, 1, notes[0].TodoID)
	assert.Equal(t, 2, notes[1].ID)
	assert.Equal(t, "Second note", notes[1].Note)
	assert.Equal(t, 1, notes[1].TodoID)
}
