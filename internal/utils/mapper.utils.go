package utils

import (
	"database/sql"

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