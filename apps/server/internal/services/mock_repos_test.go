package services

import (
	"context"
	"database/sql"

	"time"

	"github.com/giridhar-m-a/custom_form_app/internal/db/sqlc"
	"github.com/giridhar-m-a/custom_form_app/internal/dto"
	"github.com/giridhar-m-a/custom_form_app/internal/repositories"
	"github.com/google/uuid"
	"github.com/resend/resend-go/v3"
	"github.com/stretchr/testify/mock"
)

// Mock service for Users
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(ctx context.Context, userData map[string]any) (sqlc.User, error) {
	args := m.Called(ctx, userData)
	return args.Get(0).(sqlc.User), args.Error(1)
}

func (m *MockUserService) GetUserDetailsById(ctx context.Context, userID string) (sqlc.GetUserByIDRow, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(sqlc.GetUserByIDRow), args.Error(1)
}

func (m *MockUserService) GetUserDetailsByEmail(ctx context.Context, email string) (sqlc.GetUserByEmailRow, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(sqlc.GetUserByEmailRow), args.Error(1)
}

func (m *MockUserService) GetUserDetailsByGoogleId(ctx context.Context, googleID string) (sqlc.GetUserByGoogleIdRow, error) {
	args := m.Called(ctx, googleID)
	return args.Get(0).(sqlc.GetUserByGoogleIdRow), args.Error(1)
}

func (m *MockUserService) UpdateUser(ctx context.Context, userID string, data dto.UserUpdateDTO) (sqlc.User, error) {
	args := m.Called(ctx, userID, data)
	return args.Get(0).(sqlc.User), args.Error(1)
}

func (m *MockUserService) UpdateUserProfilePic(ctx context.Context, user string, data dto.FileUploadPayload) (sqlc.UpdateUserProfilePicRow, error) {
	args := m.Called(ctx, user, data)
	return args.Get(0).(sqlc.UpdateUserProfilePicRow), args.Error(1)
}

func (m *MockUserService) CreateUserProfilePic(ctx context.Context, user uuid.UUID, path string, size int64, fileType string) (sqlc.CreateUserProfilePicRow, error) {
	args := m.Called(ctx, user, path, size, fileType)
	return args.Get(0).(sqlc.CreateUserProfilePicRow), args.Error(1)
}

func (m *MockUserService) DeleteUser(ctx context.Context, user string) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserService) DeleteUserProfilePic(ctx context.Context, user string) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserService) GetUserPassword(ctx context.Context, userID string) (string, error) {
	args := m.Called(ctx, userID)
	return args.String(0), args.Error(1)
}

func (m *MockUserService) CreateTempUser(ctx context.Context, name string) (sqlc.User, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(sqlc.User), args.Error(1)
}

func (m *MockUserService) SoftDeleteUser(ctx context.Context, user string) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

// Mock service for JWT
type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) GenerateToken(userID string, expiresIn time.Duration, audience string) (string, error) {
	args := m.Called(userID, expiresIn, audience)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) GenerateInvitationToken(invitationID string, formID string, expiresIn time.Duration) (string, error) {
	args := m.Called(invitationID, formID, expiresIn)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) ValidateToken(token string) (string, error) {
	args := m.Called(token)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) ValidateInvitationToken(token string) (*InvitationClaims, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*InvitationClaims), args.Error(1)
}

func (m *MockJWTService) GenerateAnonymousInvitationToken(formId string, expiresIn time.Duration) (string, error) {
	args := m.Called(formId, expiresIn)
	return args.String(0), args.Error(1)
}

// Mock repos for Forms
type MockFormsRepository struct {
	mock.Mock
}

func (m *MockFormsRepository) CreateForm(params sqlc.CreateFormParams, ctx context.Context) (sqlc.Form, error) {
	args := m.Called(params, ctx)
	return args.Get(0).(sqlc.Form), args.Error(1)
}

func (m *MockFormsRepository) GetFormsList(params sqlc.ListFormsParams, ctx context.Context) ([]sqlc.ListFormsRow, error) {
	args := m.Called(params, ctx)
	return args.Get(0).([]sqlc.ListFormsRow), args.Error(1)
}

func (m *MockFormsRepository) GetFormByID(formID string, ctx context.Context) (sqlc.Form, error) {
	args := m.Called(formID, ctx)
	return args.Get(0).(sqlc.Form), args.Error(1)
}

func (m *MockFormsRepository) UpdateForm(params sqlc.UpdateFormParams, ctx context.Context) (sqlc.Form, error) {
	args := m.Called(params, ctx)
	return args.Get(0).(sqlc.Form), args.Error(1)
}

func (m *MockFormsRepository) DeleteForm(formID string, ctx context.Context) (sqlc.DeleteFormRow, error) {
	args := m.Called(formID, ctx)
	return args.Get(0).(sqlc.DeleteFormRow), args.Error(1)
}

