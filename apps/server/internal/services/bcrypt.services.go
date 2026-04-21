package services

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// BcryptService defines the interface for password hashing.
type BcryptService interface {
	HashPassword(password string) (string, error)
	ComparePassword(hashedPassword, password string) bool
}

// concrete implementation
type bcryptService struct{}

// Constructor
func NewBcryptService() BcryptService {
	return &bcryptService{}
}

// HashPassword hashes the given plain password using bcrypt.
func (b *bcryptService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(bytes), nil
}

// ComparePassword compares a bcrypt hash with the provided password.
func (b *bcryptService) ComparePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		fmt.Printf("password does not match: %v\n", err)
		return false
	}
	return true
}
