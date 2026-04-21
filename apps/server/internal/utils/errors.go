package utils

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

// Standard API error struct
type ApiError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// Generic function to emit JSON error response
func RespondError(c *gin.Context, status int, message string) {
	c.JSON(status, ApiError{
		Status:  status,
		Message: message,
	})
}

// Main error handler
func HandleError(c *gin.Context, err error) {
	fmt.Printf("Error: doing operation in path: %s with method: %s from ip: %s", c.FullPath(), c.Request.Method, c.ClientIP())
	fmt.Printf("\tError message: %s", err)
	// No error
	if err == nil {
		return
	}

	if strings.Contains(err.Error(), "json: unknown field") {
		RespondError(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if strings.Contains(err.Error(), "user ID not found in context") {
		RespondError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	// 404 — no rows found
	if errors.Is(err, sql.ErrNoRows) {
		RespondError(c, http.StatusNotFound, "resource not found")
		return
	}

	// Check PostgreSQL database error
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {

		switch pgErr.Code {

		// Duplicate key (unique constraint)
		case pgerrcode.UniqueViolation:
			RespondError(c, http.StatusBadRequest, "duplicate value")
			return

		// Foreign key violation
		case pgerrcode.ForeignKeyViolation:
			RespondError(c, http.StatusBadRequest, "invalid reference - foreign key violated")
			return

		// Check constraint
		case pgerrcode.CheckViolation:
			RespondError(c, http.StatusBadRequest, "value does not satisfy validation constraints")
			return
		}

		// Unknown DB error
		RespondError(c, http.StatusInternalServerError, "database error")
		return
	}

	// Generic fallback
	RespondError(c, http.StatusInternalServerError, err.Error())
}
