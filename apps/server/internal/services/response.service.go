package services

import (
	"context"
	"database/sql"
	"log"
	"log/slog"

	"github.com/giridhar-m-a/custom_form_app/internal/db/sqlc"
	"github.com/giridhar-m-a/custom_form_app/internal/dto"
	"github.com/giridhar-m-a/custom_form_app/internal/repositories"
	"github.com/giridhar-m-a/custom_form_app/internal/utils"
)

type ResponseService interface {
	SubmitForm(ctx context.Context, submission dto.CreateSubmissionRequest) (dto.SubmissionResponse, error)
	GetSubmissions(ctx context.Context, query dto.GetSubmissionsRequest) (dto.GetSubmissionsResponse, error)
	GetSingleSubmission(ctx context.Context, submissionID string) (dto.SubmissionResponse, error)
}

type responseService struct {
	responseRepo repositories.ResponseRepository
	db           *sql.DB
}

func NewResponseService(responseRepo repositories.ResponseRepository, db *sql.DB) ResponseService {
	return &responseService{responseRepo: responseRepo, db: db}
}

func (s *responseService) SubmitForm(ctx context.Context, submission dto.CreateSubmissionRequest) (dto.SubmissionResponse, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		slog.Error("[response service] failed to begin transaction:", "Error", err)
		return dto.SubmissionResponse{}, err
	}
	defer tx.Rollback()
	repo := s.responseRepo.ResponseRepositoryWithTx(tx)
	formId, err := utils.ConvertStringToUUID(submission.FormID)
	if err != nil {
		slog.Error("[response service] failed to convert string to uuid: ", "Error", err)
		return dto.SubmissionResponse{}, err
	}
	respondent := utils.ConvertStringToNullUUID(*submission.RespondentID)

	newSubmission, err := repo.CreateSubmission(ctx, sqlc.CreateSubmissionParams{
		FormID:       formId,
		RespondentID: respondent,
	})
	if err != nil {
		slog.Error("[response service] failed to create submission", "Error", err)
		return dto.SubmissionResponse{}, err
	}
	slog.Info("[response service] submission created successfully, Continue creating responses")
	responses := make([]dto.SubmissionResponseField, 0)
	for i, response := range submission.Responses {
		formFieldId, err := utils.ConvertStringToUUID(response.FormFieldID)
		if err != nil {
			slog.Error("[response service] [index: %d] failed to convert string to uuid: %v", i, err)
			return dto.SubmissionResponse{}, err
		}
		responseText := utils.ConvertStringToNullString(*response.ResponseText)
		responseParams := sqlc.CreateResponseParams{
			SubmissionID: newSubmission.SubmissionID,
			FormFieldID:  formFieldId,
			ResponseText: responseText,
		}
		newResponse, err := repo.CreateResponse(ctx, responseParams)
		if err != nil {
			slog.Error("[response service] [index: %d] failed to create response: %v", i, err)
			return dto.SubmissionResponse{}, err
		}
		responseOptions := make([]dto.SubmissionResponseOption, 0)
		slog.Info("[response service] [index: %d] response created successfully, Continue creating response options", i)
		for j, option := range response.ResponseOptions {
			optionId, err := utils.ConvertStringToUUID(option.OptionID)
			if err != nil {
				slog.Error("[response service] [index: %d] [option index: %d] failed to convert string to uuid: %v", i, j, err)
				return dto.SubmissionResponse{}, err
			}
			responseOptionParams := sqlc.CreateResponseOptionParams{
				ResponseID:   newResponse.ResponseID,
				FormOptionID: optionId,
			}
			newResponseOption, err := repo.CreateResponseOption(ctx, responseOptionParams)
			if err != nil {
				slog.Error("[response service] [index: %d] [option index: %d] failed to create response option: %v", i, j, err)
				return dto.SubmissionResponse{}, err
			}
			slog.Info("[response service] [index: %d] [option index: %d] response option created successfully, Continue creating response files", i, j)
			responseOptions = append(responseOptions, dto.SubmissionResponseOption{
				ID:           newResponseOption.ID.String(),
				FormOptionID: newResponseOption.FormOptionID.String(),
			})
		}
		slog.Info("[response service] [index: %d] response options created successfully, Continue creating response files", i)
		responseFiles := make([]dto.SubmissionResponseFile, 0)
		for k, file := range response.ResponseFiles {
			param := sqlc.CreateResponseFilesParams{
				ResponseID: newResponse.ResponseID,
				FileName:   file.FileName,
				FilePath:   file.FilePath,
				FileSize:   file.FileSize,
				FileType:   file.FileType,
				FormID:     formId,
			}
			newResponseFile, err := repo.CreateResponseFiles(ctx, param)
			if err != nil {
				slog.Error("[response service] [index: %d] [file index: %d] failed to create response file: %v", i, k, err)
				return dto.SubmissionResponse{}, err
			}
			slog.Info("[response service] [index: %d] [file index: %d] response file created successfully, Continue creating response files", i, k)
			responseFiles = append(responseFiles, dto.SubmissionResponseFile{
				ResponseFileID: newResponseFile.ResponseFileID.String(),
				FileName:       newResponseFile.FileName,
				FilePath:       newResponseFile.FilePath,
				FileSize:       newResponseFile.FileSize,
				FileType:       newResponseFile.FileType,
			})
		}
		slog.Info("[response service] [index: %d] response files created successfully, Continue creating response files", i)
		responses = append(responses, dto.SubmissionResponseField{
			ResponseID:      newResponse.ResponseID.String(),
			FormFieldID:     newResponse.FormFieldID.String(),
			ResponseText:    newResponse.ResponseText.String,
			ResponseOptions: responseOptions,
			ResponseFiles:   responseFiles,
		})
		slog.Info("[response service] [index: %d] response created successfully, Continue creating response options", i)
	}
	if err := tx.Commit(); err != nil {
		slog.Error("[response service] failed to commit transaction:", "Error",err)
		return dto.SubmissionResponse{}, err
	}
	slog.Info("[response service] transaction committed successfully")
	return dto.SubmissionResponse{
		SubmissionID: newSubmission.SubmissionID.String(),
		FormID:       newSubmission.FormID.String(),
		RespondentID: newSubmission.RespondentID.UUID.String(),
		Responses:    responses,
	}, nil
}

