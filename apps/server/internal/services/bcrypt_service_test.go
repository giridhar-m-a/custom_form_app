package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBcryptService(t *testing.T) {
	service := NewBcryptService()
	password := "password123"

	t.Run("Hash and Compare Success", func(t *testing.T) {
		hash, err := service.HashPassword(password)
		assert.NoError(t, err)
		assert.NotEmpty(t, hash)
		assert.NotEqual(t, password, hash)

		match := service.ComparePassword(hash, password)
		assert.True(t, match)
	})

	t.Run("Compare Failure", func(t *testing.T) {
		hash, _ := service.HashPassword(password)
		match := service.ComparePassword(hash, "wrong_password")
		assert.False(t, match)
	})
}
