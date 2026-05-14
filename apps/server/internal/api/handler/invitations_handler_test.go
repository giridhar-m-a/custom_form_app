package handler

import (
	"context"
	"mime/multipart"
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

type MockInvitationService struct {
	mock.Mock
}

func (m *MockInvitationService) CreateInvitation(file *multipart.FileHeader, formID uuid.UUID, userID uuid.UUID, ctx context.Context) (int, int, error) {
	args := m.Called(file, formID, userID, ctx)
	return args.Int(0), args.Int(1), args.Error(2)
}

func (m *MockInvitationService) CreateSingleInvitation(data dto.CreateInvitationDTO, userID string, ctx context.Context) (sqlc.CreateInvitationRow, error) {
	args := m.Called(data, userID, ctx)
	return args.Get(0).(sqlc.CreateInvitationRow), args.Error(1)
}

func (m *MockInvitationService) UpdateInvitationStatus(status dto.UpdateInvitationDTO, resendID uuid.UUID, ctx context.Context) (sqlc.UpdateInvitationStatusRow, error) {
	args := m.Called(status, resendID, ctx)
	return args.Get(0).(sqlc.UpdateInvitationStatusRow), args.Error(1)
}

func (m *MockInvitationService) GetInvitationByFormId(params dto.InvitationListQueryDto, ctx context.Context) (dto.InvitationListDto, error) {
	args := m.Called(params, ctx)
	return args.Get(0).(dto.InvitationListDto), args.Error(1)
}

func (m *MockInvitationService) DeleteInvitation(invitationID uuid.UUID, ctx context.Context) error {
	args := m.Called(invitationID, ctx)
	return args.Error(0)
}

func TestDeleteInvitationHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockInvitationService)
	handler := &invitationHandler{svc: mockService}

	invitationID := uuid.New()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: invitationID.String()}}

	mockService.On("DeleteInvitation", invitationID, mock.Anything).Return(nil)

	handler.DeleteInvitation(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}
