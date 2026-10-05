package webhook

import (
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/giridhar-m-a/custom_form_app/internal/db"
	"github.com/giridhar-m-a/custom_form_app/internal/db/sqlc"
	"github.com/giridhar-m-a/custom_form_app/internal/dto"
	"github.com/giridhar-m-a/custom_form_app/internal/repositories"
	"github.com/giridhar-m-a/custom_form_app/internal/utils"
	"github.com/resend/resend-go/v3"
)

func HandleInvitationWebhook(c *gin.Context) {
	id := c.Request.Header.Get("svix-id")
	timestamp := c.Request.Header.Get("svix-timestamp")
	signature := c.Request.Header.Get("svix-signature")
	if id == "" || timestamp == "" || signature == "" {
		log.Printf("[Webhook: Invitation] Missing required headers")
		c.JSON(400, gin.H{"error": "Missing required headers"})
		return
	}
	body, err := c.GetRawData()
	if err != nil {
		log.Printf("[Webhook: Invitation] Error getting raw data: %v", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[Webhook: Invitation] Webhook received: %s\n", string(body))

	secrete := utils.GetEnv("RESEND_WEBHOOK_SIGNATURE", "")

	er := utils.ResendClient.Webhooks.Verify(&resend.VerifyWebhookOptions{
		Payload: string(body),
		Headers: resend.WebhookHeaders{
			Id:        id,
			Timestamp: timestamp,
			Signature: signature,
		},
		WebhookSecret: secrete,
	})
	if er != nil {
		log.Printf("[Webhook: Invitation] Error verifying webhook: %v", er)
		c.JSON(400, gin.H{"error": er.Error()})
		return
	}

	var req dto.ResendWebhookRequest
	if err := json.Unmarshal(body, &req); err != nil {
		log.Printf("[Webhook: Invitation] Error binding JSON: %v", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	log.Printf("[Webhook: Invitation] Webhook received: %+v\n", req)
	invitationRepo := repositories.NewInvitationRepository(db.Queries)

	invitationId, err := utils.ConvertStringToUUID(req.Data.Tags.InvitationId)
	if err != nil {
		log.Printf("[Webhook: Invitation] Error converting invitation ID: %v", err)
		return
	}

	resp, err := invitationRepo.UpdateInvitationStatus(sqlc.UpdateInvitationStatusParams{
		InvitationID: invitationId,
		Status:       getInvitationStatus(req.Type),
	}, c)
	if err != nil {
		log.Printf("[Webhook: Invitation] Error updating invitation status: %v", err)
		return
	}
	log.Printf("[Webhook: Invitation] Invitation status updated successfully: %+v\n", resp)
	c.JSON(200, gin.H{
		"status":  200,
		"message": "Updated the invitation",
		"data":    resp,
	})
}

func getInvitationStatus(event string) sqlc.InvitationStatus {
	switch event {
	case "email.sent":
		return sqlc.InvitationStatusSubmitted
	case "email.delivered":
		return sqlc.InvitationStatusDelivered
	case "email.delivery_delayed":
		return sqlc.InvitationStatusDelayed
	case "email.bounced":
		return sqlc.InvitationStatusBounced
	case "email.complained":
		return sqlc.InvitationStatusComplained
	case "email.opened":
		return sqlc.InvitationStatusOpened
	case "email.clicked":
		return sqlc.InvitationStatusClicked
	default:
		return sqlc.InvitationStatusPending
	}
}
