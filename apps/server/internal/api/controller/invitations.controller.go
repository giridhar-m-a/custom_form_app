package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/giridhar-m-a/custom_form_app/internal/api/handler"
	"github.com/giridhar-m-a/custom_form_app/internal/api/middleware"
)

func RegisterInvitationsRoutes(rg *gin.RouterGroup) {
	invitationsHandler := handler.NewInvitationHandler()
	api := rg.Group("/invitations")
	api.POST("/verify", invitationsHandler.VerifyInvitation)
	protected := api.Use(middleware.AuthMiddleware())
	protected.POST("/:formId", invitationsHandler.CreateInvitation)
	protected.POST("", invitationsHandler.CreateSingleInvitation)
	protected.DELETE("/:id", invitationsHandler.DeleteInvitation)
	protected.GET("", invitationsHandler.GetInvitationByFormId)
	protected.POST("/anonymous", invitationsHandler.GenerateAnonymousInvitationToken)
}
