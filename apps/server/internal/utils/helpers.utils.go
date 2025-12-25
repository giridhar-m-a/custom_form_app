package utils

import (
	"database/sql"
	"log"
	"os"
	"strconv"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func GetEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
		log.Printf("Warning: Invalid integer value for %s: '%s', using default: %d",
			key, value, defaultValue)
	}
	return defaultValue
}

func GetEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
		log.Printf("Warning: Invalid boolean value for %s: '%s', using default: %t",
			key, value, defaultValue)
	}
	return defaultValue
}

func ConvertStringToNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{String: "", Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

func ConvertIntToNullInt(i int) sql.NullInt64 {
	if i == 0 {
		return sql.NullInt64{Int64: 0, Valid: false}
	}
	return sql.NullInt64{Int64: int64(i), Valid: true}
}

func ConvertIntToNullInt32(i int) sql.NullInt32 {
	if i == 0 {
		return sql.NullInt32{Int32: 0, Valid: false}
	}
	return sql.NullInt32{Int32: (int32)(i), Valid: true}
}

func ConvertBoolToNullBool(b bool) sql.NullBool {
	if !b {
		return sql.NullBool{Bool: false, Valid: false}
	}
	return sql.NullBool{Bool: true, Valid: true}
}

func ConvertIntToInt32(i int) int32 {
	return int32(i)
}

func ConvertInt32ToInt(i int32) int {
	return int(i)
}

func ConvertStringToUUID(id string) (uuid.UUID, error) {
	if id == "" {
		return uuid.Nil, nil
	}
	return uuid.Parse(id)
}

func ConvertStringToNullUUID(id string) uuid.NullUUID {
	if id == "" {
		return uuid.NullUUID{
			UUID:  uuid.Nil,
			Valid: false,
		}
	}

	parsed, err := ConvertStringToUUID(id)
	if err != nil {
		// If parsing fails, return NullUUID invalid
		return uuid.NullUUID{
			UUID:  uuid.Nil,
			Valid: false,
		}
	}

	return uuid.NullUUID{
		UUID:  parsed,
		Valid: true,
	}
}

func ToPgText(s string) pgtype.Text {
	if s == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: s, Valid: true}
}
