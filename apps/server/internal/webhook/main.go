package webhook

import (
	"github.com/gin-gonic/gin"
)

func RegisterWebhookRoutes(r *gin.Engine) {
	hook := r.Group("/webhook")
	{
		hook.POST("/invitation", HandleInvitationWebhook)
	}
}
