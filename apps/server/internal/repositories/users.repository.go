package repositories

import (
	"context"
	"database/sql"

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
	UpdateUserProfile(ctx context.Context, data sqlc.UpdateUserProfilePicParams) (sqlc.UpdateUserProfilePicRow, error)
	CreateUserProfilePic(ctx context.Context, data sqlc.CreateUserProfilePicParams) (sqlc.CreateUserProfilePicRow, error)
	DeleteUserProfilePic(ctx context.Context, user uuid.UUID) error
	DeleteUser(ctx context.Context, user uuid.UUID) error
	GetProfilePic(ctx context.Context, userID uuid.UUID) (sqlc.UserImage, error)
	GetUserPassword(ctx context.Context, userID uuid.UUID) (sql.NullString, error)
	SoftDeleteUser(ctx context.Context, userID uuid.UUID) error
	CreateTempUser(ctx context.Context, userFullName string) (sqlc.User, error)
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
	return r.q.GetUserByEmail(ctx, utils.ConvertStringToNullString(email))
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
		UserID:        user.UserID,
		UserFullName:  user.UserFullName,
		UserEmail:     user.UserEmail,
		UserGoogleID:  user.UserGoogleID,
		UserCreatedAt: user.UserCreatedAt,
		UserUpdatedAt: user.UserUpdatedAt,
	}, nil
}

func (r *SQLCUserRepository) UpdateUser(ctx context.Context, data sqlc.UpdateUserParams) (sqlc.User, error) {
	user, err := r.q.UpdateUser(ctx, data)
	if err != nil {
		return sqlc.User{}, err
	}
	return sqlc.User{
		UserID:        user.UserID,
		UserFullName:  user.UserFullName,
		UserEmail:     user.UserEmail,
		UserCreatedAt: user.UserCreatedAt,
		UserUpdatedAt: user.UserUpdatedAt,
	}, nil
}

func (r *SQLCUserRepository) UpdateUserProfile(ctx context.Context, data sqlc.UpdateUserProfilePicParams) (sqlc.UpdateUserProfilePicRow, error) {
	return r.q.UpdateUserProfilePic(ctx, data)
}

func (r *SQLCUserRepository) CreateUserProfilePic(ctx context.Context, data sqlc.CreateUserProfilePicParams) (sqlc.CreateUserProfilePicRow, error) {
	return r.q.CreateUserProfilePic(ctx, data)
}

func (r *SQLCUserRepository) DeleteUserProfilePic(ctx context.Context, user uuid.UUID) error {
	return r.q.DeleteUserProfilePic(ctx, user)
}

func (r *SQLCUserRepository) DeleteUser(ctx context.Context, user uuid.UUID) error {
	return r.q.DeleteUser(ctx, user)
}

func (r *SQLCUserRepository) GetProfilePic(ctx context.Context, userID uuid.UUID) (sqlc.UserImage, error) {
	return r.q.GetUserProfilePic(ctx, userID)
}

func (r *SQLCUserRepository) GetUserPassword(ctx context.Context, userID uuid.UUID) (sql.NullString, error) {
	return r.q.GetUserPassword(ctx, userID)
}

func (r *SQLCUserRepository) SoftDeleteUser(ctx context.Context, userID uuid.UUID) error {
	return r.q.SoftDeleteUser(ctx, userID)
}

func (r *SQLCUserRepository) CreateTempUser(ctx context.Context, userFullName string) (sqlc.User, error) {
	return r.q.CreateTempUser(ctx, userFullName)
}