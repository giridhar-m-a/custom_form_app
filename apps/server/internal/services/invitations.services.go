package services

import (
	"context"
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"sync"
	"time"

	"github.com/giridhar-m-a/custom_form_app/constants"
	"github.com/giridhar-m-a/custom_form_app/internal/db/sqlc"
	"github.com/giridhar-m-a/custom_form_app/internal/dto"
	"github.com/giridhar-m-a/custom_form_app/internal/repositories"
	"github.com/giridhar-m-a/custom_form_app/internal/scheduler"
	"github.com/giridhar-m-a/custom_form_app/internal/services/templates"
	"github.com/giridhar-m-a/custom_form_app/internal/utils"
	"github.com/google/uuid"
	"github.com/resend/resend-go/v3"
)

type InvitationService interface {
	CreateInvitation(fileHeader *multipart.FileHeader, formID, userID uuid.UUID, ctx context.Context) (successCount int, failedCount int, err error)
	CreateSingleInvitation(invitation dto.CreateInvitationDTO, userID string, ctx context.Context) (sqlc.CreateInvitationRow, error)
	UpdateInvitationStatus(status dto.UpdateInvitationDTO, resendID uuid.UUID, ctx context.Context) (sqlc.UpdateInvitationStatusRow, error)
	DeleteInvitation(invitationID uuid.UUID, ctx context.Context) error
	GetInvitationByFormId(query dto.InvitationListQueryDto, ctx context.Context) (dto.InvitationListDto, error)
}

type invitationService struct {
	repo repositories.InvitationRepository
	form FormService
	db   *sql.DB
}

func NewInvitationService(repo repositories.InvitationRepository, form FormService, db *sql.DB) InvitationService {
	return &invitationService{repo: repo, form: form, db: db}
}

func (s *invitationService) CreateSingleInvitation(invitation dto.CreateInvitationDTO, userID string, ctx context.Context) (sqlc.CreateInvitationRow, error) {
	formId, err := utils.ConvertStringToUUID(invitation.FormID)
	if err != nil {
		return sqlc.CreateInvitationRow{}, err
	}

	user, err := utils.ConvertStringToUUID(userID)
	if err != nil {
		return sqlc.CreateInvitationRow{}, err
	}

	form, err := s.form.GetSingleForm(ctx, invitation.FormID)
	if form.FormStatus.FormStatus == sqlc.FormStatusClosed {
		return sqlc.CreateInvitationRow{}, errors.New("form is closed")
	}
	if err != nil {
		log.Printf("[Invitation Service] Error getting form: %v", err)
		return sqlc.CreateInvitationRow{}, err
	}
	if form.FormStatus.FormStatus == sqlc.FormStatusClosed || (form.ClosingTime.Valid && form.ClosingTime.Time.Before(time.Now())) {
		return sqlc.CreateInvitationRow{}, errors.New("form is closed")
	}

	createdInvitation, err := s.repo.CreateSingleInvitation(sqlc.CreateInvitationParams{
		FormID:    formId,
		Email:     invitation.Email,
		Name:      invitation.Name,
		InvitedBy: user,
	}, ctx)
	if err != nil {
		log.Printf("[Invitation Service] Error creating invitation: %v", err)
		return sqlc.CreateInvitationRow{}, err
	}
	now := time.Now()

	isScheduled := form.IsScheduled.Valid && form.IsScheduled.Bool

	shouldInviteNow := false

	if isScheduled && form.ScheduledTime.Valid && form.ScheduledTime.Time.Before(now) {
		shouldInviteNow = true
	} else if isScheduled && form.ScheduledTime.Valid && form.InvitationScheduleGap.Valid &&
		time.Since(form.ScheduledTime.Time) < time.Duration(form.InvitationScheduleGap.Int32)*time.Minute {
		shouldInviteNow = true
	}

	if (form.IsScheduled.Valid && !form.IsScheduled.Bool) || shouldInviteNow {
		mail := NewMailService(utils.ResendClient)
		templateService := templates.NewService()
		jwtService := NewJWTService()
		token, err := jwtService.GenerateInvitationToken(createdInvitation.InvitationID.String(), formId.String(), time.Until(form.ClosingTime.Time))
		if err != nil {
			log.Printf("[Invitation Service] Error generating invitation token: %v", err)
			return sqlc.CreateInvitationRow{}, err
		}
		frontendURL := utils.GetEnv("FRONTEND_URL", "")
		fromEmail := utils.GetEnv("SENDER_EMAIL_ADDRESS", "form-genius-no-reply@giridhar.dev")
		sender := fmt.Sprintf("Form Genius <%s>", fromEmail)
		template, err := templateService.Render(constants.InvitationTemplate, dto.InvitationEmailParams{
			PlatformName:   "Form Genius",
			UserName:       invitation.Name,
			Title:          form.FormTitle,
			InvitationURL:  fmt.Sprintf("%s/user-response?token=%s", frontendURL, token),
			Year:           time.Now().Year(),
			CompanyAddress: "Form Genius Inc",
		})
		if err != nil {
			log.Printf("[Invitation Service] Error generating invitation template: %v", err)
			return sqlc.CreateInvitationRow{}, err
		}
		resendRes, err := mail.SendEmail(resend.SendEmailRequest{
			From:    sender,
			To:      []string{invitation.Email},
			Subject: fmt.Sprintf("Invitation to fill form: %s", form.FormTitle),
			Html:    template,
			Tags: []resend.Tag{
				{
					Name:  "invitation",
					Value: "true",
				},
				{
					Name:  "invitation_id",
					Value: createdInvitation.InvitationID.String(),
				},
			},
		})
		if err != nil {
			log.Printf("[Invitation Service] Error sending email: %v", err)
			return sqlc.CreateInvitationRow{}, err
		}
		log.Printf("[Invitation Service] email sent successfully: %v, %v", resendRes.Id, createdInvitation.InvitationID.String())

	}
	return createdInvitation, nil
}