func (s *responseService) GetSubmissions(ctx context.Context, query dto.GetSubmissionsRequest) (dto.GetSubmissionsResponse, error) {
	formId, err := utils.ConvertStringToUUID(query.FormID)
	if err != nil {
		slog.Error("[response service] failed to convert string to uuid: ", "error", err)
		return dto.GetSubmissionsResponse{}, err
	}
	search := utils.ConvertStringToNullString(query.Search)
	limit := int32(query.Limit)
	if limit <= 0 {
		limit = 10
	}
	page := int32(query.Page)
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit
	submissions, err := s.responseRepo.GetSubmissions(ctx, sqlc.GetSubmissionsParams{
		FormID:      formId,
		Search:      search,
		OffsetCount: offset,
		LimitCount:  limit,
	})
	count, err := s.responseRepo.GetSubmissionCount(ctx, sqlc.GetSubmissionCountParams{
		FormID: formId,
		Search: search,
	})
	if err != nil {
		slog.Error("[response service] failed to get submissions:", "error", err)
		return dto.GetSubmissionsResponse{}, err
	}
	slog.Info("[response service] submissions fetched successfully")
	total := int(count)

	responseSubmissions := make([]dto.SubmissionList, 0)
	for _, submission := range submissions {
		responseSubmissions = append(responseSubmissions, dto.SubmissionList{
			SubmissionID: submission.SubmissionID.String(),
			SubmittedAt:  submission.SubmittedAt.Time.String(),
			RespondentID: submission.RespondentID.UUID.String(),
			InvitedEmail: submission.InvitedEmail.String,
			InvitedName: submission.InvitedName.String,
		})
	}

	totalPages := total / int(limit)
	if total%int(limit) != 0 {
		totalPages++
	}
	if totalPages == 0 {
		totalPages = 1
	}
	log.Printf("[response service] submissions fetched successfully: %v", responseSubmissions)

	return dto.GetSubmissionsResponse{
		Submissions:  responseSubmissions,
		Total:        total,
		Page:         int(page),
		Limit:        int(limit),
		Pages:        totalPages,
		TotalRecords: total,
	}, nil
}

func (s *responseService) GetSingleSubmission(ctx context.Context, submissionID string) (dto.SubmissionResponse, error) {
	submissionId, err := utils.ConvertStringToUUID(submissionID)
	if err != nil {
		slog.Error("[response service] failed to convert string to uuid", "error", err.Error())
		return dto.SubmissionResponse{}, err
	}
	submission, err := s.responseRepo.GetSubmissionById(ctx, submissionId)
	log.Printf("submission: %v", submission)
	if err != nil {
		slog.Error("[response service] failed to get single submission", "error", err.Error())
		return dto.SubmissionResponse{}, err
	}
	responseFields := make([]dto.SubmissionResponseField, 0)
	for _, response := range submission.Responses {
		responseOptions := make([]dto.SubmissionResponseOption, 0)
		responseFiles := make([]dto.SubmissionResponseFile, 0)
		for _, option := range response.FormFieldOptions {
			responseOptions = append(responseOptions, dto.SubmissionResponseOption{
				ID:           option.ID,
				FormOptionID: option.FormOptionID,
			})
		}
		for _, file := range response.FormResponseFiles {
			responseFiles = append(responseFiles, dto.SubmissionResponseFile{
				ResponseFileID: file.ResponseFileID,
				FileName:       file.FileName,
				FilePath:       file.FilePath,
				FileSize:       file.FileSize,
				FileType:       file.FileType,
			})
		}
		var responseText string
		if response.ResponseText != nil {
			responseText = *response.ResponseText
		}
		responseFields = append(responseFields, dto.SubmissionResponseField{
			ResponseID:      response.ResponseID,
			FormFieldID:     response.FormFieldID,
			ResponseText:    responseText,
			ResponseOptions: responseOptions,
			ResponseFiles:   responseFiles,
		})
	}
	return dto.SubmissionResponse{
		SubmissionID: submission.SubmissionID,
		FormID:       submission.FormID,
		RespondentID: submission.RespondentID,
		Responses:    responseFields,
	}, nil
}
