package repositories

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/giridhar-m-a/custom_form_app/internal/db/sqlc"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateSubmission(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	queries := sqlc.New(db)
	repo := NewResponseRepository(queries)

	formID := uuid.New()
	respondentID := uuid.New()
	arg := sqlc.CreateSubmissionParams{
		FormID:       formID,
		RespondentID: uuid.NullUUID{UUID: respondentID, Valid: true},
	}

	now := time.Now()
	rows := sqlmock.NewRows([]string{"submission_id", "form_id", "submitted_at", "respondent_id"}).
		AddRow(uuid.New(), formID, now, respondentID)

	mock.ExpectQuery("INSERT INTO form_submissions").
		WithArgs(arg.FormID, arg.RespondentID).
		WillReturnRows(rows)

	submission, err := repo.CreateSubmission(context.Background(), arg)

	assert.NoError(t, err)
	assert.Equal(t, formID, submission.FormID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetSubmissionCount(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	queries := sqlc.New(db)
	repo := NewResponseRepository(queries)

	formID := uuid.New()
	arg := sqlc.GetSubmissionCountParams{
		FormID: formID,
		Search: sql.NullString{String: "test", Valid: true},
	}

	rows := sqlmock.NewRows([]string{"count"}).AddRow(int64(10))

	mock.ExpectQuery("WITH filtered AS").
		WithArgs(arg.FormID, arg.Search).
		WillReturnRows(rows)

	count, err := repo.GetSubmissionCount(context.Background(), arg)

	assert.NoError(t, err)
	assert.Equal(t, int64(10), count)
	assert.NoError(t, mock.ExpectationsWereMet())
}
