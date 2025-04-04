package database

import (
	"log"
	"my-go-project/models"
	"my-go-project/utils"

	"gorm.io/gorm"
)

// PopulateDatabase populates the database with fixture data
func PopulateDatabase(db *gorm.DB) {
	todos := []models.Todo{
		{Subject: "Buy groceries", Completed: false},
		{Subject: "Read a book", Completed: true},
		{Subject: "Write some code", Completed: false},
		{Subject: "Due tomorrow", Completed: false, DueDate: utils.ParseDate("2023-10-01")}, // Updated to use utils.ParseDate
		{Subject: "Some notes", Completed: false,
			Notes: []models.Note{
				{Note: "Note 1"},
				{Note: "Note 2"}},
		},
	}

	db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&todos).Error; err != nil {
			return err
		}
		return nil
	})
	log.Println("Database populated with fixture data.")
}
