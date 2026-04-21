package services

import (
	"context"
	"fmt"
	"log"

	"github.com/resend/resend-go/v3"
)

type MailService interface {
	SendEmail(resend.SendEmailRequest) (*resend.SendEmailResponse, error)
	SendBulk(ctx context.Context, params []*resend.SendEmailRequest) ([]resend.SendEmailResponse, error)
}

type mailService struct {
	client *resend.Client
}

func NewMailService(client *resend.Client) MailService {
	return &mailService{client: client}
}

func (ms *mailService) SendEmail(params resend.SendEmailRequest) (*resend.SendEmailResponse, error) {

	sent, err := ms.client.Emails.Send(&params)
	if err != nil {
		log.Printf("[Mail Service] Error sending email: %v\n", err)
		return nil, err
	}
	fmt.Printf("[Mail Service] email sent id: %+v\n", sent.Id)

	return sent, nil
}

func (ms *mailService) SendBulk(ctx context.Context, params []*resend.SendEmailRequest) ([]resend.SendEmailResponse, error) {

	sent, err := ms.client.Batch.SendWithContext(ctx, params)
	if err != nil {
		log.Printf("[Mail Service] Error sending bulk emails: %v\n", err)
		return nil, err
	}
	fmt.Printf("[Mail Service] bulk emails sent ids: %+v\n", sent.Data)

	return sent.Data, nil
}
