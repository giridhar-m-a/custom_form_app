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

func TestCreateSingleInvitation(t *testing.T) {
	mockRepo := new(MockInvitationRepository)
	mockForm := new(MockFormService)
	service := NewInvitationService(mockRepo, mockForm, nil)

	userID := uuid.New().String()
	invitationDTO := dto.CreateInvitationDTO{
		FormID: uuid.New().String(),
		Email:  "test@example.com",
		Name:   "Test User",
	}

	form := sqlc.Form{
		FormID: uuid.MustParse(invitationDTO.FormID),
		FormStatus: sqlc.NullFormStatus{FormStatus: sqlc.FormStatusDraft, Valid: true},
	}

	expectedInvitation := sqlc.CreateInvitationRow{
		InvitationID: uuid.New(),
		InvitedEmail: invitationDTO.Email,
	}

	mockForm.On("GetSingleForm", mock.Anything, invitationDTO.FormID).Return(form, nil)
	mockRepo.On("CreateSingleInvitation", mock.Anything, mock.Anything).Return(expectedInvitation, nil)

	res, err := service.CreateSingleInvitation(invitationDTO, userID, context.Background())

	assert.NoError(t, err)
	assert.Equal(t, expectedInvitation.InvitedEmail, res.InvitedEmail)
	mockRepo.AssertExpectations(t)
	mockForm.AssertExpectations(t)
}

func TestGetInvitationByFormId(t *testing.T) {
	mockRepo := new(MockInvitationRepository)
	service := NewInvitationService(mockRepo, nil, nil)

	formID := uuid.New().String()
	query := dto.InvitationListQueryDto{
		FormId: formID,
		Query: dto.Query{
			Page:   1,
			Limit:  10,
		},
	}

	expectedInvitations := []sqlc.Invitation{
		{
			InvitationID: uuid.New(),
			InvitedEmail: "test@example.com",
		},
	}

	mockRepo.On("GetInvitationByFormId", mock.Anything, mock.Anything).Return(expectedInvitations, nil)
	mockRepo.On("CountInvitationsByFormId", mock.Anything, mock.Anything).Return(int64(1), nil)

	res, err := service.GetInvitationByFormId(query, context.Background())

	assert.NoError(t, err)
	assert.Len(t, res.Invitations, 1)
	mockRepo.AssertExpectations(t)
}

func TestDeleteInvitation(t *testing.T) {
	mockRepo := new(MockInvitationRepository)
	service := NewInvitationService(mockRepo, nil, nil)

	invitationID := uuid.New()

	mockRepo.On("DeleteInvitation", invitationID, mock.Anything).Return(nil)

	err := service.DeleteInvitation(invitationID, context.Background())

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
