package dto

import (
	"database/sql"

	"github.com/google/uuid"
)

type CreateFormFieldsDTO struct {
	FormID     string               `json:"formId" validate:"required,uuid"`
	FormFields []CreateFormFieldDTO `json:"formFields" validate:"required,dive"`
}

type CreateFormFieldDTO struct {
	FieldLabel string                     `json:"fieldLabel" validate:"required"`
	FieldType  string                     `json:"fieldType" validate:"required"`
	IsRequired bool                       `json:"isRequired"`
	Ordering   int                        `json:"ordering" validate:"required"`
	Options    []CreateFormFieldOptionDTO `json:"options" validate:"dive"`
}

type CreateFormFieldOptionDTO struct {
	OptionLabel string `json:"optionLabel" validate:"required"`
	Ordering    int    `json:"ordering" validate:"required"`
	IsAnswer    bool   `json:"isAnswer"`
}

type CreatedFormFieldDTO struct {
	FormId     uuid.UUID                   `json:"formId"`
	FieldID    uuid.UUID                   `json:"fieldId"`
	FieldLabel string                      `json:"fieldLabel"`
	FieldType  string                      `json:"fieldType"`
	IsRequired sql.NullBool                `json:"isRequired"`
	Ordering   int32                       `json:"ordering"`
	Options    []CreatedFormFieldOptionDTO `json:"options"`
}

type CreatedFormFieldOptionDTO struct {
	OptionID    uuid.UUID     `json:"optionId"`
	FieldId     uuid.NullUUID `json:"fieldId"`
	OptionLabel string        `json:"optionLabel"`
	Ordering    int32         `json:"ordering"`
	IsAnswer    sql.NullBool  `json:"isAnswer"`
}
