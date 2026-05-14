package services

import (
	"context"
	"testing"

	"github.com/giridhar-m-a/custom_form_app/internal/db/sqlc"
	"github.com/giridhar-m-a/custom_form_app/internal/dto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateForm(t *testing.T) {
	mockFormRepo := new(MockFormsRepository)
	mockFieldRepo := new(MockFormFieldsRepository)
	mockOptionRepo := new(MockFormFieldOptionsRepository)
	
	service := NewFormService(mockFormRepo, mockFieldRepo, mockOptionRepo, nil)

	userID := uuid.New().String()
	formDTO := dto.CreateFormDTO{
		Title: "Test Form",
		Description: func(s string) *string { return &s }("Test Description"),
	}

	expectedForm := sqlc.Form{
		FormID: uuid.New(),
		FormTitle: "Test Form",
	}

	mockFormRepo.On("CreateForm", mock.Anything, mock.Anything).Return(expectedForm, nil)

	form, err := service.CreateForm(context.Background(), formDTO, userID)

	assert.NoError(t, err)
	assert.Equal(t, expectedForm.FormTitle, form.FormTitle)
	mockFormRepo.AssertExpectations(t)
}

func TestGetForms(t *testing.T) {
	mockFormRepo := new(MockFormsRepository)
	service := NewFormService(mockFormRepo, nil, nil, nil)

	userID := uuid.New().String()
	query := dto.ListFormQuery{
		Query: dto.Query{
			Page: 1,
			Limit: 10,
		},
	}

	expectedForms := []sqlc.ListFormsRow{
		{
			FormID: uuid.New(),
			FormTitle: "Form 1",
			TotalCount: 1,
		},
	}

	mockFormRepo.On("GetFormsList", mock.Anything, mock.Anything).Return(expectedForms, nil)

	res, err := service.GetForms(context.Background(), userID, query)

	assert.NoError(t, err)
	assert.Len(t, res.Forms, 1)
	assert.Equal(t, "Form 1", res.Forms[0].FormTitle)
	mockFormRepo.AssertExpectations(t)
}

func TestGetSingleForm(t *testing.T) {
	mockFormRepo := new(MockFormsRepository)
	service := NewFormService(mockFormRepo, nil, nil, nil)

	formID := uuid.New().String()
	expectedForm := sqlc.Form{
		FormID: uuid.MustParse(formID),
		FormTitle: "Single Form",
	}

	mockFormRepo.On("GetFormByID", formID, mock.Anything).Return(expectedForm, nil)

	form, err := service.GetSingleForm(context.Background(), formID)

	assert.NoError(t, err)
	assert.Equal(t, expectedForm.FormTitle, form.FormTitle)
	mockFormRepo.AssertExpectations(t)
}

func TestDeleteForm(t *testing.T) {
	mockFormRepo := new(MockFormsRepository)
	service := NewFormService(mockFormRepo, nil, nil, nil)

	formID := uuid.New().String()
	form := sqlc.Form{
		FormID: uuid.MustParse(formID),
	}
	expectedResult := sqlc.DeleteFormRow{
		FormID: form.FormID,
	}

	mockFormRepo.On("GetFormByID", formID, mock.Anything).Return(form, nil)
	mockFormRepo.On("DeleteForm", formID, mock.Anything).Return(expectedResult, nil)

	res, err := service.DeleteForm(context.Background(), formID)

	assert.NoError(t, err)
	assert.Equal(t, expectedResult.FormID, res.FormID)
	mockFormRepo.AssertExpectations(t)
}
