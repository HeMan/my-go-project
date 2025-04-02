package tests

import (
	"testing"

	"my-go-project/models"

	"github.com/stretchr/testify/assert"
)

func TestUserStruct(t *testing.T) {
	user := models.User{
		ID:       1,
		Name:     "John Doe",
		Email:    "john.doe@example.com",
		Password: "securepassword",
	}

	assert.Equal(t, uint(1), user.ID)
	assert.Equal(t, "John Doe", user.Name)
	assert.Equal(t, "john.doe@example.com", user.Email)
	assert.Equal(t, "securepassword", user.Password)
}
