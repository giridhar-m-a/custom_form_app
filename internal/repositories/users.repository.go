package repositories

import (
	"context"

	"github.com/giridhar-m-a/custom_form_app/internal/db/sqlc"
	"github.com/google/uuid"
)

// UserRepository defines the contract for user data access.
type UserRepository interface {
	GetByGoogleID(ctx context.Context, googleID string) (sqlc.GetUserByGoogleIdRow, error)
	GetByEmail(ctx context.Context, email string) (sqlc.GetUserByEmailRow, error)
	GetByID(ctx context.Context, userID uuid.UUID) (sqlc.GetUserByIDRow, error)
	Create(ctx context.Context, params sqlc.CreateUserParams) (sqlc.User, error)
}

// SQLCUserRepository implements UserRepository using sqlc.
type SQLCUserRepository struct {
	q *sqlc.Queries
}

// NewSQLCUserRepository creates a new SQLCUserRepository that uses the provided sqlc.Queries for database operations.
func NewSQLCUserRepository(q *sqlc.Queries) *SQLCUserRepository {
	return &SQLCUserRepository{q: q}
}

func (r *SQLCUserRepository) GetByGoogleID(ctx context.Context, googleID string) (sqlc.GetUserByGoogleIdRow, error) {
	return r.q.GetUserByGoogleId(ctx, googleID)
}

func (r *SQLCUserRepository) GetByEmail(ctx context.Context, email string) (sqlc.GetUserByEmailRow, error) {
	return r.q.GetUserByEmail(ctx, email)
}

func (r *SQLCUserRepository) GetByID(ctx context.Context, userID uuid.UUID) (sqlc.GetUserByIDRow, error) {
	return r.q.GetUserByID(ctx, userID)
}

func (r *SQLCUserRepository) Create(ctx context.Context, params sqlc.CreateUserParams) (sqlc.User, error) {
	return r.q.CreateUser(ctx, params)
}