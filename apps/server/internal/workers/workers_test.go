package workers

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/giridhar-m-a/custom_form_app/internal/db/sqlc"
	"github.com/giridhar-m-a/custom_form_app/internal/dto/scheduler_dto"
	"github.com/giridhar-m-a/custom_form_app/internal/repositories"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockFormsRepository struct {
	mock.Mock
}

func (m *MockFormsRepository) CreateForm(form sqlc.CreateFormParams, ctx context.Context) (sqlc.Form, error) {
	args := m.Called(form, ctx)
	return args.Get(0).(sqlc.Form), args.Error(1)
}

func (m *MockFormsRepository) UpdateForm(form sqlc.UpdateFormParams, ctx context.Context) (sqlc.Form, error) {
	args := m.Called(form, ctx)
	return args.Get(0).(sqlc.Form), args.Error(1)
}

func (m *MockFormsRepository) GetFormByID(id string, ctx context.Context) (sqlc.Form, error) {
	args := m.Called(id, ctx)
	return args.Get(0).(sqlc.Form), args.Error(1)
}

func (m *MockFormsRepository) GetFormsList(params sqlc.ListFormsParams, ctx context.Context) ([]sqlc.ListFormsRow, error) {
	args := m.Called(params, ctx)
	return args.Get(0).([]sqlc.ListFormsRow), args.Error(1)
}

func (m *MockFormsRepository) DeleteForm(id string, ctx context.Context) (sqlc.DeleteFormRow, error) {
	args := m.Called(id, ctx)
	return args.Get(0).(sqlc.DeleteFormRow), args.Error(1)
}

func (m *MockFormsRepository) FormRepoWithTx(tx *sql.Tx) repositories.FormsRepository {
	return m
}

func (m *MockFormsRepository) SoftDeleteForm(id uuid.UUID, ctx context.Context) error {
	args := m.Called(id, ctx)
	return args.Error(0)
}

func TestFormWorker_HandleFormStatusUpdate(t *testing.T) {
	mockRepo := new(MockFormsRepository)
	worker := NewFormWorker(mockRepo)

	formID := uuid.New()
	payload := scheduler_dto.InvitationSchedulerPayload{
		FormID: formID.String(),
	}
	payloadBytes, _ := json.Marshal(payload)
	task := asynq.NewTask("update_status", payloadBytes)

	mockRepo.On("UpdateForm", mock.MatchedBy(func(p sqlc.UpdateFormParams) bool {
		return p.FormID == formID && p.FormStatus.FormStatus == sqlc.FormStatusPublished
	}), mock.Anything).Return(sqlc.Form{}, nil)

	handler := worker.HandleFormStatusUpdate()
	err := handler(context.Background(), task)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
