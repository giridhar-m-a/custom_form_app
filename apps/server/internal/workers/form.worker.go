package workers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/giridhar-m-a/custom_form_app/internal/db/sqlc"
	"github.com/giridhar-m-a/custom_form_app/internal/dto"
	"github.com/giridhar-m-a/custom_form_app/internal/dto/scheduler_dto"
	"github.com/giridhar-m-a/custom_form_app/internal/services"
	"github.com/giridhar-m-a/custom_form_app/internal/utils"
	"github.com/hibiken/asynq"
)

type FormWorker struct {
	formService services.FormService
}

func NewFormWorker(service services.FormService) *FormWorker {
	return &FormWorker{
		formService: service,
	}
}

func (w *FormWorker) HandleFormStatusUpdate() asynq.HandlerFunc {
	return func(ctx context.Context, task *asynq.Task) error {
		var payload scheduler_dto.InvitationSchedulerPayload
		if err := json.Unmarshal(task.Payload(), &payload); err != nil {
			fmt.Printf("cannot unmarshal task payload: %v", err)
			return err
		}

		formId, err := utils.ConvertStringToUUID(payload.FormID)
		if err != nil {
			fmt.Printf("cannot convert string to uuid: %v", err)
			return err
		}

		status := sqlc.FormStatusPublished
		_, err = w.formService.UpdateForm(ctx, dto.UpdateFormDTO{
			Status: &status,
		}, formId.String())
		if err != nil {
			fmt.Printf("cannot update form status: %v", err)
			return err
		}

		return nil
	}
}
