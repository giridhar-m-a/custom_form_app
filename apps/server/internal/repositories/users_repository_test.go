package repositories

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/giridhar-m-a/custom_form_app/internal/db/sqlc"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	queries := sqlc.New(db)
	repo := NewSQLCUserRepository(queries)

	email := "test@example.com"
	userID := uuid.New()

	rows := sqlmock.NewRows([]string{"user_id", "user_full_name", "user_email", "user_created_at", "user_updated_at", "user_password", "file_name"}).
		AddRow(userID, "Test User", email, nil, nil, "password", nil)

	mock.ExpectQuery("SELECT (.+) FROM users u").
		WithArgs(sql.NullString{String: email, Valid: true}).
		WillReturnRows(rows)

	user, err := repo.GetByEmail(context.Background(), email)

	assert.NoError(t, err)
	assert.Equal(t, userID, user.UserID)
	assert.Equal(t, email, user.UserEmail.String)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateUserRepo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	queries := sqlc.New(db)
	repo := NewSQLCUserRepository(queries)

	params := sqlc.CreateUserParams{
		UserFullName: "New User",
		UserEmail:    sql.NullString{String: "new@example.com", Valid: true},
	}

	userID := uuid.New()
	rows := sqlmock.NewRows([]string{"user_id", "user_full_name", "user_email", "user_google_id", "user_created_at", "user_updated_at"}).
		AddRow(userID, params.UserFullName, params.UserEmail.String, nil, nil, nil)

	mock.ExpectQuery("INSERT INTO users").
		WithArgs(params.UserFullName, params.UserEmail, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	user, err := repo.Create(context.Background(), params)

	assert.NoError(t, err)
	assert.Equal(t, userID, user.UserID)
	assert.NoError(t, mock.ExpectationsWereMet())
}
func TestSoftDeleteUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	queries := sqlc.New(db)
	repo := NewSQLCUserRepository(queries)

	userID := uuid.New()

	mock.ExpectExec("UPDATE users SET is_deleted = TRUE WHERE user_id = (.+)").
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.SoftDeleteUser(context.Background(), userID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
