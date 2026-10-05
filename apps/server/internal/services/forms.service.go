package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/giridhar-m-a/custom_form_app/internal/db/sqlc"
	"github.com/giridhar-m-a/custom_form_app/internal/dto"
	"github.com/giridhar-m-a/custom_form_app/internal/repositories"
	"github.com/giridhar-m-a/custom_form_app/internal/utils"
	"github.com/google/uuid"
)

type FormService interface {
	CreateForm(ctx context.Context, form dto.CreateFormDTO, userID string) (sqlc.Form, error)
	CreateFormFields(ctx context.Context, form dto.CreateFormFieldsDTO, userID string) ([]dto.CreatedFormFieldDTO, error)
	GetForms(ctx context.Context, userID string, query dto.ListFormQuery) (dto.FormListResponse, error)
	GetSingleForm(ctx context.Context, formID string) (sqlc.Form, error)
	UpdateForm(ctx context.Context, form dto.UpdateFormDTO, formID string) (sqlc.Form, error)
	DeleteForm(ctx context.Context, formID string) (sqlc.DeleteFormRow, error)
	GetFormFieldsByFormId(ctx context.Context, formId string) ([]dto.CreatedFormFieldDTO, error)
	UpdateFormFields(ctx context.Context, form dto.UpdateFormFieldsDTO) ([]dto.CreatedFormFieldDTO, error)
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

func (s *formService) CreateForm(ctx context.Context, form dto.CreateFormDTO, userID string) (sqlc.Form, error) {
	// Parse user ID
	user, err := uuid.Parse(userID)
	if err != nil {
		return sqlc.Form{}, err
	}

	CreatedBy := uuid.NullUUID{
		UUID:  user,
		Valid: true,
	}

	// Convert optional description
	var formDescription sql.NullString
	if form.Description != nil {
		formDescription = sql.NullString{
			String: *form.Description,
			Valid:  true,
		}
	}

	// Convert optional times
	var scheduledTime, closingTime sql.NullTime
	if form.ScheduledTime != nil {
		scheduledTime = sql.NullTime{Time: *form.ScheduledTime, Valid: true}
	}
	if form.ClosingTime != nil {
		closingTime = sql.NullTime{Time: *form.ClosingTime, Valid: true}
	}

	formAccess := sqlc.NullFormAccess{
		FormAccess: form.FormAccess,
		Valid:      form.FormAccess != "",
	}

	if form.FormAccess == "" {
		formAccess = sqlc.NullFormAccess{
			FormAccess: sqlc.FormAccessRestricted,
			Valid:      true,
		}
	}

	isScheduled := utils.BoolPtrToNullBool(form.IsScheduled)

	return s.formRepo.CreateForm(sqlc.CreateFormParams{
		FormTitle:       form.Title,
		FormDescription: formDescription,
		CreatedBy:       CreatedBy,
		FormAccess:      formAccess,
		ScheduledTime:   scheduledTime,
		ClosingTime:     closingTime,
		IsScheduled:     isScheduled,
	}, ctx)
}

func (s *formService) CreateFormFields(ctx context.Context, form dto.CreateFormFieldsDTO, userID string) ([]dto.CreatedFormFieldDTO, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return []dto.CreatedFormFieldDTO{}, err
	}
	defer tx.Rollback()
	var createdFormFields []dto.CreatedFormFieldDTO

	formID, err := uuid.Parse(form.FormID)
	if err != nil {
		return []dto.CreatedFormFieldDTO{}, err
	}

	for _, field := range form.FormFields {
		var formField dto.CreatedFormFieldDTO
		createdFormField, err := s.createFormField(field, tx, formID, ctx)
		if err != nil {
			return []dto.CreatedFormFieldDTO{}, err
		}
		formField = dto.CreatedFormFieldDTO{
			FormId:     createdFormField.FormId,
			FieldID:    createdFormField.FieldID,
			FieldLabel: createdFormField.FieldLabel,
			FieldType:  createdFormField.FieldType,
			IsRequired: createdFormField.IsRequired,
			Ordering:   createdFormField.Ordering,
			Options:    []dto.CreatedFormFieldOptionDTO{},
		}
		for _, option := range field.Options {
			createdOption, err := s.createFormFieldOption(option, tx, createdFormField.FieldID, ctx)
			if err != nil {
				return []dto.CreatedFormFieldDTO{}, err
			}
			formField.Options = append(formField.Options, dto.CreatedFormFieldOptionDTO{
				OptionID:    createdOption.OptionID,
				FieldId:     createdOption.FieldId,
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
		CreatedBy:  utils.ConvertStringToNullUUID(user.String()),
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

	return dto.FormListResponse{
		Forms: forms,
		Total: int(totalCount),
		Page:  page,
		Limit: limit,
		Pages: totalPages,
	}, nil
}

func (s *formService) GetSingleForm(ctx context.Context, formID string) (sqlc.Form, error) {
	return s.formRepo.GetFormByID(formID, ctx)
}

func (s *formService) UpdateForm(ctx context.Context, form dto.UpdateFormDTO, formID string) (sqlc.Form, error) {
	// Convert formID to uuid.UUID
	id, err := utils.ConvertStringToUUID(formID)
	if err != nil {
		return sqlc.Form{}, err
	}

	// Convert optional strings
	formTitle := utils.ConvertStringToNullString(*form.Title)
	formDescription := utils.ConvertStringToNullString(*form.Description)

	// Convert optional FormStatus
	var formStatus sqlc.NullFormStatus
	if form.Status != nil {
		formStatus = sqlc.NullFormStatus{
			FormStatus: *form.Status,
			Valid:      true,
		}
	}

	// Convert optional FormAccess
	var formAccess sqlc.NullFormAccess
	if form.Access != nil {
		formAccess = sqlc.NullFormAccess{
			FormAccess: *form.Access,
			Valid:      true,
		}
	}

	// Convert optional scheduling fields
	var schedulingID uuid.NullUUID
	if form.SchedulingID != nil {
		u, err := uuid.Parse(*form.SchedulingID)
		if err == nil {
			schedulingID = uuid.NullUUID{UUID: u, Valid: true}
		}
	}

	var scheduledTime, closingTime sql.NullTime
	if form.ScheduledTime != nil {
		scheduledTime = sql.NullTime{Time: *form.ScheduledTime, Valid: true}
	}
	if form.ClosingTime != nil {
		closingTime = sql.NullTime{Time: *form.ClosingTime, Valid: true}
	}

	var isScheduleCompleted, isScheduled sql.NullBool
	if form.IsScheduleCompleted != nil {
		isScheduleCompleted = sql.NullBool{Bool: *form.IsScheduleCompleted, Valid: true}
	}
	if form.IsScheduled != nil {
		isScheduled = sql.NullBool{Bool: *form.IsScheduled, Valid: true}
	}

	// Prepare payload
	formPayload := sqlc.UpdateFormParams{
		FormID:              id,
		FormTitle:           formTitle,
		FormDescription:     formDescription,
		FormStatus:          formStatus,
		FormAccess:          formAccess,
		SchedulingID:        schedulingID,
		ScheduledTime:       scheduledTime,
		ClosingTime:         closingTime,
		IsScheduleCompleted: isScheduleCompleted,
		IsScheduled:         isScheduled,
	}

	return s.formRepo.UpdateForm(formPayload, ctx)
}

func (s *formService) DeleteForm(ctx context.Context, formID string) (sqlc.DeleteFormRow, error) {
	return s.formRepo.DeleteForm(formID, ctx)
}

func (s *formService) GetFormFieldsByFormId(ctx context.Context, formId string) ([]dto.CreatedFormFieldDTO, error) {
	fields, err := s.fieldRepo.GetFormFieldsByFormId(formId, ctx)
	if err != nil {
		return nil, err
	}
	var formFields []dto.CreatedFormFieldDTO
	for _, field := range fields {
		var options []dto.CreatedFormFieldOptionDTO
		err := json.Unmarshal(field.Options, &options)
		if err != nil {
			return nil, err
		}

		var marshalledOptions []dto.CreatedFormFieldOptionDTO
		for _, opt := range options {
			marshalledOptions = append(marshalledOptions, dto.CreatedFormFieldOptionDTO{
				OptionID:    opt.OptionID,
				FieldId:     opt.FieldId,
				OptionLabel: opt.OptionLabel,
				Ordering:    opt.Ordering,
				IsAnswer:    opt.IsAnswer,
			})
		}

		formFields = append(formFields, dto.CreatedFormFieldDTO{
			FormId:     field.FormId,
			FieldID:    field.FieldId,
			FieldLabel: field.FieldLabel,
			FieldType:  field.FieldType.FormFieldType,
			IsRequired: field.IsRequired,
			Ordering:   field.Ordering,
			Options:    marshalledOptions,
		})
	}
	return formFields, nil
}

func (s *formService) UpdateFormFields(ctx context.Context, form dto.UpdateFormFieldsDTO) ([]dto.CreatedFormFieldDTO, error) {
	formId, err := utils.ConvertStringToUUID(form.FormID)
	if err != nil {
		return nil, err
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	var createdFormFields []dto.CreatedFormFieldDTO
	for _, field := range form.FormFields {
		var formField dto.CreatedFormFieldDTO
		if field.FieldId != "" {
			fieldId, err := utils.ConvertStringToUUID(field.FieldId)
			if err != nil {
				return nil, err
			}
			updatedFormField, err := s.updateFormField(field, tx, fieldId, ctx)
			if err != nil {
				return nil, err
			}
			formField = dto.CreatedFormFieldDTO{
				FormId:     updatedFormField.FormId,
				FieldID:    updatedFormField.FieldID,
				FieldLabel: updatedFormField.FieldLabel,
				FieldType:  updatedFormField.FieldType,
				IsRequired: updatedFormField.IsRequired,
				Options:    []dto.CreatedFormFieldOptionDTO{},
			}
			for _, option := range field.Options {

				if option.OptionId != "" {

					optionId, err := utils.ConvertStringToUUID(option.OptionId)
					if err != nil {
						return nil, err
					}

					updatedOption, err := s.updateFormFieldOption(option, tx, optionId, ctx)
					if err != nil {
						return nil, err
					}
					formField.Options = append(formField.Options, dto.CreatedFormFieldOptionDTO{
						OptionID:    updatedOption.OptionID,
						FieldId:     updatedOption.FieldId,
						OptionLabel: updatedOption.OptionLabel,
						Ordering:    updatedOption.Ordering,
						IsAnswer:    updatedOption.IsAnswer,
					})
				} else {
					createdOption, err := s.createFormFieldOption(option, tx, fieldId, ctx)
					if err != nil {
						return nil, err
					}
					formField.Options = append(formField.Options, dto.CreatedFormFieldOptionDTO{
						OptionID:    createdOption.OptionID,
						FieldId:     createdOption.FieldId,
						OptionLabel: createdOption.OptionLabel,
						Ordering:    createdOption.Ordering,
						IsAnswer:    createdOption.IsAnswer,
					})
				}
			}
		} else {
			createdFormField, err := s.createFormField(field, tx, formId, ctx)
			if err != nil {
				return nil, err
			}
			formField = createdFormField
			for _, option := range field.Options {
				createdOption, err := s.createFormFieldOption(option, tx, formField.FieldID, ctx)
				if err != nil {
					return nil, err
				}
				formField.Options = append(formField.Options, dto.CreatedFormFieldOptionDTO{
					OptionID:    createdOption.OptionID,
					FieldId:     createdOption.FieldId,
					OptionLabel: createdOption.OptionLabel,
					Ordering:    createdOption.Ordering,
					IsAnswer:    createdOption.IsAnswer,
				})
			}

		}
		createdFormFields = append(createdFormFields, formField)
	}

	for _, option := range form.RemovedOptions {
		_, err = s.deleteFormFieldOption(option, tx, ctx)
		if err != nil {
			return nil, err
		}
	}

	for _, field := range form.RemovedFields {
		_, err = s.deleteFormField(field, tx, ctx)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return createdFormFields, nil

}

func (s *formService) createFormField(field dto.CreateFormFieldDTO, tx *sql.Tx, formId uuid.UUID, ctx context.Context) (dto.CreatedFormFieldDTO, error) {
	formRepo := s.fieldRepo.FormFieldRepoWithTx(tx)
	createdFormFieldParams := sqlc.CreateFormFieldParams{
		FormID:     formId,
		FieldLabel: field.FieldLabel,
		FieldType:  sqlc.NullFormFieldType{FormFieldType: field.FieldType, Valid: field.FieldType != ""},
		IsRequired: utils.ConvertBoolToNullBool(field.IsRequired),
		Ordering:   utils.ConvertIntToInt32(field.Ordering),
	}

	createdFormField, err := formRepo.CreateFormField(createdFormFieldParams, ctx)
	if err != nil {
		return dto.CreatedFormFieldDTO{}, err
	}
	return dto.CreatedFormFieldDTO{
		FormId:     createdFormField.FormID,
		FieldID:    createdFormField.FieldID,
		FieldLabel: createdFormField.FieldLabel,
		FieldType:  createdFormField.FieldType.FormFieldType,
		IsRequired: createdFormField.IsRequired,
		Ordering:   createdFormField.Ordering,
		Options:    []dto.CreatedFormFieldOptionDTO{},
	}, nil
}

func (s *formService) createFormFieldOption(option dto.CreateFormFieldOptionDTO, tx *sql.Tx, fieldId uuid.UUID, ctx context.Context) (dto.CreatedFormFieldOptionDTO, error) {
	optionRepo := s.fieldOptionRepo.FormFieldOptionsRepoWithTx(tx)
	createdFormFieldOptionParams := sqlc.CreateFieldOptionParams{
		FieldID:     utils.ConvertStringToNullUUID(fieldId.String()),
		OptionLabel: option.OptionLabel,
		Ordering:    utils.ConvertIntToInt32(option.Ordering),
		IsAnswer:    utils.ConvertBoolToNullBool(option.IsAnswer),
	}

	createdFormFieldOption, err := optionRepo.CreateFieldOption(createdFormFieldOptionParams, ctx)
	if err != nil {
		return dto.CreatedFormFieldOptionDTO{}, err
	}
	return dto.CreatedFormFieldOptionDTO{
		OptionID:    createdFormFieldOption.OptionID,
		FieldId:     createdFormFieldOption.FieldID,
		OptionLabel: createdFormFieldOption.OptionLabel,
		Ordering:    createdFormFieldOption.Ordering,
		IsAnswer:    createdFormFieldOption.IsAnswer.Bool,
	}, nil
}

func (s *formService) updateFormFieldOption(option dto.CreateFormFieldOptionDTO, tx *sql.Tx, optionId uuid.UUID, ctx context.Context) (dto.CreatedFormFieldOptionDTO, error) {
	optionRepo := s.fieldOptionRepo.FormFieldOptionsRepoWithTx(tx)
	updatedFormFieldOptionParams := sqlc.UpdateFieldOptionParams{
		OptionLabel: option.OptionLabel,
		Ordering:    utils.ConvertIntToInt32(option.Ordering),
		IsAnswer:    utils.ConvertBoolToNullBool(option.IsAnswer),
		OptionID:    optionId,
	}
	updatedFormFieldOption, err := optionRepo.UpdateFieldOption(updatedFormFieldOptionParams, ctx)
	if err != nil {
		return dto.CreatedFormFieldOptionDTO{}, err
	}
	return dto.CreatedFormFieldOptionDTO{
		OptionID:    updatedFormFieldOption.OptionID,
		FieldId:     updatedFormFieldOption.FieldID,
		OptionLabel: updatedFormFieldOption.OptionLabel,
		Ordering:    updatedFormFieldOption.Ordering,
		IsAnswer:    updatedFormFieldOption.IsAnswer.Bool,
	}, nil
}

func (s *formService) updateFormField(option dto.CreateFormFieldDTO, tx *sql.Tx, fieldId uuid.UUID, ctx context.Context) (dto.CreatedFormFieldDTO, error) {
	fieldRepo := s.fieldRepo.FormFieldRepoWithTx(tx)
	updatedFormFieldParams := sqlc.UpdateFormFieldParams{
		FieldLabel: option.FieldLabel,
		FieldType:  sqlc.NullFormFieldType{FormFieldType: option.FieldType, Valid: option.FieldType != ""},
		IsRequired: utils.ConvertBoolToNullBool(option.IsRequired),
		Ordering:   utils.ConvertIntToInt32(option.Ordering),
		FieldID:    fieldId,
	}
	updatedFormField, err := fieldRepo.UpdateFormField(updatedFormFieldParams, ctx)
	if err != nil {
		return dto.CreatedFormFieldDTO{}, err
	}
	return dto.CreatedFormFieldDTO{
		FormId:     updatedFormField.FormID,
		FieldID:    updatedFormField.FieldID,
		FieldLabel: updatedFormField.FieldLabel,
		FieldType:  updatedFormField.FieldType.FormFieldType,
		IsRequired: updatedFormField.IsRequired,
		Ordering:   updatedFormField.Ordering,
		Options:    []dto.CreatedFormFieldOptionDTO{},
	}, nil
}

func (s *formService) deleteFormFieldOption(optionId string, tx *sql.Tx, ctx context.Context) (sqlc.DeleteFieldOptionRow, error) {
	optionRepo := s.fieldOptionRepo.FormFieldOptionsRepoWithTx(tx)
	return optionRepo.DeleteFieldOption(optionId, ctx)
}

func (s *formService) deleteFormField(fieldId string, tx *sql.Tx, ctx context.Context) (sqlc.DeleteFormFieldRow, error) {
	fieldRepo := s.fieldRepo.FormFieldRepoWithTx(tx)
	return fieldRepo.DeleteFormField(fieldId, ctx)
}
