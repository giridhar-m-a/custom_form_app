package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/giridhar-m-a/custom_form_app/internal/dto"
	"github.com/giridhar-m-a/custom_form_app/internal/dto/scheduler_dto"
	"github.com/giridhar-m-a/custom_form_app/internal/services"
	"github.com/giridhar-m-a/custom_form_app/internal/utils"
	"github.com/hibiken/asynq"
)

type UserWorker struct {
	userService   services.UserService
	formService   services.FormService
}

func NewUserWorker(service services.UserService, formService services.FormService) *UserWorker {
	return &UserWorker{
		userService:   service,
		formService:   formService,
	}
}

func (w *UserWorker) HandleDeleteUserAccount() asynq.HandlerFunc {
	return func(ctx context.Context, task *asynq.Task) error {
		var payload scheduler_dto.UserSchedulerPayload
		if err := json.Unmarshal(task.Payload(), &payload); err != nil {
			log.Printf("[User Worker] Error unmarshalling task payload: %v", err)
			return fmt.Errorf("cannot unmarshal task payload: %w", err)
		}

		userId, err := utils.ConvertStringToUUID(payload.UserID)
		if err != nil {
			log.Printf("[User Worker] Error converting string to uuid: %v", err)
			return nil
		}

		bucket := utils.GetEnv("MINIO_BUCKET_NAME", "custom-form-app")

		// --- Delete form files in paginated manner with goroutines ---
		page := 1
		limit := 10
		const maxWorkers = 5

		for {
			formList, err := w.formService.GetForms(ctx, userId.String(), dto.ListFormQuery{
				Query: dto.Query{
					Page:  page,
					Limit: limit,
				},
			})
			if err != nil {
				log.Printf("[User Worker] Error fetching forms page %d: %v", page, err)
				return nil
			}

			if len(formList.Forms) == 0 {
				break
			}

			// Use a semaphore channel to limit concurrent goroutines
			sem := make(chan struct{}, maxWorkers)
			var wg sync.WaitGroup
			var mu sync.Mutex
			var deleteErrors []error

			for _, form := range formList.Forms {
				formId := form.FormID.String() // adjust field name to match your struct

				wg.Add(1)
				sem <- struct{}{} // acquire slot

				go func(fId string) {
					defer wg.Done()
					defer func() { <-sem }() // release slot

					if err := services.DeleteFolderBulk(bucket, fmt.Sprintf("forms/%s", fId), ctx); err != nil {
						log.Printf("[User Worker] Error deleting form folder %s: %v", fId, err)
						mu.Lock()
						deleteErrors = append(deleteErrors, err)
						mu.Unlock()
					}
				}(formId)
			}

			wg.Wait()

			if len(deleteErrors) > 0 {
				log.Printf("[User Worker] %d form folder deletions failed on page %d", len(deleteErrors), page)
				// decide: return nil to continue, or return error to retry task
			}

			// If we've fetched all pages, stop
			if page >= formList.Pages {
				break
			}
			page++
		}

		// --- Delete the user's own folder ---
		if err := services.DeleteFolderBulk(bucket, fmt.Sprintf("users/%s", userId.String()), ctx); err != nil {
			log.Printf("[User Worker] Error deleting user folder: %v", err)
			return nil
		}

		// --- Delete the user from DB ---
		if err := w.userService.DeleteUser(ctx, userId.String()); err != nil {
			log.Printf("[User Worker] Error deleting user: %v", err)
			return nil
		}

		log.Printf("[User Worker] Successfully deleted user %s and all associated data", userId.String())
		return nil
	}
}
