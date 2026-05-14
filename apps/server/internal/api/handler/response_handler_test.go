package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/giridhar-m-a/custom_form_app/internal/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockResponseService struct {
	mock.Mock
}

func (m *MockResponseService) SubmitForm(ctx context.Context, req dto.CreateSubmissionRequest) (dto.SubmissionResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(dto.SubmissionResponse), args.Error(1)
}

func (m *MockResponseService) GetSubmissions(ctx context.Context, req dto.GetSubmissionsRequest) (dto.GetSubmissionsResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(dto.GetSubmissionsResponse), args.Error(1)
}

func (m *MockResponseService) GetSingleSubmission(ctx context.Context, submissionID string) (dto.SubmissionResponse, error) {
	args := m.Called(ctx, submissionID)
	return args.Get(0).(dto.SubmissionResponse), args.Error(1)
}

func TestCreateSubmissionHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockResponseService)
	handler := &responseHandler{responseService: mockService}

	req := dto.CreateSubmissionRequest{
		FormID: "form-123",
	}
	body, _ := json.Marshal(req)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/response", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	mockService.On("SubmitForm", mock.Anything, req).Return(dto.SubmissionResponse{}, nil)

	handler.CreateSubmission(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestGetSubmissionsHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockResponseService)
	handler := &responseHandler{responseService: mockService}

	formId := "form-123"
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "formId", Value: formId}}
	c.Request = httptest.NewRequest("GET", "/response/"+formId+"?page=1&limit=10", nil)

	mockService.On("GetSubmissions", mock.Anything, mock.MatchedBy(func(req dto.GetSubmissionsRequest) bool {
		return req.FormID == formId
	})).Return(dto.GetSubmissionsResponse{
		Submissions: []dto.SubmissionList{},
	}, nil)

	handler.GetSubmissions(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestGetSingleSubmissionHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockResponseService)
	handler := &responseHandler{responseService: mockService}

	submissionId := "sub-123"
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "submissionId", Value: submissionId}}

	mockService.On("GetSingleSubmission", mock.Anything, submissionId).Return(dto.SubmissionResponse{}, nil)

	handler.GetSingleSubmission(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}
