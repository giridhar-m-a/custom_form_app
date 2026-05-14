package repositories

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/giridhar-m-a/custom_form_app/internal/db/sqlc"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetTotalFormsRepo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	queries := sqlc.New(db)
	repo := NewDashboardRepository(queries)

	userID := uuid.NullUUID{UUID: uuid.New(), Valid: true}

	rows := sqlmock.NewRows([]string{"count"}).AddRow(int64(5))

	mock.ExpectQuery("SELECT COUNT(.+) FROM forms").
		WithArgs(userID).
		WillReturnRows(rows)

	count, err := repo.GetTotalForms(context.Background(), userID)

	assert.NoError(t, err)
	assert.Equal(t, int64(5), count)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTotalSubmissionsRepo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	queries := sqlc.New(db)
	repo := NewDashboardRepository(queries)

	userID := uuid.NullUUID{UUID: uuid.New(), Valid: true}

	rows := sqlmock.NewRows([]string{"count"}).AddRow(int64(20))

	mock.ExpectQuery("SELECT COUNT(.+) FROM form_submissions").
		WithArgs(userID).
		WillReturnRows(rows)

	count, err := repo.GetTotalSubmissions(context.Background(), userID)

	assert.NoError(t, err)
	assert.Equal(t, int64(20), count)
	assert.NoError(t, mock.ExpectationsWereMet())
}
