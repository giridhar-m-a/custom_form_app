package services

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"sync"

	"github.com/giridhar-m-a/custom_form_app/internal/db/sqlc"
	"github.com/giridhar-m-a/custom_form_app/internal/dto"
	"github.com/giridhar-m-a/custom_form_app/internal/repositories"
	"github.com/giridhar-m-a/custom_form_app/internal/utils"
	"github.com/google/uuid"
)

type InvitationService interface {
	CreateInvitation(fileHeader *multipart.FileHeader, formID, userID uuid.UUID, ctx context.Context) (successCount int, failedCount int, err error)
	CreateSingleInvitation(invitation dto.CreateInvitationDTO, userID string, ctx context.Context) (sqlc.CreateInvitationRow, error)
	UpdateInvitationStatus(status dto.UpdateInvitationDTO, invitationID uuid.UUID, ctx context.Context) (sqlc.UpdateInvitationStatusRow, error)
	DeleteInvitation(invitationID uuid.UUID, ctx context.Context) error
	GetInvitationByFormId(query dto.InvitationListQueryDto, ctx context.Context) (dto.InvitationListDto, error)
}

type invitationService struct {
	repo repositories.InvitationRepository
	db   *sql.DB
}

func NewInvitationService(repo repositories.InvitationRepository, db *sql.DB) InvitationService {
	return &invitationService{repo: repo, db: db}
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

	return s.repo.CreateSingleInvitation(sqlc.CreateInvitationParams{
		FormID:    formId,
		Email:     invitation.Email,
		Name:      invitation.Name,
		InvitedBy: user,
	}, ctx)
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
	exclude := sqlc.NullInvitationStatus{
		InvitationStatus: query.Exclude,
		Valid:            query.Exclude != "",
	}
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

	log.Println("Count: ", count)
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

	return totalSuccess, totalFailed, nil
}
