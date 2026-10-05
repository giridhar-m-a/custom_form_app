package repositories

import (
	"context"

	"github.com/giridhar-m-a/custom_form_app/internal/db/sqlc"
	"github.com/google/uuid"
)

type DashboardRepository interface {
	GetTotalForms(ctx context.Context, userID uuid.NullUUID) (int64, error)
	GetTotalSubmissions(ctx context.Context, userID uuid.NullUUID) (int64, error)
	GetTotalActiveForms(ctx context.Context, userID uuid.NullUUID) (int64, error)
	GetTotalClosedForms(ctx context.Context, userID uuid.NullUUID) (int64, error)
	GetTotalInvitations(ctx context.Context, userID uuid.NullUUID) (int64, error)
	GetFormSubmissionsByMonth(ctx context.Context, userID uuid.NullUUID) ([]sqlc.GetFormSubmissionsByMonthRow, error)
}

type dashboardRepository struct {
	q *sqlc.Queries
}

func NewDashboardRepository(q *sqlc.Queries) DashboardRepository {
	return &dashboardRepository{
		q: q,
	}
}

func (r *dashboardRepository) GetTotalForms(ctx context.Context, userID uuid.NullUUID) (int64, error) {
	return r.q.GetTotalForms(ctx, userID)
}

func (r *dashboardRepository) GetTotalSubmissions(ctx context.Context, userID uuid.NullUUID) (int64, error) {
	return r.q.GetTotalSubmissions(ctx, userID)
}

func (r *dashboardRepository) GetTotalActiveForms(ctx context.Context, userID uuid.NullUUID) (int64, error) {
	return r.q.GetTotalActiveForms(ctx, userID)
}

func (r *dashboardRepository) GetTotalClosedForms(ctx context.Context, userID uuid.NullUUID) (int64, error) {
	return r.q.GetTotalClosedForms(ctx, userID)
}

func (r *dashboardRepository) GetTotalInvitations(ctx context.Context, userID uuid.NullUUID) (int64, error) {
	return r.q.GetTotalInvitations(ctx, userID)
}

func (r *dashboardRepository) GetFormSubmissionsByMonth(ctx context.Context, userID uuid.NullUUID) ([]sqlc.GetFormSubmissionsByMonthRow, error) {
	return r.q.GetFormSubmissionsByMonth(ctx, userID)
}
	