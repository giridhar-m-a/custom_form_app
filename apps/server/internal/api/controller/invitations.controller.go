package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/giridhar-m-a/custom_form_app/internal/api/handler"
)

func RegisterInvitationsRoutes(rg *gin.RouterGroup) {
	invitationsHandler := handler.NewInvitationHandler()
	api := rg.Group("/invitations")

	api.POST("/:formId", invitationsHandler.CreateInvitation)
	api.POST("", invitationsHandler.CreateSingleInvitation)
	api.DELETE("/:id", invitationsHandler.DeleteInvitation)
	api.GET("", invitationsHandler.GetInvitationByFormId)

}
