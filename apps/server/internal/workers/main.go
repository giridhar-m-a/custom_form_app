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
			Username: cache.Client.Options().Username,
		},
		asynq.Config{
			Concurrency: concurrency,
			Queues: map[string]int{
				constants.QueueInvitations: 10,
				constants.QueueFormStatus:  10,
				constants.QueueFormDelete:  10,
				constants.QueueDeleteUser:  10,
				"default":                  5,
			},
		},
	)

	mux := asynq.NewServeMux()
	formRepo := repositories.NewFormsRepository(db.Queries)
	formService := services.NewFormService(
		formRepo,
		repositories.NewFormFieldsRepository(db.Queries),
		repositories.NewFormFieldOptionsRepository(db.Queries),
		db.Connection,
	)
	userRepo := repositories.NewSQLCUserRepository(db.Queries)
	userService := services.UserServiceProvider(userRepo)

	formWorker := NewFormWorker(formRepo)

	invitationWorker := NewInvitationWorker(
		services.NewInvitationService(
			repositories.NewInvitationRepository(db.Queries),
			formService,
			db.Connection,
		),
		formService,
	)

	userWorker := NewUserWorker(
		userService,
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

	mux.HandleFunc(
		constants.TaskTypeFormDelete,
		formWorker.HandleFormDelete(),
	)

	mux.HandleFunc(
		constants.TaskTypeDeleteUser,
		userWorker.HandleDeleteUserAccount(),
	)

	log.Println("Asynq worker starting...")

	go func() {
		if err := srv.Run(mux); err != nil {
			log.Fatal(err)
		}
	}()
}
