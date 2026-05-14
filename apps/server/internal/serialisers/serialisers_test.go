package serializers

import (
	"database/sql"
	"testing"
	"time"

	"github.com/giridhar-m-a/custom_form_app/internal/db/sqlc"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMapUser(t *testing.T) {
	id := uuid.New()
	now := time.Now()
	u := sqlc.User{
		UserID:       id,
		UserEmail:    sql.NullString{String: "test@example.com", Valid: true},
		UserFullName: "Test User",
		UserCreatedAt: sql.NullTime{Time: now, Valid: true},
		UserUpdatedAt: sql.NullTime{Time: now, Valid: true},
	}

	res := MapUser(u)
	assert.Equal(t, id.String(), res.UserID)
	assert.Equal(t, "test@example.com", res.UserEmail)
	assert.Equal(t, "Test User", res.UserFullName)
	assert.Equal(t, now, res.UserCreatedAt)
}

func TestMapGetUserByEmailRow(t *testing.T) {
	id := uuid.New()
	now := time.Now()
	u := sqlc.GetUserByEmailRow{
		UserID:       id,
		UserEmail:    sql.NullString{String: "test@example.com", Valid: true},
		UserFullName: "Test User",
		UserCreatedAt: sql.NullTime{Time: now, Valid: true},
		UserUpdatedAt: sql.NullTime{Time: now, Valid: true},
	}

	res := MapGetUserByEmailRow(u)
	assert.Equal(t, id.String(), res.UserID)
	assert.Equal(t, "test@example.com", res.UserEmail)
}
