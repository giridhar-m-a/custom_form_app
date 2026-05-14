package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/giridhar-m-a/custom_form_app/internal/db/sqlc"
	"github.com/giridhar-m-a/custom_form_app/internal/dto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockFormService struct {
	mock.Mock
}

func (m *MockFormService) CreateForm(ctx context.Context, form dto.CreateFormDTO, userID string) (sqlc.Form, error) {
	args := m.Called(ctx, form, userID)
	return args.Get(0).(sqlc.Form), args.Error(1)
}

func (m *MockFormService) CreateFormFields(ctx context.Context, form dto.CreateFormFieldsDTO, userID string) ([]dto.CreatedFormFieldDTO, error) {
	args := m.Called(ctx, form, userID)
	return args.Get(0).([]dto.CreatedFormFieldDTO), args.Error(1)
}

func (m *MockFormService) GetForms(ctx context.Context, userID string, query dto.ListFormQuery) (dto.FormListResponse, error) {
	args := m.Called(ctx, userID, query)
	return args.Get(0).(dto.FormListResponse), args.Error(1)
}

func (m *MockFormService) GetSingleForm(ctx context.Context, formID string) (sqlc.Form, error) {
	args := m.Called(ctx, formID)
	return args.Get(0).(sqlc.Form), args.Error(1)
}

func (m *MockFormService) UpdateForm(ctx context.Context, form dto.UpdateFormDTO, formID string) (sqlc.Form, error) {
	args := m.Called(ctx, form, formID)
	return args.Get(0).(sqlc.Form), args.Error(1)
}

func (m *MockFormService) DeleteForm(ctx context.Context, formID string) (sqlc.DeleteFormRow, error) {
	args := m.Called(ctx, formID)
	return args.Get(0).(sqlc.DeleteFormRow), args.Error(1)
}

func (m *MockFormService) GetFormFieldsByFormId(ctx context.Context, formId string) ([]dto.CreatedFormFieldDTO, error) {
	args := m.Called(ctx, formId)
	return args.Get(0).([]dto.CreatedFormFieldDTO), args.Error(1)
}

func (m *MockFormService) UpdateFormFields(ctx context.Context, form dto.UpdateFormFieldsDTO) ([]dto.CreatedFormFieldDTO, error) {
	args := m.Called(ctx, form)
	return args.Get(0).([]dto.CreatedFormFieldDTO), args.Error(1)
}

func (m *MockFormService) UpdateFormScheduleId(formID uuid.UUID, scheduleID uuid.NullUUID, invitationId uuid.NullUUID, ctx context.Context) (sqlc.Form, error) {
	args := m.Called(formID, scheduleID, invitationId, ctx)
	return args.Get(0).(sqlc.Form), args.Error(1)
}

func (m *MockFormService) SoftDeleteForm(formID string, ctx context.Context) error {
	args := m.Called(formID, ctx)
	return args.Error(0)
}

func TestGetSingleFormHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockFormService)
	handler := &formHandler{formService: mockService}

	formID := uuid.New().String()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "formID", Value: formID}}

	expectedForm := sqlc.Form{
		FormID:    uuid.MustParse(formID),
		FormTitle: "Test Form",
	}

	mockService.On("GetSingleForm", mock.Anything, formID).Return(expectedForm, nil)

	handler.GetSingleForm(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Form retrieved successfully", response["message"])
	mockService.AssertExpectations(t)
}

func TestCreateFormHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockFormService)
	handler := &formHandler{formService: mockService}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	payload := dto.CreateFormDTO{
		Title: "New Form",
	}
	body, _ := json.Marshal(payload)
	c.Request = httptest.NewRequest("POST", "/forms", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	userID := uuid.New().String()
	c.Set("userID", userID)

	expectedForm := sqlc.Form{
		FormID:    uuid.New(),
		FormTitle: "New Form",
	}

	mockService.On("CreateForm", mock.Anything, payload, userID).Return(expectedForm, nil)

	handler.CreateForm(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}
