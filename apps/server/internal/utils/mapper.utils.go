package utils

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

func NullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func NullUUIDToString(nu uuid.NullUUID) string {
	if nu.Valid {
		return nu.UUID.String()
	}
	return ""
}

func NullStringToPtr(ns sql.NullString) *string {
	if ns.Valid {
		s := ns.String
		return &s
	}
	return nil
}

func NullUUIDToPtr(nu uuid.NullUUID) *string {
	if nu.Valid {
		id := nu.UUID.String()
		return &id
	}
	return nil
}

func BoolToNullBool(b *bool) sql.NullBool {
	if b == nil {
		return sql.NullBool{Valid: false}
	}
	return sql.NullBool{Bool: *b, Valid: true}
}

func TimeToNullTime(t time.Time) sql.NullTime {
	if t.IsZero() {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: t, Valid: true}
}

func Int64ToNullInt64(i *int64) sql.NullInt64 {
	if i == nil {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Int64: *i, Valid: true}
}
