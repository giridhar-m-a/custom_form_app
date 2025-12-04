package services

import (
	"context"
	"database/sql"

	"github.com/giridhar-m-a/custom_form_app/internal/db/sqlc"
	"github.com/giridhar-m-a/custom_form_app/internal/dto"
	"github.com/giridhar-m-a/custom_form_app/internal/repositories"
	"github.com/giridhar-m-a/custom_form_app/internal/utils"
	"github.com/google/uuid"
)

type FormService interface {
	CreateForm(ctx context.Context, form dto.CreateFormDTO, userID string) (sqlc.CreateFormRow, error)
	CreateFormFields(ctx context.Context, form dto.CreateFormFieldsDTO, userID string) ([]dto.CreatedFormFieldDTO, error)
}

type formService struct {
	formRepo        repositories.FormsRepository
	fieldRepo       repositories.FormFieldsRepository
	fieldOptionRepo repositories.FormFieldOptionsRepository
	db              *sql.DB
}

func NewFormService(formRepo repositories.FormsRepository, fieldRepo repositories.FormFieldsRepository, fieldOptionRepo repositories.FormFieldOptionsRepository, db *sql.DB) FormService {
	return &formService{formRepo: formRepo, fieldRepo: fieldRepo, fieldOptionRepo: fieldOptionRepo, db: db}
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

func (s *formService) CreateFormFields(ctx context.Context, form dto.CreateFormFieldsDTO, userID string) ([]dto.CreatedFormFieldDTO, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return []dto.CreatedFormFieldDTO{}, err
	}
	defer tx.Rollback()
	fieldRepo := s.fieldRepo.FormFieldRepoWithTx(tx)
	optionRepo := s.fieldOptionRepo.FormFieldOptionsRepoWithTx(tx)
	var createdFormFields []dto.CreatedFormFieldDTO

	formID, err := uuid.Parse(form.FormID)
	if err != nil {
		return []dto.CreatedFormFieldDTO{}, err
	}

	for _, field := range form.FormFields {
		var formField dto.CreatedFormFieldDTO
		formFieldParams := sqlc.CreateFormFieldParams{
			FormID:     formID,
			FieldLabel: field.FieldLabel,
			FieldType:  field.FieldType,
			IsRequired: utils.ConvertBoolToNullBool(field.IsRequired),
			Ordering:   utils.ConvertIntToInt32(field.Ordering),
		}
		createdFormField, err := fieldRepo.CreateFormField(formFieldParams, ctx)
		if err != nil {
			return []dto.CreatedFormFieldDTO{}, err
		}
		formField = dto.CreatedFormFieldDTO{
			FormId:     createdFormField.FormID,
			FieldID:    createdFormField.FieldID,
			FieldLabel: createdFormField.FieldLabel,
			FieldType:  createdFormField.FieldType,
			IsRequired: createdFormField.IsRequired,
			Ordering:   createdFormField.Ordering,
			Options:    []dto.CreatedFormFieldOptionDTO{},
		}
		for _, option := range field.Options {
			fieldOptionParams := sqlc.CreateFieldOptionParams{
				FieldID:     utils.ConvertUUIDToNullUUID(createdFormField.FieldID.String()),
				OptionLabel: option.OptionLabel,
				Ordering:    utils.ConvertIntToInt32(option.Ordering),
				IsAnswer:    utils.ConvertBoolToNullBool(option.IsAnswer),
			}
			createdOption, err := optionRepo.CreateFieldOption(fieldOptionParams, ctx)
			if err != nil {
				return []dto.CreatedFormFieldDTO{}, err
			}
			formField.Options = append(formField.Options, dto.CreatedFormFieldOptionDTO{
				OptionID:    createdOption.OptionID,
				FieldId:     createdOption.FieldID,
				OptionLabel: createdOption.OptionLabel,
				Ordering:    createdOption.Ordering,
				IsAnswer:    createdOption.IsAnswer,
			})
		}
		createdFormFields = append(createdFormFields, formField)
	}

	if err := tx.Commit(); err != nil {
		return []dto.CreatedFormFieldDTO{}, err
	}

	return createdFormFields, nil
}
