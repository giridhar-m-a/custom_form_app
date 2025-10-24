package utils

import (
	"database/sql"

	"github.com/google/uuid"
)

// NullStringToString converts an sql.NullString to a plain string.
// If ns.Valid is true it returns the contained string, otherwise it returns the empty string.
func NullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

// NullUUIDToString converts a uuid.NullUUID to its string representation when valid; otherwise it returns the empty string.
func NullUUIDToString(nu uuid.NullUUID) string {
	if nu.Valid {
		return nu.UUID.String()
	}
	return ""
}

// NullStringToPtr returns a pointer to ns.String when ns.Valid is true, otherwise nil.
// The returned pointer points to a newly allocated string containing the same value as ns.String.
func NullStringToPtr(ns sql.NullString) *string {
	if ns.Valid {
		s := ns.String
		return &s
	}
	return nil
}

// NullUUIDToPtr returns a pointer to the UUID's string when nu.Valid is true, otherwise nil.
// The pointer refers to a newly allocated string containing the UUID's canonical representation.
func NullUUIDToPtr(nu uuid.NullUUID) *string {
	if nu.Valid {
		id := nu.UUID.String()
		return &id
	}
	return nil
}