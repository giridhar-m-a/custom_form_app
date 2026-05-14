package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/giridhar-m-a/custom_form_app/internal/dto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDashboardService struct {
	mock.Mock
}

func (m *MockDashboardService) GetDashboardData(ctx context.Context, userID uuid.NullUUID) (dto.DashboardResponse, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(dto.DashboardResponse), args.Error(1)
}

func TestGetDashboardDataHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockDashboardService)
	handler := &dashboardHandler{dashboardService: mockService}

	userID := uuid.New()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", userID.String())

	expectedData := dto.DashboardResponse{
		TotalForms: 5,
	}

	mockService.On("GetDashboardData", mock.Anything, uuid.NullUUID{UUID: userID, Valid: true}).Return(expectedData, nil)

	handler.GetDashboardData(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}
