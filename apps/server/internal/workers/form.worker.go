package workers

import (
	"context"
	"encoding/json"
	"log"

	"github.com/giridhar-m-a/custom_form_app/internal/db/sqlc"
	"github.com/giridhar-m-a/custom_form_app/internal/dto/scheduler_dto"
	"github.com/giridhar-m-a/custom_form_app/internal/repositories"
	"github.com/giridhar-m-a/custom_form_app/internal/services"
	"github.com/giridhar-m-a/custom_form_app/internal/utils"
	"github.com/hibiken/asynq"
)

type FormWorker struct {
	repo repositories.FormsRepository
}

func NewFormWorker(repo repositories.FormsRepository) *FormWorker {
	return &FormWorker{
		repo: repo,
	}
}

func (w *FormWorker) HandleFormStatusUpdate() asynq.HandlerFunc {
	return func(ctx context.Context, task *asynq.Task) error {
		log.Printf("[form worker] initiated form worker")
		var payload scheduler_dto.InvitationSchedulerPayload
		if err := json.Unmarshal(task.Payload(), &payload); err != nil {
			log.Printf("cannot unmarshal task payload: %v", err)
			return err
		}

		formId, err := utils.ConvertStringToUUID(payload.FormID)
		log.Printf("[form worker] form id: %s", formId.String())
		if err != nil {
			log.Printf("cannot convert string to uuid: %v", err)
			return err
		}

		status := sqlc.FormStatusPublished
		_, err = w.repo.UpdateForm(sqlc.UpdateFormParams{
			FormStatus: sqlc.NullFormStatus{
				FormStatus: status,
				Valid:      true,
			},
			FormID: formId,
		}, ctx)
		if err != nil {
			log.Printf("cannot update form status for %s: %v", formId, err)
			return err
		}
		log.Printf("[form worker] form status updated for %s", formId.String())
		return nil
	}
}

func (w *FormWorker) HandleFormDelete() asynq.HandlerFunc {
	return func(ctx context.Context, task *asynq.Task) error {
		log.Printf("[form worker] initiated form worker")
		var payload scheduler_dto.InvitationSchedulerPayload
		if err := json.Unmarshal(task.Payload(), &payload); err != nil {
			log.Printf("cannot unmarshal task payload: %v", err)
			return err
		}
		formId:= payload.FormID
		_, err := w.repo.DeleteForm(formId, ctx)
		if err != nil {
			log.Printf("cannot delete form for %s: %v", formId, err)
			return err
		}
		bucketName:= utils.GetEnv("MINIO_BUCKET", "custom-forms-bucket")
		err = services.DeleteFolderBulk(bucketName, "forms/"+formId, ctx)
		if err != nil {
			log.Printf("cannot delete form folder for %s: %v", formId, err)
			return err
		}
		log.Printf("[form worker] form deleted for %s", formId)
		return nil
	}
}