package workers

import (
	"log"

	"github.com/giridhar-m-a/custom_form_app/constants"
	"github.com/giridhar-m-a/custom_form_app/internal/cache"
	"github.com/giridhar-m-a/custom_form_app/internal/db"
	"github.com/giridhar-m-a/custom_form_app/internal/repositories"
	"github.com/giridhar-m-a/custom_form_app/internal/services"
	"github.com/hibiken/asynq"
)

func Start(concurrency int) {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     cache.Client.Options().Addr,
			Password: cache.Client.Options().Password,
			DB:       cache.Client.Options().DB,
		},
		asynq.Config{
			Concurrency: concurrency,
		},
	)

	mux := asynq.NewServeMux()

	formService := services.NewFormService(
		repositories.NewFormsRepository(db.Queries),
		repositories.NewFormFieldsRepository(db.Queries),
		repositories.NewFormFieldOptionsRepository(db.Queries),
		db.Connection,
	)

	formWorker := NewFormWorker(formService)

	invitationWorker := NewInvitationWorker(
		services.NewInvitationService(repositories.NewInvitationRepository(db.Queries), formService, db.Connection),
		formService,
	)

	mux.HandleFunc(
		constants.TaskTypeFormStatusUpdate,
		formWorker.HandleFormStatusUpdate(),
	)
	mux.HandleFunc(
		constants.TaskTypeInvitationSchedule,
		invitationWorker.HandleInvitationsSchedule(),
	)

	go func() {
		if err := srv.Run(mux); err != nil {
			log.Fatal(err)
		}
		log.Println("Scheduler started successfully")
	}()
}
