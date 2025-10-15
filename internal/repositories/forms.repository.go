package repositories

import (
	"context"

	"github.com/giridhar-m-a/custom_form_app/internal/db/sqlc"
	"github.com/google/uuid"
)

type FormsRepository interface {
	// Define methods for form data access here
	CreateForm(form sqlc.CreateFormParams) (sqlc.CreateFormRow, error)
	UpdateForm(form sqlc.UpdateFormParams) (sqlc.UpdateFormRow, error)
	GetFormByID(id string) (sqlc.Form, error)
	GetFormsList(params sqlc.ListFormsParams) ([]sqlc.Form, error)
	DeleteForm(id string) (sqlc.DeleteFormRow, error)
}

type formsRepository struct {
	q *sqlc.Queries
	ctx context.Context
}

func NewFormsRepository(q *sqlc.Queries, ctx context.Context) FormsRepository {
	return &formsRepository{
		q:   q,
		ctx: ctx,
	}
}

func (r *formsRepository) CreateForm(form sqlc.CreateFormParams) (sqlc.CreateFormRow, error) {
	return r.q.CreateForm(r.ctx, form)
}

func (r *formsRepository) UpdateForm(form sqlc.UpdateFormParams) (sqlc.UpdateFormRow, error) {
	return r.q.UpdateForm(r.ctx, form)
}

func (r *formsRepository) GetFormByID(id string) (sqlc.Form, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return sqlc.Form{}, err
	}
	return r.q.GetFormByID(r.ctx, uid)
}

func (r *formsRepository) GetFormsList(params sqlc.ListFormsParams) ([]sqlc.Form, error) {
	return r.q.ListForms(r.ctx, params)
}

func (r *formsRepository) DeleteForm(id string) (sqlc.DeleteFormRow, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return sqlc.DeleteFormRow{}, err
	}
	return r.q.DeleteForm(r.ctx, uid)
}