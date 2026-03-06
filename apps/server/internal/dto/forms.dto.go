package dto

import (
	"time"

	"github.com/giridhar-m-a/custom_form_app/internal/db/sqlc"
)

type CreateFormDTO struct {
	Title                 string          `json:"title" binding:"required" form:"title" message:"title is required"`
	Description           *string         `json:"description,omitempty" form:"description"`
	FormAccess            sqlc.FormAccess `json:"form_access,omitempty" form:"form_access"` // default 'restricted'
	ScheduledTime         *time.Time      `json:"scheduledTime,omitempty" form:"scheduled_time"`
	ClosingTime           *time.Time      `json:"closingTime,omitempty" form:"closing_time"`
	IsScheduled           *bool           `json:"isScheduled,omitempty" form:"is_scheduled"`
	InvitationScheduleGap *int32          `json:"invitationScheduleGap,omitempty" form:"invitation_schedule_gap"`
}

type FormResponse struct {
	ID                  string `json:"id"`
	Title               string `json:"title"`
	Description         string `json:"description,omitempty"`
	CreatedBy           string `json:"createdBy,omitempty"`
	Status              string `json:"status,omitempty"`
	CreatedAt           string `json:"createdAt,omitempty"`
	UpdatedAt           string `json:"updatedAt,omitempty"`
	Access              string `json:"access,omitempty"`
	SchedulingID        string `json:"schedulingId,omitempty"`
	ScheduledTime       string `json:"scheduledTime,omitempty"`
	ClosingTime         string `json:"closingTime,omitempty"`
	IsScheduleCompleted bool   `json:"isScheduleCompleted,omitempty"`
	IsScheduled         bool   `json:"isScheduled,omitempty"`
}

type ListFormQuery struct {
	Query
	Status string `form:"status"`
	Access string `form:"access"`
}

type FormListResponse struct {
	Forms []sqlc.ListFormsRow `json:"forms"`
	Total int                 `json:"total"`
	Page  int                 `json:"page"`
	Limit int                 `json:"limit"`
	Pages int                 `json:"pages"`
}

type UpdateFormDTO struct {
	Title                 *string          `json:"title,omitempty" form:"title"`
	Description           *string          `json:"description,omitempty" form:"description"`
	Access                *sqlc.FormAccess `json:"access,omitempty" form:"access"`
	Status                *sqlc.FormStatus `json:"status,omitempty" form:"status"`
	SchedulingID          *string          `json:"schedulingId,omitempty" form:"scheduling_id"`
	ScheduledTime         *time.Time       `json:"scheduledTime,omitempty" form:"scheduled_time"`
	ClosingTime           *time.Time       `json:"closingTime,omitempty" form:"closing_time"`
	IsScheduleCompleted   *bool            `json:"isScheduleCompleted,omitempty" form:"is_schedule_completed"`
	IsScheduled           *bool            `json:"isScheduled,omitempty" form:"is_scheduled"`
	InvitationScheduleGap *int32           `json:"invitationScheduleGap,omitempty" form:"invitation_schedule_gap"`
}
