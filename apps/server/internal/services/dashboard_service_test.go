package services

import (
	"context"
	"testing"

	"github.com/giridhar-m-a/custom_form_app/internal/db/sqlc"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetDashboardData(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	service := NewDashboardService(mockRepo)

	userID := uuid.NullUUID{UUID: uuid.New(), Valid: true}

	mockRepo.On("GetTotalForms", mock.Anything, userID).Return(int64(10), nil)
	mockRepo.On("GetTotalSubmissions", mock.Anything, userID).Return(int64(50), nil)
	mockRepo.On("GetTotalActiveForms", mock.Anything, userID).Return(int64(5), nil)
	mockRepo.On("GetTotalClosedForms", mock.Anything, userID).Return(int64(2), nil)
	mockRepo.On("GetTotalInvitations", mock.Anything, userID).Return(int64(100), nil)
	
	expectedSubmissions := []sqlc.GetFormSubmissionsByMonthRow{
		{
			Month:            "January",
			TotalSubmissions: 10,
		},
	}
	mockRepo.On("GetFormSubmissionsByMonth", mock.Anything, userID).Return(expectedSubmissions, nil)

	res, err := service.GetDashboardData(context.Background(), userID)

	assert.NoError(t, err)
	assert.Equal(t, int64(10), res.TotalForms)
	assert.Equal(t, int64(50), res.TotalSubmissions)
	assert.Equal(t, "January", res.SubmissionsByMonth[0].Month)
	mockRepo.AssertExpectations(t)
}
