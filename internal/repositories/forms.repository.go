package repositories

import (
	"context"

	"github.com/giridhar-m-a/custom_form_app/internal/db/sqlc"
	"github.com/google/uuid"
)

type FormsRepository interface {
	// Define methods for form data access here
	CreateForm(form sqlc.CreateFormParams, ctx context.Context) (sqlc.CreateFormRow, error)
	UpdateForm(form sqlc.UpdateFormParams, ctx context.Context) (sqlc.UpdateFormRow, error)
	GetFormByID(id string, ctx context.Context) (sqlc.Form, error)
	GetFormsList(params sqlc.ListFormsParams, ctx context.Context) ([]sqlc.Form, error)
	DeleteForm(id string, ctx context.Context) (sqlc.DeleteFormRow, error)
}

type formsRepository struct {
	q *sqlc.Queries
}

// NewFormsRepository creates a FormsRepository backed by the given *sqlc.Queries.
// The returned repository delegates form persistence operations to the provided Queries instance.
func NewFormsRepository(q *sqlc.Queries) FormsRepository {
	return &formsRepository{
		q: q,
	}
}

func (r *formsRepository) CreateForm(form sqlc.CreateFormParams, ctx context.Context) (sqlc.CreateFormRow, error) {
	return r.q.CreateForm(ctx, form)
}

func (r *formsRepository) UpdateForm(form sqlc.UpdateFormParams, ctx context.Context) (sqlc.UpdateFormRow, error) {
	return r.q.UpdateForm(ctx, form)
}

func (r *formsRepository) GetFormByID(id string, ctx context.Context) (sqlc.Form, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return sqlc.Form{}, err
	}
	return r.q.GetFormByID(ctx, uid)
}

func (r *formsRepository) GetFormsList(params sqlc.ListFormsParams, ctx context.Context) ([]sqlc.Form, error) {
	return r.q.ListForms(ctx, params)
}

func (r *formsRepository) DeleteForm(id string, ctx context.Context) (sqlc.DeleteFormRow, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return sqlc.DeleteFormRow{}, err
	}
	return r.q.DeleteForm(ctx, uid)
}