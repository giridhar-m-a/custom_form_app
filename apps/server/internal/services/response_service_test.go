package services

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/giridhar-m-a/custom_form_app/internal/db/sqlc"
	"github.com/giridhar-m-a/custom_form_app/internal/dto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSubmitForm(t *testing.T) {
	db, mockDB, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockRepo := new(MockResponseRepository)
	service := NewResponseService(mockRepo, db)

	formID := uuid.New()
	respondentID := uuid.New().String()
	submissionReq := dto.CreateSubmissionRequest{
		FormID:       formID.String(),
		RespondentID: &respondentID,
		Responses: []dto.ResponseRequest{
			{
				FormFieldID: uuid.New().String(),
				ResponseText: func(s string) *string { return &s }("Answer"),
			},
		},
	}

	mockDB.ExpectBegin()
	mockRepo.On("ResponseRepositoryWithTx", mock.Anything).Return(mockRepo)
	
	submission := sqlc.FormSubmission{
		SubmissionID: uuid.New(),
		FormID:       formID,
	}
	mockRepo.On("CreateSubmission", mock.Anything, mock.Anything).Return(submission, nil)
	
	response := sqlc.FormResponse{
		ResponseID: uuid.New(),
		FormFieldID: uuid.MustParse(submissionReq.Responses[0].FormFieldID),
	}
	mockRepo.On("CreateResponse", mock.Anything, mock.Anything).Return(response, nil)
	
	mockDB.ExpectCommit()

	res, err := service.SubmitForm(context.Background(), submissionReq)

	assert.NoError(t, err)
	assert.Equal(t, submission.SubmissionID.String(), res.SubmissionID)
	mockRepo.AssertExpectations(t)
	assert.NoError(t, mockDB.ExpectationsWereMet())
}

func TestGetSubmissions(t *testing.T) {
	mockRepo := new(MockResponseRepository)
	service := NewResponseService(mockRepo, nil)

	formID := uuid.New()
	query := dto.GetSubmissionsRequest{
		FormID: formID.String(),
		ResponseQuery: dto.ResponseQuery{
			Page: 1,
			Limit: 10,
		},
	}

	expectedSubmissions := []sqlc.GetSubmissionsRow{
		{
			SubmissionID: uuid.New(),
		},
	}

	mockRepo.On("GetSubmissions", mock.Anything, mock.Anything).Return(expectedSubmissions, nil)
	mockRepo.On("GetSubmissionCount", mock.Anything, mock.Anything).Return(int64(1), nil)

	res, err := service.GetSubmissions(context.Background(), query)

	assert.NoError(t, err)
	assert.Len(t, res.Submissions, 1)
	mockRepo.AssertExpectations(t)
}
