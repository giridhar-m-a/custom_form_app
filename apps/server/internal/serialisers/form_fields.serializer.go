package serializers

import (
	"github.com/giridhar-m-a/custom_form_app/internal/dto"
	"github.com/giridhar-m-a/custom_form_app/internal/utils"
)

func MapCreatedFormFieldToResponse(input dto.CreatedFormFieldDTO) dto.FormFieldResponseDto {
	options := make([]dto.FormFieldOptionResponseDto, len(input.Options))
	for i, opt := range input.Options {
		options[i] = MapCreatedFormFieldOptionToResponse(opt)
	}

	return dto.FormFieldResponseDto{
		FormId:     input.FormId.String(),
		FieldID:    input.FieldID.String(),
		FieldLabel: input.FieldLabel,
		FieldType:  string(input.FieldType),
		IsRequired: input.IsRequired.Bool,
		Ordering:   input.Ordering,
		Options:    options,
	}
}

func MapCreatedFormFieldOptionToResponse(input dto.CreatedFormFieldOptionDTO) dto.FormFieldOptionResponseDto {

	return dto.FormFieldOptionResponseDto{
		OptionID:    input.OptionID.String(),
		FieldId:     utils.NullUUIDToString(input.FieldId),
		OptionLabel: input.OptionLabel,
		Ordering:    input.Ordering,
		IsAnswer:    input.IsAnswer,
	}
}
