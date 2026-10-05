package dto

type ResendWebhookRequest struct {
	Type      string                `json:"type"`
	Data      InvitationWebhookData `json:"data"`
	CreatedAt string                `json:"created_at"`
}

type InvitationWebhookTag struct {
	Invitation   string `json:"invitation"`
	InvitationId string `json:"invitation_id"`
}

type InvitationWebhookData struct {
	CreatedAt string               `json:"created_at"`
	EmailId   string               `json:"email_id"`
	From      string               `json:"from"`
	Subject   string               `json:"subject"`
	To        *[]string            `json:"to"`
	Tags      InvitationWebhookTag `json:"tags"`
}
