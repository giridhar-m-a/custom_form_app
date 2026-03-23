package dto

import (
	"github.com/giridhar-m-a/custom_form_app/internal/db/sqlc"
	"github.com/google/uuid"
)

type CreateInvitationDTO struct {
	FormID string `json:"form_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
}

type InvitationResponseDTO struct {
	InvitationID string `json:"invitation_id"`
	Email        string `json:"email"`
	Name         string `json:"name"`
	Status       string `json:"status"`
	InvitedAt    string `json:"invited_at"`
}

type UpdateInvitationDTO struct {
	Status sqlc.InvitationStatus `json:"status"`
}

type InvitationListQueryDto struct {
	FormId  string                  `form:"formId"`
	Exclude []sqlc.InvitationStatus `form:"exclude"`
	Status  sqlc.InvitationStatus   `form:"status"`
	Query
}

type InvitationListDto struct {
	Invitations []sqlc.Invitation `json:"invitations"`
	Total       int               `json:"total"`
	Page        int               `json:"page"`
	Limit       int               `json:"limit"`
	Pages       int               `json:"pages"`
}

type InvitationResponseDto struct {
	InvitationID string `json:"invitationId"`
	FormID       string `json:"formId"`
	InvitedEmail string `json:"invitedEmail"`
	Status       string `json:"status"`
	InvitedBy    string `json:"invitedBy"`
	InvitedName  string `json:"invitedName"`
}

type InvitationListResponse struct {
	Status     int                     `json:"status" example:"200"`
	Message    string                  `json:"message" example:"Invitations retrieved successfully"`
	Data       []InvitationResponseDto `json:"data"`
	Pagination PaginationResponse      `json:"pagination"`
}

type UpdateInvitationsResendParams struct {
	ResendID     uuid.UUID `json:"resend_id"`
	InvitationID uuid.UUID `json:"invitation_id"`
}

type InvitationEmailParams struct {
	PlatformName   string
	UserName       string
	Title          string
	InvitationURL  string
	Year           int
	CompanyAddress string
}
