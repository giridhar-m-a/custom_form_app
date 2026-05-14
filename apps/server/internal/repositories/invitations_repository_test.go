package repositories

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/giridhar-m-a/custom_form_app/internal/db/sqlc"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeleteInvitationRepo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	queries := sqlc.New(db)
	repo := NewInvitationRepository(queries)

	invitationID := uuid.New()

	mock.ExpectExec("DELETE FROM invitations WHERE invitation_id = ?").
		WithArgs(invitationID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.DeleteInvitation(invitationID, context.Background())

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateSingleInvitationRepo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	queries := sqlc.New(db)
	repo := NewInvitationRepository(queries)

	params := sqlc.CreateInvitationParams{
		FormID:    uuid.New(),
		Email:     "test@example.com",
		Name:      "Test Invitee",
		InvitedBy: uuid.New(),
	}

	rows := sqlmock.NewRows([]string{"invitation_id", "invited_email", "invited_name"}).
		AddRow(uuid.New(), params.Email, params.Name)

	mock.ExpectQuery("INSERT INTO invitations").
		WithArgs(params.FormID, params.Email, params.Name, params.InvitedBy).
		WillReturnRows(rows)

	invitation, err := repo.CreateSingleInvitation(params, context.Background())

	assert.NoError(t, err)
	assert.Equal(t, params.Email, invitation.InvitedEmail)
	assert.NoError(t, mock.ExpectationsWereMet())
}
