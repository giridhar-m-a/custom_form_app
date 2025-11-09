// internal/services/user_service.go
package services

import (
	"context"
	"encoding/json"

	"github.com/giridhar-m-a/custom_form_app/internal/cache"
	"github.com/giridhar-m-a/custom_form_app/internal/db/sqlc"
	"github.com/giridhar-m-a/custom_form_app/internal/repositories"
	"github.com/giridhar-m-a/custom_form_app/internal/utils"
	"github.com/google/uuid"
)

type UserService interface {
	GetUserDetailsById(ctx context.Context, userID string) (sqlc.GetUserByIDRow, error)
	GetUserDetailsByEmail(ctx context.Context, email string) (sqlc.GetUserByEmailRow, error)
	GetUserDetailsByGoogleId(ctx context.Context, googleID string) (sqlc.GetUserByGoogleIdRow, error)
	CreateUser(ctx context.Context, data map[string]any) (sqlc.User, error)
}

type userService struct {
	repo repositories.UserRepository
}

func UserServiceProvider(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(ctx context.Context, data map[string]any) (sqlc.User, error) {
	newUser, err := s.repo.Create(ctx, sqlc.CreateUserParams{
		UserFullName: data["name"].(string),
		UserEmail:    data["email"].(string),
		UserGoogleID: utils.ConvertStringToNullString(data["id"].(string)),
		UserPassword: utils.ConvertStringToNullString(data["password"].(string)),
	})
	return newUser, err
}

func (s *userService) GetUserDetailsById(ctx context.Context, userID string) (sqlc.GetUserByIDRow, error) {
	user, err := uuid.Parse(userID)
	if err != nil {
		return sqlc.GetUserByIDRow{}, err
	}
	return s.repo.GetByID(ctx, user)
}

func (s *userService) GetUserDetailsByEmail(ctx context.Context, email string) (sqlc.GetUserByEmailRow, error) {
	return s.repo.GetByEmail(ctx, email)
}

func (s *userService) GetUserDetailsByGoogleId(ctx context.Context, googleID string) (sqlc.GetUserByGoogleIdRow, error) {
	key := "user:google_id:" + googleID

	// 1. Try cache
	cachedUser, err := cache.Get(ctx, key)
	if err == nil && cachedUser != "" {
		var user sqlc.GetUserByGoogleIdRow
		if err := json.Unmarshal([]byte(cachedUser), &user); err == nil {
			// ✅ Cache hit, return immediately
			return user, nil
		}
	}

	// 2. Fallback to DB
	user, err := s.repo.GetByGoogleID(ctx, googleID)
	if err != nil {
		return sqlc.GetUserByGoogleIdRow{}, err
	}

	// 3. Save to cache (async or ignore error if you like)
	userJSON, _ := json.Marshal(user)
	_ = cache.Set(ctx, key, string(userJSON))

	return user, nil
}