func (m *MockFormsRepository) SoftDeleteForm(formID uuid.UUID, ctx context.Context) error {
	args := m.Called(formID, ctx)
	return args.Error(0)
}

func (m *MockFormsRepository) FormRepoWithTx(tx *sql.Tx) repositories.FormsRepository {
	args := m.Called(tx)
	return args.Get(0).(repositories.FormsRepository)
}

// Mock repo for Form Fields
type MockFormFieldsRepository struct {
	mock.Mock
}

func (m *MockFormFieldsRepository) CreateFormField(params sqlc.CreateFormFieldParams, ctx context.Context) (sqlc.CreateFormFieldRow, error) {
	args := m.Called(params, ctx)
	return args.Get(0).(sqlc.CreateFormFieldRow), args.Error(1)
}

func (m *MockFormFieldsRepository) UpdateFormField(params sqlc.UpdateFormFieldParams, ctx context.Context) (sqlc.UpdateFormFieldRow, error) {
	args := m.Called(params, ctx)
	return args.Get(0).(sqlc.UpdateFormFieldRow), args.Error(1)
}

func (m *MockFormFieldsRepository) DeleteFormField(fieldId string, ctx context.Context) (sqlc.DeleteFormFieldRow, error) {
	args := m.Called(fieldId, ctx)
	return args.Get(0).(sqlc.DeleteFormFieldRow), args.Error(1)
}

func (m *MockFormFieldsRepository) FormFieldRepoWithTx(tx *sql.Tx) repositories.FormFieldsRepository {
	args := m.Called(tx)
	return args.Get(0).(repositories.FormFieldsRepository)
}

func (m *MockFormFieldsRepository) GetFormFieldsByFormId(formId string, ctx context.Context) ([]sqlc.GetFormFieldsWithOptionsRow, error) {
	args := m.Called(formId, ctx)
	return args.Get(0).([]sqlc.GetFormFieldsWithOptionsRow), args.Error(1)
}

// Mock repo for Form Field Options
type MockFormFieldOptionsRepository struct {
	mock.Mock
}

func (m *MockFormFieldOptionsRepository) CreateFieldOption(params sqlc.CreateFieldOptionParams, ctx context.Context) (sqlc.CreateFieldOptionRow, error) {
	args := m.Called(params, ctx)
	return args.Get(0).(sqlc.CreateFieldOptionRow), args.Error(1)
}

func (m *MockFormFieldOptionsRepository) UpdateFieldOption(params sqlc.UpdateFieldOptionParams, ctx context.Context) (sqlc.UpdateFieldOptionRow, error) {
	args := m.Called(params, ctx)
	return args.Get(0).(sqlc.UpdateFieldOptionRow), args.Error(1)
}

func (m *MockFormFieldOptionsRepository) DeleteFieldOption(optionId string, ctx context.Context) (sqlc.DeleteFieldOptionRow, error) {
	args := m.Called(optionId, ctx)
	return args.Get(0).(sqlc.DeleteFieldOptionRow), args.Error(1)
}

func (m *MockFormFieldOptionsRepository) FormFieldOptionsRepoWithTx(tx *sql.Tx) repositories.FormFieldOptionsRepository {
	args := m.Called(tx)
	return args.Get(0).(repositories.FormFieldOptionsRepository)
}

// Mock repo for Responses
type MockResponseRepository struct {
	mock.Mock
}

func (m *MockResponseRepository) CreateSubmission(ctx context.Context, arg sqlc.CreateSubmissionParams) (sqlc.FormSubmission, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(sqlc.FormSubmission), args.Error(1)
}

func (m *MockResponseRepository) CreateResponse(ctx context.Context, arg sqlc.CreateResponseParams) (sqlc.FormResponse, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(sqlc.FormResponse), args.Error(1)
}

func (m *MockResponseRepository) CreateResponseOption(ctx context.Context, arg sqlc.CreateResponseOptionParams) (sqlc.ResponseOption, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(sqlc.ResponseOption), args.Error(1)
}

func (m *MockResponseRepository) CreateResponseFiles(ctx context.Context, arg sqlc.CreateResponseFilesParams) (sqlc.FormResponseFile, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(sqlc.FormResponseFile), args.Error(1)
}

func (m *MockResponseRepository) GetSubmissionById(ctx context.Context, submissionID uuid.UUID) (repositories.SubmissionResponse, error) {
	args := m.Called(ctx, submissionID)
	return args.Get(0).(repositories.SubmissionResponse), args.Error(1)
}

func (m *MockResponseRepository) GetSubmissions(ctx context.Context, arg sqlc.GetSubmissionsParams) ([]sqlc.GetSubmissionsRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]sqlc.GetSubmissionsRow), args.Error(1)
}

