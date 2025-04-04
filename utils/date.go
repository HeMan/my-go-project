package utils

import (
	"log"
	"time"
)

// ParseDate parses a date string in the format "2006-01-02" and returns a pointer to the time.Time object.
func ParseDate(dateStr string) *time.Time {
	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		log.Fatalf("Failed to parse date: %s", err)
	}
	return &parsedDate
}
