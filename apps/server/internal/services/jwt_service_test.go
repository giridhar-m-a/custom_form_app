package services

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJWTService(t *testing.T) {
	service := NewJWTService()

	t.Run("Generate and Validate Token", func(t *testing.T) {
		userID := "user-123"
		token, err := service.GenerateToken(userID, time.Hour, "audience")
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		sub, err := service.ValidateToken(token)
		assert.NoError(t, err)
		assert.Equal(t, userID, sub)
	})

	t.Run("Generate and Validate Invitation Token", func(t *testing.T) {
		invitationID := "inv-123"
		formID := "form-123"
		token, err := service.GenerateInvitationToken(invitationID, formID, time.Hour)
		assert.NoError(t, err)

		claims, err := service.ValidateInvitationToken(token)
		assert.NoError(t, err)
		assert.Equal(t, invitationID, *claims.InvitationID)
		assert.Equal(t, formID, claims.FormID)
	})

	t.Run("Validate Token with Bearer prefix", func(t *testing.T) {
		userID := "user-456"
		token, _ := service.GenerateToken(userID, time.Hour, "audience")
		bearerToken := "Bearer " + token

		sub, err := service.ValidateToken(bearerToken)
		assert.NoError(t, err)
		assert.Equal(t, userID, sub)
	})

	t.Run("Invalid Token", func(t *testing.T) {
		_, err := service.ValidateToken("invalid-token")
		assert.Error(t, err)
	})
}
