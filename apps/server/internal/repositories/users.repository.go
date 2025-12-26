package repositories

import (
	"context"

	"github.com/giridhar-m-a/custom_form_app/internal/db/sqlc"
	"github.com/giridhar-m-a/custom_form_app/internal/utils"
	"github.com/google/uuid"
)

// UserRepository defines the contract for user data access.
type UserRepository interface {
	GetByGoogleID(ctx context.Context, googleID string) (sqlc.GetUserByGoogleIdRow, error)
	GetByEmail(ctx context.Context, email string) (sqlc.GetUserByEmailRow, error)
	GetByID(ctx context.Context, userID uuid.UUID) (sqlc.GetUserByIDRow, error)
	Create(ctx context.Context, params sqlc.CreateUserParams) (sqlc.User, error)
	UpdateUser(ctx context.Context, data sqlc.UpdateUserParams) (sqlc.User, error)
}

// SQLCUserRepository implements UserRepository using sqlc.
type SQLCUserRepository struct {
	q *sqlc.Queries
}

func NewSQLCUserRepository(q *sqlc.Queries) *SQLCUserRepository {
	return &SQLCUserRepository{q: q}
}

func (r *SQLCUserRepository) GetByGoogleID(ctx context.Context, googleID string) (sqlc.GetUserByGoogleIdRow, error) {
	return r.q.GetUserByGoogleId(ctx, utils.ConvertStringToNullString(googleID))
}

func (r *SQLCUserRepository) GetByEmail(ctx context.Context, email string) (sqlc.GetUserByEmailRow, error) {
	return r.q.GetUserByEmail(ctx, email)
}

func (r *SQLCUserRepository) GetByID(ctx context.Context, userID uuid.UUID) (sqlc.GetUserByIDRow, error) {
	return r.q.GetUserByID(ctx, userID)
}

func (r *SQLCUserRepository) Create(ctx context.Context, params sqlc.CreateUserParams) (sqlc.User, error) {
	user, err := r.q.CreateUser(ctx, params)
	if err != nil {
		return sqlc.User{}, err
	}
	return sqlc.User{
		UserID:           user.UserID,
		UserFullName:     user.UserFullName,
		UserEmail:        user.UserEmail,
		UserGoogleID:     user.UserGoogleID,
		UserProfilePicID: user.UserProfilePicID,
		UserCreatedAt:    user.UserCreatedAt,
		UserUpdatedAt:    user.UserUpdatedAt,
	}, nil
}

func (r *SQLCUserRepository) UpdateUser(ctx context.Context, data sqlc.UpdateUserParams) (sqlc.User, error) {
	user, err := r.q.UpdateUser(ctx, data)
	if err != nil {
		return sqlc.User{}, err
	}
	return sqlc.User{
		UserID:           user.UserID,
		UserFullName:     user.UserFullName,
		UserEmail:        user.UserEmail,
		UserProfilePicID: user.UserProfilePicID,
		UserCreatedAt:    user.UserCreatedAt,
		UserUpdatedAt:    user.UserUpdatedAt,
	}, nil
}
