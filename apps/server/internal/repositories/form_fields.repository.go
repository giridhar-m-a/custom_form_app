package repositories

import (
	"context"
	"database/sql"

	"github.com/giridhar-m-a/custom_form_app/internal/db/sqlc"
	"github.com/google/uuid"
)

type FormFieldsRepository interface {
	CreateFormField(form sqlc.CreateFormFieldParams, ctx context.Context) (sqlc.CreateFormFieldRow, error)
	UpdateFormField(form sqlc.UpdateFormFieldParams, ctx context.Context) (sqlc.UpdateFormFieldRow, error)
	DeleteFormField(id string, ctx context.Context) (sqlc.DeleteFormFieldRow, error)
	FormFieldRepoWithTx(tx *sql.Tx) FormFieldsRepository
}

type formFieldsRepository struct {
	q *sqlc.Queries
}

func NewFormFieldsRepository(q *sqlc.Queries) FormFieldsRepository {
	return &formFieldsRepository{
		q: q,
	}
}

func (r *formFieldsRepository) CreateFormField(form sqlc.CreateFormFieldParams, ctx context.Context) (sqlc.CreateFormFieldRow, error) {
	return r.q.CreateFormField(ctx, form)
}

func (r *formFieldsRepository) UpdateFormField(form sqlc.UpdateFormFieldParams, ctx context.Context) (sqlc.UpdateFormFieldRow, error) {
	return r.q.UpdateFormField(ctx, form)
}

func (r *formFieldsRepository) DeleteFormField(id string, ctx context.Context) (sqlc.DeleteFormFieldRow, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return sqlc.DeleteFormFieldRow{}, err
	}
	return r.q.DeleteFormField(ctx, uid)
}

func (r *formFieldsRepository) FormFieldRepoWithTx(tx *sql.Tx) FormFieldsRepository {
	return &formFieldsRepository{
		q: r.q.WithTx(tx),
	}
}