func (s *invitationService) UpdateInvitationStatus(status dto.UpdateInvitationDTO, invitationID uuid.UUID, ctx context.Context) (sqlc.UpdateInvitationStatusRow, error) {

	return s.repo.UpdateInvitationStatus(sqlc.UpdateInvitationStatusParams{
		InvitationID: invitationID,
		Status:       status.Status,
	}, ctx)
}

func (s *invitationService) DeleteInvitation(invitationID uuid.UUID, ctx context.Context) error {
	return s.repo.DeleteInvitation(invitationID, ctx)
}

func (s *invitationService) GetInvitationByFormId(query dto.InvitationListQueryDto, ctx context.Context) (dto.InvitationListDto, error) {

	formID, err := utils.ConvertStringToUUID(query.FormId)
	fmt.Printf("Form ID: %v\n", query.FormId)
	if err != nil {
		return dto.InvitationListDto{}, err
	}
	search := utils.ConvertStringToNullString(query.Search)
	exclude := query.Exclude
	status := sqlc.NullInvitationStatus{
		InvitationStatus: query.Status,
		Valid:            query.Status != "",
	}
	page := query.Page
	if page <= 0 {
		page = 1
	}
	limitVal := query.Limit
	if limitVal <= 0 {
		limitVal = 10
	}
	offsetVal := (page - 1) * limitVal

	offset := utils.ConvertIntToNullInt32(offsetVal)
	limit := utils.ConvertIntToNullInt32(limitVal)

	invitations, err := s.repo.GetInvitationByFormId(sqlc.GetInvitationByFormIdParams{
		FormID:        formID,
		Search:        search,
		Status:        status,
		OffsetVal:     offset,
		LimitVal:      limit,
		ExcludeStatus: exclude,
	}, ctx)
	if err != nil {
		return dto.InvitationListDto{}, err
	}
	count, err := s.repo.CountInvitationsByFormId(sqlc.CountInvitationsByFormIdParams{
		FormID:        formID,
		Search:        search,
		Status:        status,
		ExcludeStatus: exclude,
	}, ctx)
	if err != nil {
		return dto.InvitationListDto{}, err
	}

	var limitInt int
	if limit.Valid {
		limitInt = int(limit.Int32)
	} else {
		limitInt = 10
	}

	totalCount := int(count)
	totalPages := totalCount / limitInt
	if totalCount%limitInt != 0 {
		totalPages++
	}
	if totalPages == 0 {
		totalPages = 1
	}

	return dto.InvitationListDto{
		Invitations: invitations,
		Total:       totalCount,
		Page:        query.Page,
		Limit:       query.Limit,
		Pages:       totalPages,
	}, nil
}

func (s *invitationService) DeleteInvitationByFormId(formID string, ctx context.Context) error {
	id, err := utils.ConvertStringToUUID(formID)
	if err != nil {
		return err
	}
	return s.repo.DeleteInvitation(id, ctx)
}

