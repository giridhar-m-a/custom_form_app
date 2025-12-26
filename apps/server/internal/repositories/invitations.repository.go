package repositories

import (
	"context"
	"database/sql"

	"github.com/giridhar-m-a/custom_form_app/internal/db/sqlc"
	"github.com/google/uuid"
)

type InvitationRepository interface {
	CreateInvitation(invitations sqlc.CreateManyInvitationsParams, ctx context.Context) ([]sqlc.CreateManyInvitationsRow, error)
	UpdateInvitationStatus(status sqlc.UpdateInvitationStatusParams, ctx context.Context) (sqlc.UpdateInvitationStatusRow, error)
	DeleteInvitation(invitationID uuid.UUID, ctx context.Context) error
	GetInvitationByFormId(query sqlc.GetInvitationByFormIdParams, ctx context.Context) ([]sqlc.Invitation, error)
	CreateSingleInvitation(invitation sqlc.CreateInvitationParams, ctx context.Context) (sqlc.CreateInvitationRow, error)
	InvitationRepositoryWithTx(tx *sql.Tx) InvitationRepository
}

type invitationRepository struct {
	q *sqlc.Queries
}

func NewInvitationRepository(q *sqlc.Queries) InvitationRepository {
	return &invitationRepository{q: q}
}

func (r *invitationRepository) CreateSingleInvitation(invitation sqlc.CreateInvitationParams, ctx context.Context) (sqlc.CreateInvitationRow, error) {
	return r.q.CreateInvitation(ctx, invitation)
}

func (r *invitationRepository) CreateInvitation(invitations sqlc.CreateManyInvitationsParams, ctx context.Context) ([]sqlc.CreateManyInvitationsRow, error) {
	return r.q.CreateManyInvitations(ctx, invitations)
}

func (r *invitationRepository) UpdateInvitationStatus(status sqlc.UpdateInvitationStatusParams, ctx context.Context) (sqlc.UpdateInvitationStatusRow, error) {
	return r.q.UpdateInvitationStatus(ctx, status)
}

func (r *invitationRepository) DeleteInvitation(invitationID uuid.UUID, ctx context.Context) error {
	return r.q.DeleteInvitation(ctx, invitationID)
}

func (r *invitationRepository) GetInvitationByFormId(query sqlc.GetInvitationByFormIdParams, ctx context.Context) ([]sqlc.Invitation, error) {
	return r.q.GetInvitationByFormId(ctx, query)
}

func (r *invitationRepository) InvitationRepositoryWithTx(tx *sql.Tx) InvitationRepository {
	return &invitationRepository{
		q: r.q.WithTx(tx),
	}
}
