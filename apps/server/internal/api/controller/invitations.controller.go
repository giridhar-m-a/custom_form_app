package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/giridhar-m-a/custom_form_app/internal/api/handler"
)

func RegisterInvitationsRoutes(rg *gin.RouterGroup) {
	invitationsHandler := handler.NewInvitationHandler()

	rg.POST("/invitations/:formId", invitationsHandler.CreateInvitation)
	rg.POST("/invitations", invitationsHandler.CreateSingleInvitation)
	rg.DELETE("/invitations/:id", invitationsHandler.DeleteInvitation)
	rg.GET("/invitations", invitationsHandler.GetInvitationByFormId)

}