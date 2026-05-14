package middleware

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/giridhar-m-a/custom_form_app/internal/db"
	"github.com/giridhar-m-a/custom_form_app/internal/db/sqlc"
	"github.com/giridhar-m-a/custom_form_app/internal/services"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	os.Setenv("JWT_SECRET", "test_secret")
	os.Setenv("JWT_ISSUER", "test_issuer")
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func TestAuthMiddleware(t *testing.T) {
	jwtSvc := services.NewJWTService()
	validToken, _ := jwtSvc.GenerateToken("user-123", time.Hour, "test_audience")

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		expectedUserID string
	}{
		{
			name:           "Valid token",
			authHeader:     "Bearer " + validToken,
			expectedStatus: http.StatusOK,
			expectedUserID: "user-123",
		},
		{
			name:           "Invalid token",
			authHeader:     "Bearer invalid-token",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Empty token",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest("GET", "/", nil)
			c.Request.Header.Set("Authorization", tt.authHeader)

			handler := AuthMiddleware()
			handler(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedUserID != "" {
				userID, exists := c.Get("userID")
				assert.True(t, exists)
				assert.Equal(t, tt.expectedUserID, userID)
			}
		})
	}
}

func TestResponseMiddleware(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer sqlDB.Close()

	// Overwrite global db.Queries
	oldQueries := db.Queries
	db.Queries = sqlc.New(sqlDB)
	defer func() { db.Queries = oldQueries }()

	jwtSvc := services.NewJWTService()
	formID := uuid.New()
	invitationID := uuid.New()
	validToken, _ := jwtSvc.GenerateInvitationToken(invitationID.String(), formID.String(), time.Hour)

	formColumns := []string{
		"form_id", "form_title", "form_description", "form_status", "form_access",
		"form_created_at", "form_updated_at", "created_by", "scheduling_id",
		"scheduled_time", "closing_time", "is_schedule_completed", "is_scheduled",
		"invitation_schedule_gap", "invitation_schedule_id", "is_deleted",
	}

	tests := []struct {
		name           string
		authHeader     string
		mockSetup      func()
		expectedStatus int
	}{
		{
			name:       "Valid published form",
			authHeader: "Bearer " + validToken,
			mockSetup: func() {
				rows := sqlmock.NewRows(formColumns).AddRow(
					formID, "Test Form", "Desc", sqlc.FormStatusPublished, sqlc.FormAccessPublic,
					time.Now(), time.Now(), uuid.New(), uuid.Nil, time.Now(), sql.NullTime{Valid: false},
					false, false, 0, uuid.Nil, false,
				)
				mock.ExpectQuery("SELECT (.+) FROM forms WHERE form_id = (.+)").
					WithArgs(formID).
					WillReturnRows(rows)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:       "Form not published",
			authHeader: "Bearer " + validToken,
			mockSetup: func() {
				rows := sqlmock.NewRows(formColumns).AddRow(
					formID, "Test Form", "Desc", sqlc.FormStatusDraft, sqlc.FormAccessPublic,
					time.Now(), time.Now(), uuid.New(), uuid.Nil, time.Now(), sql.NullTime{Valid: false},
					false, false, 0, uuid.Nil, false,
				)
				mock.ExpectQuery("SELECT (.+) FROM forms WHERE form_id = (.+)").
					WithArgs(formID).
					WillReturnRows(rows)
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:       "Form closed",
			authHeader: "Bearer " + validToken,
			mockSetup: func() {
				rows := sqlmock.NewRows(formColumns).AddRow(
					formID, "Test Form", "Desc", sqlc.FormStatusPublished, sqlc.FormAccessPublic,
					time.Now(), time.Now(), uuid.New(), uuid.Nil, time.Now(), sql.NullTime{Time: time.Now().Add(-time.Hour), Valid: true},
					false, false, 0, uuid.Nil, false,
				)
				mock.ExpectQuery("SELECT (.+) FROM forms WHERE form_id = (.+)").
					WithArgs(formID).
					WillReturnRows(rows)
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest("GET", "/", nil)
			c.Request.Header.Set("Authorization", tt.authHeader)

			handler := ResponseMiddleware()
			handler(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
