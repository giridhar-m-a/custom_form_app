package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/giridhar-m-a/custom_form_app/internal/db/sqlc"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetFormByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	queries := sqlc.New(db)
	repo := NewFormsRepository(queries)

	formID := uuid.New()
	now := time.Now()
	
	rows := sqlmock.NewRows([]string{
		"form_id", "form_title", "form_description", "form_status", "form_access", 
		"form_created_at", "form_updated_at", "created_by", "scheduling_id", 
		"scheduled_time", "closing_time", "is_schedule_completed", "is_scheduled", 
		"invitation_schedule_gap", "invitation_schedule_id", "is_deleted",
	}).AddRow(
		formID, "Test Form", "Description", sqlc.FormStatusDraft, sqlc.FormAccessPublic,
		now, now, uuid.New(), uuid.Nil, now, now, false, false, 0, uuid.Nil, false,
	)

	mock.ExpectQuery("SELECT (.+) FROM forms WHERE form_id = (.+) AND is_deleted = FALSE").
		WithArgs(formID).
		WillReturnRows(rows)

	form, err := repo.GetFormByID(formID.String(), context.Background())

	assert.NoError(t, err)
	assert.Equal(t, "Test Form", form.FormTitle)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSoftDeleteForm(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	queries := sqlc.New(db)
	repo := NewFormsRepository(queries)

	formID := uuid.New()

	mock.ExpectExec("UPDATE forms SET is_deleted = TRUE WHERE form_id = (.+)").
		WithArgs(formID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.SoftDeleteForm(formID, context.Background())

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