func (m *MockResponseRepository) GetSubmissionCount(ctx context.Context, arg sqlc.GetSubmissionCountParams) (int64, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockResponseRepository) ResponseRepositoryWithTx(tx *sql.Tx) repositories.ResponseRepository {
	args := m.Called(tx)
	return args.Get(0).(repositories.ResponseRepository)
}

// Mock repo for Invitations
type MockInvitationRepository struct {
	mock.Mock
}

func (m *MockInvitationRepository) CreateInvitation(invitations sqlc.CreateManyInvitationsParams, ctx context.Context) ([]sqlc.CreateManyInvitationsRow, error) {
	args := m.Called(invitations, ctx)
	return args.Get(0).([]sqlc.CreateManyInvitationsRow), args.Error(1)
}

func (m *MockInvitationRepository) UpdateInvitationStatus(status sqlc.UpdateInvitationStatusParams, ctx context.Context) (sqlc.UpdateInvitationStatusRow, error) {
	args := m.Called(status, ctx)
	return args.Get(0).(sqlc.UpdateInvitationStatusRow), args.Error(1)
}

func (m *MockInvitationRepository) DeleteInvitation(invitationID uuid.UUID, ctx context.Context) error {
	args := m.Called(invitationID, ctx)
	return args.Error(0)
}

func (m *MockInvitationRepository) GetInvitationByFormId(query sqlc.GetInvitationByFormIdParams, ctx context.Context) ([]sqlc.Invitation, error) {
	args := m.Called(query, ctx)
	return args.Get(0).([]sqlc.Invitation), args.Error(1)
}

func (m *MockInvitationRepository) CreateSingleInvitation(invitation sqlc.CreateInvitationParams, ctx context.Context) (sqlc.CreateInvitationRow, error) {
	args := m.Called(invitation, ctx)
	return args.Get(0).(sqlc.CreateInvitationRow), args.Error(1)
}

func (m *MockInvitationRepository) InvitationRepositoryWithTx(tx *sql.Tx) repositories.InvitationRepository {
	args := m.Called(tx)
	return args.Get(0).(repositories.InvitationRepository)
}

func (m *MockInvitationRepository) CountInvitationsByFormId(params sqlc.CountInvitationsByFormIdParams, ctx context.Context) (int64, error) {
	args := m.Called(params, ctx)
	return args.Get(0).(int64), args.Error(1)
}

// Mock service for Forms
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

func (m *MockFormService) SoftDeleteForm(formID string, ctx context.Context) error {
	args := m.Called(formID, ctx)
	return args.Error(0)
}

func (m *MockFormService) UpdateFormScheduleId(formID uuid.UUID, scheduleID uuid.NullUUID, invitationId uuid.NullUUID, ctx context.Context) (sqlc.Form, error) {
	args := m.Called(formID, scheduleID, invitationId, ctx)
	return args.Get(0).(sqlc.Form), args.Error(1)
}

// Mock repo for Dashboard
type MockDashboardRepository struct {
	mock.Mock
}

func (m *MockDashboardRepository) GetTotalForms(ctx context.Context, userID uuid.NullUUID) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockDashboardRepository) GetTotalSubmissions(ctx context.Context, userID uuid.NullUUID) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockDashboardRepository) GetTotalActiveForms(ctx context.Context, userID uuid.NullUUID) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockDashboardRepository) GetTotalClosedForms(ctx context.Context, userID uuid.NullUUID) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockDashboardRepository) GetTotalInvitations(ctx context.Context, userID uuid.NullUUID) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockDashboardRepository) GetFormSubmissionsByMonth(ctx context.Context, userID uuid.NullUUID) ([]sqlc.GetFormSubmissionsByMonthRow, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]sqlc.GetFormSubmissionsByMonthRow), args.Error(1)
}

// Mock service for Google Auth
type MockGoogleAuthService struct {
	mock.Mock
}

func (m *MockGoogleAuthService) Authenticate(ctx context.Context, code string) (sqlc.GetUserByGoogleIdRow, error) {
	args := m.Called(ctx, code)
	return args.Get(0).(sqlc.GetUserByGoogleIdRow), args.Error(1)
}

// Mock service for Bcrypt
type MockBcryptService struct {
	mock.Mock
}

func (m *MockBcryptService) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockBcryptService) ComparePassword(hashedPassword, password string) bool {
	args := m.Called(hashedPassword, password)
	return args.Bool(0)
}

// Mock service for Mail
type MockMailService struct {
	mock.Mock
}

func (m *MockMailService) SendEmail(params resend.SendEmailRequest) (*resend.SendEmailResponse, error) {
	args := m.Called(params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*resend.SendEmailResponse), args.Error(1)
}

func (m *MockMailService) SendBulk(ctx context.Context, params []*resend.SendEmailRequest) ([]resend.SendEmailResponse, error) {
	args := m.Called(ctx, params)
	return args.Get(0).([]resend.SendEmailResponse), args.Error(1)
}
