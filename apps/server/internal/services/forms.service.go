package services

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/giridhar-m-a/custom_form_app/internal/db/sqlc"
	"github.com/giridhar-m-a/custom_form_app/internal/dto"
	"github.com/giridhar-m-a/custom_form_app/internal/repositories"
	"github.com/giridhar-m-a/custom_form_app/internal/utils"
	"github.com/google/uuid"
)

type FormService interface {
	CreateForm(ctx context.Context, form dto.CreateFormDTO, userID string) (sqlc.CreateFormRow, error)
	CreateFormFields(ctx context.Context, form dto.CreateFormFieldsDTO, userID string) ([]dto.CreatedFormFieldDTO, error)
	GetForms(ctx context.Context, userID string, query dto.ListFormQuery) (dto.FormListResponse, error)
	GetSingleForm(ctx context.Context, formID string) (dto.FormResponse, error)
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
			FieldType:  sqlc.NullFormFieldType{FormFieldType: field.FieldType},
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
			FieldType:  createdFormField.FieldType.FormFieldType,
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

func (s *formService) GetForms(ctx context.Context, userID string, query dto.ListFormQuery) (dto.FormListResponse, error) {

	user, err := utils.ConvertStringToUUID(userID)
	if err != nil {
		return dto.FormListResponse{}, err
	}

	// Set default pagination values if not provided
	page := query.Page
	if page == 0 {
		page = 1
	}
	limit := query.Limit
	if limit == 0 {
		limit = 10
	}

	// Set default sort value if not provided
	sortBy := query.Sort
	if sortBy == "" {
		sortBy = "updated"
	}

	// Now we use the properly typed parameters from the regenerated sqlc code
	var search sql.NullString
	if query.Search != "" {
		search = sql.NullString{String: query.Search, Valid: true}
	}

	var status sqlc.NullFormStatus
	if query.Status != "" {
		status = sqlc.NullFormStatus{FormStatus: sqlc.FormStatus(query.Status), Valid: true}
	}

	var access sqlc.NullFormAccess
	if query.Access != "" {
		access = sqlc.NullFormAccess{FormAccess: sqlc.FormAccess(query.Access), Valid: true}
	}

	params := sqlc.ListFormsParams{
		CreatedBy:  utils.ConvertUUIDToNullUUID(user.String()),
		Offset:     sql.NullInt32{Int32: int32((page - 1) * limit), Valid: true},
		Limit:      sql.NullInt32{Int32: int32(limit), Valid: true},
		Search:     search,
		ShortBy:    sql.NullString{String: sortBy, Valid: true},
		FormStatus: status,
		FormAccess: access,
	}

	forms, err := s.formRepo.GetFormsList(params, ctx)
	if err != nil {
		fmt.Printf("Error getting forms: %v", err)
		return dto.FormListResponse{}, err
	}

	// Get total count from the first row (all rows have the same total_count due to window function)
	var totalCount int64 = 0
	if len(forms) > 0 {
		totalCount = forms[0].TotalCount
	}

	// Calculate total pages
	totalPages := int(totalCount) / limit
	if int(totalCount)%limit != 0 {
		totalPages++
	}
	if totalPages == 0 {
		totalPages = 1
	}

	var formResponses []dto.FormResponse
	for _, form := range forms {
		formResponses = append(formResponses, dto.FormResponse{
			ID:          form.FormID.String(),
			Title:       form.FormTitle,
			Description: form.FormDescription.String,
			CreatedBy:   form.CreatedBy.UUID.String(),
			Status:      string(form.FormStatus.FormStatus),
			CreatedAt:   form.FormCreatedAt.Time.String(),
			UpdatedAt:   form.FormUpdatedAt.Time.String(),
			Access:      string(form.FormAccess.FormAccess),
		})
	}

	return dto.FormListResponse{
		Forms: formResponses,
		Total: int(totalCount),
		Page:  page,
		Limit: limit,
		Pages: totalPages,
	}, nil
}

func (s *formService) GetSingleForm(ctx context.Context, formID string) (dto.FormResponse, error) {
	form, err := s.formRepo.GetFormByID(formID, ctx)
	if err != nil {
		return dto.FormResponse{}, err
	}
	return dto.FormResponse{
		ID:          form.FormID.String(),
		Title:       form.FormTitle,
		Description: form.FormDescription.String,
		CreatedBy:   form.CreatedBy.UUID.String(),
		Status:      string(form.FormStatus.FormStatus),
		CreatedAt:   form.FormCreatedAt.Time.String(),
		UpdatedAt:   form.FormUpdatedAt.Time.String(),
		Access:      string(form.FormAccess.FormAccess),
	}, nil
}
