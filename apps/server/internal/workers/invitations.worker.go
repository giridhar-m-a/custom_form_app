package workers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/giridhar-m-a/custom_form_app/constants"
	"github.com/giridhar-m-a/custom_form_app/internal/db/sqlc"
	"github.com/giridhar-m-a/custom_form_app/internal/dto"
	"github.com/giridhar-m-a/custom_form_app/internal/dto/scheduler_dto"
	"github.com/giridhar-m-a/custom_form_app/internal/services"
	"github.com/giridhar-m-a/custom_form_app/internal/services/templates"
	"github.com/giridhar-m-a/custom_form_app/internal/utils"
	"github.com/hibiken/asynq"
	"github.com/resend/resend-go/v3"
)

type InvitationWorker struct {
	invitationService services.InvitationService
	formService       services.FormService
}

func NewInvitationWorker(service services.InvitationService, formService services.FormService) *InvitationWorker {
	return &InvitationWorker{
		invitationService: service,
		formService:       formService,
	}
}

func (w *InvitationWorker) HandleInvitationsSchedule() asynq.HandlerFunc {
	return func(ctx context.Context, task *asynq.Task) error {
		var payload scheduler_dto.InvitationSchedulerPayload
		if err := json.Unmarshal(task.Payload(), &payload); err != nil {
			log.Printf("[Invitation Worker] Error unmarshalling task payload: %v", err)
			return fmt.Errorf("cannot unmarshal task payload: %w", err)
		}

		formId, err := utils.ConvertStringToUUID(payload.FormID)
		if err != nil {
			log.Printf("[Invitation Worker] Error converting string to uuid: %v", err)
			return nil
		}

		form, err := w.formService.GetSingleForm(ctx, formId.String())
		if err != nil {
			log.Printf("[Invitation Worker] Error getting form %s: %v",formId, err)
			return nil
		}
		if form.FormStatus.FormStatus == sqlc.FormStatusClosed || (form.ClosingTime.Valid && form.ClosingTime.Time.Before(time.Now())) {
			log.Printf("[Invitation Worker] Form is closed or closing time is reached: %v", formId)
			return nil
		}

		const pageSize = 50
		baseParams := dto.InvitationListQueryDto{
			FormId:  formId.String(),
			Exclude: []sqlc.InvitationStatus{sqlc.InvitationStatusBounced, sqlc.InvitationStatusClicked, sqlc.InvitationStatusOpened, sqlc.InvitationStatusComplained, sqlc.InvitationStatusDelayed, sqlc.InvitationStatusFailed, sqlc.InvitationStatusDelivered, sqlc.InvitationStatusSubmitted},
			Query:   dto.Query{Limit: pageSize},
		}

		firstPage, err := w.invitationService.GetInvitationByFormId(dto.InvitationListQueryDto{
			FormId:  baseParams.FormId,
			Exclude: baseParams.Exclude,
			Query:   dto.Query{Page: 1, Limit: pageSize},
		}, ctx)
		if err != nil {
			log.Printf("[Invitation Worker] Error getting invitation list: %v", err)
			return fmt.Errorf("cannot get invitation list: %w", err)
		}
		if len(firstPage.Invitations) == 0 {
			log.Printf("[Invitation Worker] No invitations found for form: %v", formId)
			return nil
		}

		var (
			mailService     = services.NewMailService(utils.ResendClient)
			jwtService      = services.NewJWTService()
			templateService = templates.NewService()
			closingTime     = form.ClosingTime.Time
			frontendURL     = utils.GetEnv("FRONTEND_URL", "")
			fromEmail       = utils.GetEnv("SENDER_EMAIL_ADDRESS", "form-genius-no-reply@giridhar.dev")
			sender          = fmt.Sprintf("Form Genius <%s>", fromEmail)
		)

		var (
			errMu   sync.Mutex
			errs    []error
			emailWg sync.WaitGroup
		)

		sendEmails := func(invitations []sqlc.Invitation) {
			defer emailWg.Done()

			emailPayloads := make([]*resend.SendEmailRequest, 0, len(invitations))

			for _, inv := range invitations {
				token, err := jwtService.GenerateInvitationToken(inv.InvitationID.String(), formId.String(), time.Until(closingTime))
				if err != nil {
					log.Printf("[Invitation Worker] Error generating token for %s: %v", inv.InvitedEmail, err)
					errMu.Lock()
					errs = append(errs, fmt.Errorf("generate token for %s: %w", inv.InvitedEmail, err))
					errMu.Unlock()
					continue
				}

				tmpl, err := templateService.Render(constants.InvitationTemplate, dto.InvitationEmailParams{
					PlatformName:   "Form Genius",
					UserName:       inv.InvitedName,
					Title:          form.FormTitle,
					InvitationURL:  fmt.Sprintf("%s/user-response?token=%s", frontendURL, token),
					Year:           time.Now().Year(),
					CompanyAddress: "Form Genius Inc",
				})
				if err != nil {
					log.Printf("[Invitation Worker] Error rendering template for %s: %v", inv.InvitedEmail, err)
					errMu.Lock()
					errs = append(errs, fmt.Errorf("render template for %s: %w", inv.InvitedEmail, err))
					errMu.Unlock()
					continue
				}

				emailPayloads = append(emailPayloads, &resend.SendEmailRequest{
					From:    sender,
					To:      []string{inv.InvitedEmail},
					Subject: fmt.Sprintf("Invitation to fill form: %s", form.FormTitle),
					Html:    tmpl,
					Tags: []resend.Tag{
						{
							Name:  "invitation",
							Value: "true",
						},
						{
							Name:  "invitation_id",
							Value: inv.InvitationID.String(),
						},
					},
				})
			}

			if len(emailPayloads) == 0 {
				return
			}

			_, err := mailService.SendBulk(ctx, emailPayloads)
			if err != nil {
				log.Printf("[Invitation Worker] Error sending bulk emails: %v", err)
				errMu.Lock()
				errs = append(errs, fmt.Errorf("send bulk emails: %w", err))
				errMu.Unlock()
				return
			}

		}

		// Send first page immediately
		emailWg.Add(1)
		go sendEmails(firstPage.Invitations)

		// Fetch and send remaining pages concurrently
		type pageTask struct {
			invitations []sqlc.Invitation
			err         error
		}

		totalPages := firstPage.Pages
		pageCh := make(chan pageTask, totalPages-1)
		dbSem := make(chan struct{}, 5)
		var fetchWg sync.WaitGroup

		for page := 2; page <= totalPages; page++ {
			fetchWg.Add(1)
			go func(p int) {
				defer fetchWg.Done()
				dbSem <- struct{}{}
				defer func() { <-dbSem }()

				result, err := w.invitationService.GetInvitationByFormId(dto.InvitationListQueryDto{
					FormId:  baseParams.FormId,
					Exclude: baseParams.Exclude,
					Query:   dto.Query{Page: p, Limit: pageSize},
				}, ctx)
				if err != nil {
					pageCh <- pageTask{err: fmt.Errorf("page %d: %w", p, err)}
					return
				}
				pageCh <- pageTask{invitations: result.Invitations}
			}(page)
		}

		go func() {
			fetchWg.Wait()
			close(pageCh)
		}()

		for pt := range pageCh {
			if pt.err != nil {
				log.Printf("[Invitation Worker] Error getting invitation list: %v", pt.err)
				errMu.Lock()
				errs = append(errs, pt.err)
				errMu.Unlock()
				continue
			}
			emailWg.Add(1)
			go sendEmails(pt.invitations)
		}

		emailWg.Wait()

		if len(errs) > 0 {
			log.Printf("[Invitation Worker] Completed with %d error(s): %v", len(errs), errors.Join(errs...))
			return fmt.Errorf("completed with %d error(s): %w", len(errs), errors.Join(errs...))
		}

		log.Printf("[Invitation Worker] All invitations sent successfully for form: %v", formId)
		return nil
	}
}
