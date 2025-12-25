package dto

import "github.com/giridhar-m-a/custom_form_app/internal/db/sqlc"

type CreateInvitationDTO struct {
	FormID    string `json:"form_id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
}

type InvitationResponseDTO struct {
	InvitationID string `json:"invitation_id"`
	Email        string `json:"email"`
	Name         string `json:"name"`
	Status       string `json:"status"`
	InvitedAt    string `json:"invited_at"`
}

type UpdateInvitationDTO struct {
	Status       sqlc.InvitationStatus `json:"status"`
}

type InvitationListQueryDto struct {
	FormId string `form:"formId"`
	Exclude sqlc.InvitationStatus `form:"exclude"`
	Status sqlc.InvitationStatus `form:"status"`
	Query
}