func (s *invitationService) CreateInvitation(
	fileHeader *multipart.FileHeader,
	formID, userID uuid.UUID,
	ctx context.Context,
) (successCount int, failedCount int, err error) {

	form, err := s.form.GetSingleForm(ctx, formID.String())
	if err != nil {
		return 0, 0, err
	}

	if form.FormStatus.FormStatus == sqlc.FormStatusClosed {
		return 0, 0, errors.New("form is closed")
	}

	if form.ClosingTime.Valid && form.ClosingTime.Time.Before(time.Now()) {
		return 0, 0, errors.New("form is closed")
	}

	// 1. Open the CSV file
	file, err := fileHeader.Open()
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	// 2. Setup Communication Channels
	type row struct{ Email, Name string }
	rowChan := make(chan row, 500)
	errChan := make(chan error, 1) // Captures fatal DB errors
	var wg sync.WaitGroup

	// Counters protected by Mutex
	var mu sync.Mutex
	var totalSuccess int
	var totalFailed int

	// 3. START THE CONSUMER (Worker)
	wg.Add(1)
	go func() {
		defer wg.Done()
		batchSize := 5000
		emails := make([]string, 0, batchSize)
		names := make([]string, 0, batchSize)

		flush := func() error {
			if len(emails) == 0 {
				return nil
			}

			// Batch Insert
			res, err := s.repo.CreateInvitation(sqlc.CreateManyInvitationsParams{
				FormID:    formID,
				InvitedBy: userID,
				Emails:    emails,
				Names:     names,
			}, ctx)

			mu.Lock()
			defer mu.Unlock()

			if err != nil {
				// If a batch fails, we treat it as a fatal error to stop the process
				return err
			}

			totalSuccess += len(res)
			// Calculate failures (usually duplicates if using ON CONFLICT DO NOTHING)
			totalFailed += (len(emails) - len(res))
			return nil
		}

		for r := range rowChan {
			emails = append(emails, r.Email)
			names = append(names, r.Name)
			if len(emails) >= batchSize {
				if err := flush(); err != nil {
					errChan <- err // Signal fatal error to producer
					return
				}
				emails = emails[:0]
				names = names[:0]
			}
		}
		// Final batch
		if err := flush(); err != nil {
			errChan <- err
		}
	}()

	// 4. START THE PRODUCER (CSV Reader)
	reader := csv.NewReader(file)
	_, _ = reader.Read() // Skip header line

ReadLoop: // This label allows us to break out of the for-loop from inside the select
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break ReadLoop
		}
		if err != nil {
			// Log and skip malformed lines so one bad row doesn't kill 100k upload
			mu.Lock()
			totalFailed++
			mu.Unlock()
			continue ReadLoop
		}

		// Check for two data columns
		if len(record) < 2 {
			mu.Lock()
			totalFailed++
			mu.Unlock()
			continue ReadLoop
		}

		select {
		case fatalErr := <-errChan:
			// The database worker failed; stop the reader
			return 0, 0, fatalErr
		case <-ctx.Done():
			// User cancelled the request (closed browser/tab)
			return 0, 0, ctx.Err()
		case rowChan <- row{Email: record[0], Name: record[1]}:
			// Successfully pushed to the conveyor belt
		}
	}

	// 5. CLEANUP & COORDINATION
	close(rowChan) // Signal worker that no more rows are coming
	wg.Wait()      // Wait for the final batch to finish

	// Check one last time if the final flush failed
	select {
	case fatalErr := <-errChan:
		return 0, 0, fatalErr
	default:
	}

	isScheduled := form.IsScheduled.Valid && form.IsScheduled.Bool
	now := time.Now()

	shouldInviteNow := false

	if form.FormStatus.FormStatus == sqlc.FormStatusPublished {
		shouldInviteNow = true
	} else if isScheduled && form.ScheduledTime.Valid && form.ScheduledTime.Time.Before(now) {
		shouldInviteNow = true
	} else if isScheduled && form.ScheduledTime.Valid && form.InvitationScheduleGap.Valid &&
		time.Since(form.ScheduledTime.Time) < time.Duration(form.InvitationScheduleGap.Int32) {
		shouldInviteNow = true
	}

	if shouldInviteNow || !isScheduled {
		log.Printf("[Invitation Service] Scheduling invitation in first block")
		scheduler.ScheduleInvitation(formID.String(), now.Add(5*time.Minute))

	} else if isScheduled && form.ScheduledTime.Valid && form.ScheduledTime.Time.After(now) {
		log.Printf("[Invitation Service] Scheduling invitation in second block")
		runAt := form.ScheduledTime.Time.Add(
			-time.Duration(form.InvitationScheduleGap.Int32) * time.Minute,
		)
		scheduler.ScheduleInvitation(formID.String(), runAt)
	}

	return totalSuccess, totalFailed, nil
}
