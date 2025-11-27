package repositories

import (
	"context"
	"database/sql"

	"github.com/giridhar-m-a/custom_form_app/internal/db/sqlc"
	"github.com/google/uuid"
)

type FormFieldOptionsRepository interface {
	CreateFieldOption(form sqlc.CreateFieldOptionParams, ctx context.Context) (sqlc.CreateFieldOptionRow, error)
	UpdateFieldOption(form sqlc.UpdateFieldOptionParams, ctx context.Context) (sqlc.UpdateFieldOptionRow, error)
	DeleteFieldOption(id string, ctx context.Context) (sqlc.DeleteFieldOptionRow, error)
	FormFieldOptionsRepoWithTx(tx *sql.Tx) FormFieldOptionsRepository
}

type formFieldOptionsRepository struct {
	q *sqlc.Queries
}

func NewFormFieldOptionsRepository(q *sqlc.Queries) FormFieldOptionsRepository {
	return &formFieldOptionsRepository{
		q: q,
	}
}

func (r *formFieldOptionsRepository) CreateFieldOption(form sqlc.CreateFieldOptionParams, ctx context.Context) (sqlc.CreateFieldOptionRow, error) {
	return r.q.CreateFieldOption(ctx, form)
}

func (r *formFieldOptionsRepository) UpdateFieldOption(form sqlc.UpdateFieldOptionParams, ctx context.Context) (sqlc.UpdateFieldOptionRow, error) {
	return r.q.UpdateFieldOption(ctx, form)
}

func (r *formFieldOptionsRepository) DeleteFieldOption(id string, ctx context.Context) (sqlc.DeleteFieldOptionRow, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return sqlc.DeleteFieldOptionRow{}, err
	}
	return r.q.DeleteFieldOption(ctx, uid)
}

func (r *formFieldOptionsRepository) FormFieldOptionsRepoWithTx(tx *sql.Tx) FormFieldOptionsRepository {
	return &formFieldOptionsRepository{
		q: r.q.WithTx(tx),
	}
}