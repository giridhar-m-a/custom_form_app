package utils

import (
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetEnv(t *testing.T) {
	t.Run("Existing env", func(t *testing.T) {
		t.Setenv("TEST_ENV", "value")
		assert.Equal(t, "value", GetEnv("TEST_ENV", "default"))
	})
	t.Run("Default value", func(t *testing.T) {
		assert.Equal(t, "default", GetEnv("NON_EXISTENT", "default"))
	})
}

func TestGetEnvAsInt(t *testing.T) {
	t.Run("Valid int", func(t *testing.T) {
		t.Setenv("TEST_INT", "123")
		assert.Equal(t, 123, GetEnvAsInt("TEST_INT", 0))
	})
	t.Run("Invalid int", func(t *testing.T) {
		t.Setenv("TEST_INT", "abc")
		assert.Equal(t, 0, GetEnvAsInt("TEST_INT", 0))
	})
}

func TestConvertStringToNullString(t *testing.T) {
	assert.True(t, ConvertStringToNullString("test").Valid)
	assert.Equal(t, "test", ConvertStringToNullString("test").String)
	assert.False(t, ConvertStringToNullString("").Valid)
}

func TestNullUUIDToStringOrEmpty(t *testing.T) {
	u := uuid.New()
	assert.Equal(t, u.String(), NullUUIDToStringOrEmpty(uuid.NullUUID{UUID: u, Valid: true}))
	assert.Equal(t, "", NullUUIDToStringOrEmpty(uuid.NullUUID{Valid: false}))
}

func TestNullTimeToString(t *testing.T) {
	now := time.Now()
	assert.Equal(t, now.Format(time.RFC3339), NullTimeToString(sql.NullTime{Time: now, Valid: true}))
	assert.Equal(t, "", NullTimeToString(sql.NullTime{Valid: false}))
}

func TestMapperFunctions(t *testing.T) {
	t.Run("NullStringToString", func(t *testing.T) {
		assert.Equal(t, "test", NullStringToString(sql.NullString{String: "test", Valid: true}))
		assert.Equal(t, "", NullStringToString(sql.NullString{Valid: false}))
	})

	t.Run("NullStringToPtr", func(t *testing.T) {
		res := NullStringToPtr(sql.NullString{String: "test", Valid: true})
		assert.NotNil(t, res)
		assert.Equal(t, "test", *res)
		assert.Nil(t, NullStringToPtr(sql.NullString{Valid: false}))
	})
}
