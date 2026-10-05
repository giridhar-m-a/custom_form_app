package dto

import (
	"database/sql"

	"github.com/giridhar-m-a/custom_form_app/internal/db/sqlc"
	"github.com/google/uuid"
)

type CreateFormFieldsDTO struct {
	FormID     string               `json:"formId" validate:"required,uuid"`
	FormFields []CreateFormFieldDTO `json:"formFields" validate:"required,dive"`
}

type UpdateFormFieldsDTO struct {
	FormID         string               `json:"formId" validate:"required,uuid"`
	FormFields     []CreateFormFieldDTO `json:"formFields" validate:"required,dive"`
	RemovedFields  []string             `json:"removedFields"`
	RemovedOptions []string             `json:"removedFieldOptions"`
}

type CreateFormFieldDTO struct {
	FieldLabel string                     `json:"fieldLabel" validate:"required"`
	FieldType  sqlc.FormFieldType         `json:"fieldType" validate:"required"`
	IsRequired bool                       `json:"isRequired"`
	Ordering   int                        `json:"ordering" validate:"required"`
	Options    []CreateFormFieldOptionDTO `json:"options" validate:"dive"`
	FieldId    string                     `json:"fieldId" validate:"omitempty,uuid"`
}

type CreateFormFieldOptionDTO struct {
	OptionLabel string `json:"optionLabel" validate:"required"`
	Ordering    int    `json:"ordering" validate:"required"`
	IsAnswer    bool   `json:"isAnswer"`
	OptionId    string `json:"optionId" validate:"omitempty,uuid"`
	FieldId     string `json:"fieldId" validate:"omitempty,uuid"`
}

type CreatedFormFieldDTO struct {
	FormId     uuid.UUID                   `json:"formId"`
	FieldID    uuid.UUID                   `json:"fieldId"`
	FieldLabel string                      `json:"fieldLabel"`
	FieldType  sqlc.FormFieldType          `json:"fieldType"`
	IsRequired sql.NullBool                `json:"isRequired"`
	Ordering   int32                       `json:"ordering"`
	Options    []CreatedFormFieldOptionDTO `json:"options"`
}

type CreatedFormFieldOptionDTO struct {
	OptionID    uuid.UUID     `json:"optionId"`
	FieldId     uuid.NullUUID `json:"fieldId"`
	OptionLabel string        `json:"optionLabel"`
	Ordering    int32         `json:"ordering"`
	IsAnswer    bool          `json:"isAnswer"`
}

type FormFieldResponseDto struct {
	FormId     string                       `json:"formId"`
	FieldID    string                       `json:"fieldId"`
	FieldLabel string                       `json:"fieldLabel"`
	FieldType  string                       `json:"fieldType"`
	IsRequired bool                         `json:"isRequired"`
	Ordering   int32                        `json:"ordering"`
	Options    []FormFieldOptionResponseDto `json:"options"`
}

type FormFieldOptionResponseDto struct {
	OptionID    string `json:"optionId"`
	FieldId     string `json:"fieldId"`
	OptionLabel string `json:"optionLabel"`
	Ordering    int32  `json:"ordering"`
	IsAnswer    bool   `json:"isAnswer"`
}
