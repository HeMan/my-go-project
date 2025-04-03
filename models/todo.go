package models

import (
	"time"
)

func init() {
	RegisterModel(&Todo{})
	RegisterModel(&Note{})
}

// Todo represents a task with a summary, dates, and completion status.
type Todo struct {
	ID        int        `gorm:"primaryKey" json:"id"`
	Subject   string     `gorm:"size:255;not null" json:"subject"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_date"`
	DueDate   *time.Time `gorm:"type:timestamp" json:"due_date,omitempty"` // Pointer to allow empty value
	Completed bool       `gorm:"default:false" json:"completed"`
	Notes     []Note     `gorm:"foreignKey:TodoID;constraint:OnDelete:CASCADE;" json:"notes"` // One-to-many relationship
}

// Note represents a note associated with a Todo.
type Note struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	Note      string    `gorm:"size:500;not null" json:"note"`
	TodoID    int       `gorm:"not null" json:"todo_id"`               // Foreign key to Todo
	Todo      Todo      `gorm:"constraint:OnDelete:CASCADE;" json:"-"` // Many-to-one relationship
}
