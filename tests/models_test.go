package tests

import (
	"testing"
	"time"

	"my-go-project/models"

	"github.com/stretchr/testify/assert"
)

func TestTodoStruct(t *testing.T) {

	dueDate, _ := time.Parse(time.RFC3339, "2025-10-10T00:00:00Z")

	todo := models.Todo{
		Subject:   "Test Todo",
		DueDate:   &dueDate,
		Completed: false,
		Notes: []models.Note{
			{Note: "First note"},
			{Note: "Second note"},
		},
	}

	assertTodo(t, todo, dueDate)
	assertNotes(t, todo.Notes)
}

func assertTodo(t *testing.T, todo models.Todo, dueDate time.Time) {
	assert.Equal(t, uint(0), todo.ID)
	assert.Equal(t, "Test Todo", todo.Subject)
	assert.Equal(t, dueDate, *todo.DueDate)
	assert.Equal(t, false, todo.Completed)
}

func assertNotes(t *testing.T, notes []models.Note) {
	assert.Len(t, notes, 2)
	assert.Equal(t, uint(0), notes[0].ID)
	assert.Equal(t, "First note", notes[0].Note)
	assert.Equal(t, uint(0), notes[0].TodoID)
	assert.Equal(t, uint(0), notes[1].ID)
	assert.Equal(t, "Second note", notes[1].Note)
	assert.Equal(t, uint(0), notes[1].TodoID)
}
