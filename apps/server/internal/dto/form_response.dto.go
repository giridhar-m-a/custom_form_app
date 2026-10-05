package dto

type ResponseFileRequest struct {
	FileName string `json:"fileName" validate:"required"`
	FilePath string `json:"filePath" validate:"required"`
	FileSize int64  `json:"fileSize" validate:"required"`
	FileType string `json:"fileType" validate:"required"`
}

type ResponseOptionRequest struct {
	OptionID string `json:"optionId" validate:"required,uuid"`
}

type ResponseRequest struct {
	FormFieldID     string                  `json:"formFieldId" validate:"required,uuid"`
	ResponseText    *string                 `json:"responseText" validate:"omitempty"`
	ResponseOptions []ResponseOptionRequest `json:"responseOptions" validate:"omitempty"`
	ResponseFiles   []ResponseFileRequest   `json:"responseFiles" validate:"omitempty"`
}

type CreateSubmissionRequest struct {
	FormID       string            `json:"formId" validate:"required,uuid"`
	RespondentID *string           `json:"respondentId" validate:"omitempty,uuid"`
	Responses    []ResponseRequest `json:"responses" validate:"required"`
}

type SubmissionResponseFile struct {
	ResponseFileID string `json:"fileId" validate:"required,uuid"`
	FileName       string `json:"fileName" validate:"required"`
	FilePath       string `json:"filePath" validate:"required"`
	FileSize       int64  `json:"fileSize" validate:"required"`
	FileType       string `json:"fileType" validate:"required"`
	FileUploadedAt string `json:"fileUploadedAt" validate:"required"`
}

type SubmissionResponseOption struct {
	ID           string `json:"id" validate:"required,uuid"`
	FormOptionID string `json:"formOptionId" validate:"required,uuid"`
}

type SubmissionResponseField struct {
	ResponseID      string                     `json:"responseId" validate:"required,uuid"`
	FormFieldID     string                     `json:"formFieldId" validate:"required,uuid"`
	ResponseText    string                     `json:"responseText" validate:"required"`
	ResponseOptions []SubmissionResponseOption `json:"responseOptions" validate:"required"`
	ResponseFiles   []SubmissionResponseFile   `json:"responseFiles" validate:"required"`
}

type SubmissionResponse struct {
	SubmissionID string                    `json:"submissionId" validate:"required,uuid"`
	FormID       string                    `json:"formId" validate:"required,uuid"`
	SubmittedAt  string                    `json:"submittedAt" validate:"required"`
	RespondentID string                    `json:"respondentId" validate:"required,uuid"`
	Responses    []SubmissionResponseField `json:"responses" validate:"required"`
}

type GetSubmissionsRequest struct {
	FormID string `json:"formId" validate:"required,uuid"`
	ResponseQuery
}

type SubmissionList struct {
	SubmissionID string `json:"submissionId" validate:"required,uuid"`
	SubmittedAt  string `json:"submittedAt" validate:"required"`
	RespondentID string `json:"respondentId" validate:"required,uuid"`
	InvitedEmail string `json:"invitedEmail" validate:"required,email"`
	InvitedName  string `json:"invitedName" validate:"required"`
}

type GetSubmissionsResponse struct {
	Submissions  []SubmissionList `json:"submissions" validate:"required"`
	Total        int              `json:"total" validate:"required"`
	Page         int              `json:"page" validate:"required"`
	Limit        int              `json:"limit" validate:"required"`
	Pages        int              `json:"pages" validate:"required"`
	TotalRecords int              `json:"totalRecords" validate:"required"`
}
