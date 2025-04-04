package models

import (
	"time"

	"gorm.io/gorm"
)

func init() {
	RegisterModel(&Todo{})
	RegisterModel(&Note{})
}

// Todo represents a task with a summary, dates, and completion status.
type Todo struct {
	gorm.Model
	Subject   string     `gorm:"size:255;not null" json:"subject"`
	DueDate   *time.Time `gorm:"type:timestamp" json:"due_date,omitempty"` // Pointer to allow empty value
	Completed bool       `gorm:"default:false" json:"completed"`
	Notes     []Note     `gorm:"foreignKey:TodoID;constraint:OnDelete:CASCADE;" json:"notes"` // One-to-many relationship
}

// Note represents a note associated with a Todo.
type Note struct {
	gorm.Model
	Note   string `gorm:"size:500;not null" json:"note"`
	TodoID uint   `gorm:"not null" json:"todo_id"` // Foreign key to Todo
}
