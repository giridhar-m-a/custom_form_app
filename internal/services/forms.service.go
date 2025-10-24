package services

import (
	"context"
	"database/sql"

	"github.com/giridhar-m-a/custom_form_app/internal/db/sqlc"
	"github.com/giridhar-m-a/custom_form_app/internal/dto"
	"github.com/giridhar-m-a/custom_form_app/internal/repositories"
	"github.com/google/uuid"
)

type FormService interface {
	CreateForm(ctx context.Context, form dto.CreateFormDTO, userID string) (sqlc.CreateFormRow, error)
}

type formService struct {
	formRepo repositories.FormsRepository
}

// NewFormService creates a FormService backed by the provided FormsRepository.
// The returned service uses formRepo to persist and retrieve form data.
func NewFormService(formRepo repositories.FormsRepository) FormService {
	return &formService{formRepo: formRepo}
}

func (s *formService) CreateForm(ctx context.Context, form dto.CreateFormDTO, userID string) (sqlc.CreateFormRow, error) {

	user, err := uuid.Parse(userID)
	if err != nil {
		return sqlc.CreateFormRow{}, err
	}

	CreatedBy := uuid.NullUUID{
		UUID:  user, // userID should be of type uuid.UUID
		Valid: true,
	}
	formDescription := sql.NullString{
		String: form.Description,
		Valid:  form.Description != "",
	}

	return s.formRepo.CreateForm(sqlc.CreateFormParams{
		FormTitle:       form.Title,
		FormDescription: formDescription,
		CreatedBy:       CreatedBy,
	}, ctx)
}