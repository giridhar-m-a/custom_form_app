package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/giridhar-m-a/custom_form_app/internal/db/sqlc"
	"github.com/google/uuid"
)

type SubmissionResponseFile struct {
	ResponseFileID string    `json:"response_file_id"`
	FileName       string    `json:"file_name"`
	FilePath       string    `json:"file_path"`
	FileSize       int64     `json:"file_size"`
	FileType       string    `json:"file_type"`
	FileUploadedAt time.Time `json:"file_uploaded_at"`
}

type SubmissionResponseOption struct {
	ID           string `json:"id"`             // maps to ro.id
	FormOptionID string `json:"form_option_id"` // maps to ro.form_option_id
}

type SubmissionResponseItem struct {
	ResponseID        string                     `json:"response_id"`
	FormFieldID       string                     `json:"form_field_id"`
	ResponseText      *string                    `json:"response_text"`
	FormFieldOptions  []SubmissionResponseOption `json:"form_field_options"`
	FormResponseFiles []SubmissionResponseFile   `json:"form_response_files"`
}

type SubmissionResponse struct {
	SubmissionID string                   `json:"submission_id"`
	FormID       string                   `json:"form_id"`
	RespondentID string                   `json:"respondent_id"`
	SubmittedAt  time.Time                `json:"submitted_at"`
	Responses    []SubmissionResponseItem `json:"responses"`
}

type ResponseRepository interface {
	CreateSubmission(ctx context.Context, arg sqlc.CreateSubmissionParams) (sqlc.FormSubmission, error)
	CreateResponse(ctx context.Context, arg sqlc.CreateResponseParams) (sqlc.FormResponse, error)
	CreateResponseOption(ctx context.Context, arg sqlc.CreateResponseOptionParams) (sqlc.ResponseOption, error)
	CreateResponseFiles(ctx context.Context, arg sqlc.CreateResponseFilesParams) (sqlc.FormResponseFile, error)
	GetSubmissionById(ctx context.Context, submissionID uuid.UUID) (SubmissionResponse, error)
	GetSubmissions(ctx context.Context, arg sqlc.GetSubmissionsParams) ([]sqlc.GetSubmissionsRow, error)
	GetSubmissionCount(ctx context.Context, arg sqlc.GetSubmissionCountParams) (int64, error)
	ResponseRepositoryWithTx(tx *sql.Tx) ResponseRepository
}

type responseRepository struct {
	queries *sqlc.Queries
}

func NewResponseRepository(queries *sqlc.Queries) ResponseRepository {
	return &responseRepository{queries: queries}
}

func (r *responseRepository) CreateSubmission(ctx context.Context, arg sqlc.CreateSubmissionParams) (sqlc.FormSubmission, error) {
	return r.queries.CreateSubmission(ctx, arg)
}

func (r *responseRepository) CreateResponse(ctx context.Context, arg sqlc.CreateResponseParams) (sqlc.FormResponse, error) {
	return r.queries.CreateResponse(ctx, arg)
}

func (r *responseRepository) CreateResponseOption(ctx context.Context, arg sqlc.CreateResponseOptionParams) (sqlc.ResponseOption, error) {
	return r.queries.CreateResponseOption(ctx, arg)
}

func (r *responseRepository) CreateResponseFiles(ctx context.Context, arg sqlc.CreateResponseFilesParams) (sqlc.FormResponseFile, error) {
	return r.queries.CreateResponseFiles(ctx, arg)
}

func (r *responseRepository) GetSubmissionById(ctx context.Context, submissionID uuid.UUID) (SubmissionResponse, error) {
	submission, err := r.queries.GetSubmissionById(ctx, submissionID)
	data, err := json.Marshal(submission)
	if err != nil {
		return SubmissionResponse{}, err
	}

	var response SubmissionResponse
	if err := json.Unmarshal(data, &response); err != nil {
		slog.Error("[response service] failed to unmarshal submission", "error", err.Error())
		return SubmissionResponse{}, err
	}
	return response, nil
}

func (r *responseRepository) GetSubmissions(ctx context.Context, arg sqlc.GetSubmissionsParams) ([]sqlc.GetSubmissionsRow, error) {
	return r.queries.GetSubmissions(ctx, arg)
}

func (r *responseRepository) GetSubmissionCount(ctx context.Context, arg sqlc.GetSubmissionCountParams) (int64, error) {
	return r.queries.GetSubmissionCount(ctx, arg)
}

func (r *responseRepository) ResponseRepositoryWithTx(tx *sql.Tx) ResponseRepository {
	return &responseRepository{queries: r.queries.WithTx(tx)}
}
