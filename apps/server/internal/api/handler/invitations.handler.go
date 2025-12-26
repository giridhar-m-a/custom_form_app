package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giridhar-m-a/custom_form_app/internal/db"
	"github.com/giridhar-m-a/custom_form_app/internal/dto"
	"github.com/giridhar-m-a/custom_form_app/internal/repositories"
	"github.com/giridhar-m-a/custom_form_app/internal/services"
	"github.com/giridhar-m-a/custom_form_app/internal/utils"
)

type InvitationHandler interface {
	CreateInvitation(c *gin.Context)
	CreateSingleInvitation(c *gin.Context)
	DeleteInvitation(c *gin.Context)
	GetInvitationByFormId(c *gin.Context)
}

type invitationHandler struct {
	svc services.InvitationService
}

func NewInvitationHandler() InvitationHandler {
	queries := db.Queries
	repo := repositories.NewInvitationRepository(queries)
	conn := db.Connection
	return &invitationHandler{svc: services.NewInvitationService(repo, conn)}
}

// NewInvitationHandler creates a new invitation handler
// @Summary NewInvitationHandler creates a new invitation handler
// @Description NewInvitationHandler creates a new invitation handler
// @Tags Invitation
// @Accept multipart/form-data
// @Produce json
// @Param formId path string true "Form ID"
// @Param file formData file true "File to upload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /invitations/{formId} [post]
// @Security BearerAuth
// @type http
// @scheme bearer
func (h *invitationHandler) CreateInvitation(c *gin.Context) {
	// 1. Get Path/Param IDs
	formIDStr := c.Param("formId")
	formID, err := utils.ConvertStringToUUID(formIDStr)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	// 2. Get UserID (Assumes you have middleware setting "user_id" in context)
	userIDStr := c.GetString("userID")
	userID, err := utils.ConvertStringToUUID(userIDStr)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	// 3. Retrieve the File from the Request
	fileHeader, err := c.FormFile("file")
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	// 4. Call the Service
	// We pass c.Request.Context() so if the user cancels the request, the service knows
	success, failed, err := h.svc.CreateInvitation(fileHeader, formID, userID, c.Request.Context())
	
	if err != nil {
		// If it's a fatal error (like DB connection loss)
		utils.HandleError(c, err)
		return
	}

	// 5. Return the Summary
	c.JSON(http.StatusOK, gin.H{
		"message": "Processing complete",
		"data": gin.H{
			"success_count": success,
			"failed_count":  failed,
			"total_rows":    success + failed,
		},
	})
}

// CreateSingleInvitation creates a single invitation
// @Summary CreateSingleInvitation creates a single invitation
// @Description CreateSingleInvitation creates a single invitation
// @Tags Invitation
// @Accept json
// @Produce json
// @Param data body dto.CreateInvitationDTO true "Invitation data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /invitations [post]
// @Security BearerAuth
// @type http
// @scheme bearer
func (h *invitationHandler) CreateSingleInvitation(c *gin.Context) {

	var data dto.CreateInvitationDTO
	if err := c.ShouldBind(&data); err != nil {
		utils.HandleError(c, err)
		return
	}

	user := c.GetString("userID")
	
	invitedUser, err := h.svc.CreateSingleInvitation(data, user, c)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Invitation created successfully",
		"data":    invitedUser,
	})
}

// DeleteInvitation deletes an invitation
// @Summary DeleteInvitation deletes an invitation
// @Description DeleteInvitation deletes an invitation
// @Tags Invitation
// @Accept json
// @Produce json
// @Param id path string true "Invitation ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /invitations/{id} [delete]
// @Security BearerAuth
// @type http
// @scheme bearer
func (h *invitationHandler) DeleteInvitation(c *gin.Context) {

	formIDStr := c.Param("id")
	formID, err := utils.ConvertStringToUUID(formIDStr)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	err = h.svc.DeleteInvitation(formID, c)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Invitation deleted successfully",
		"status":  http.StatusOK,
	})

}

// GetInvitationByFormId retrieves invitations by form ID
// @Summary GetInvitationByFormId retrieves invitations by form ID
// @Description GetInvitationByFormId retrieves invitations by form ID
// @Tags Invitation
// @Accept json
// @Produce json
// @Param   query query dto.InvitationListQueryDto true "Invitation List Query"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /invitations [get]
// @Security BearerAuth
// @type http
// @scheme bearer
func (h *invitationHandler) GetInvitationByFormId(c *gin.Context) {

	
	var params dto.InvitationListQueryDto
	if err := c.ShouldBindQuery(&params); err != nil {
		utils.HandleError(c, err)
		return
	}

	invitations, err := h.svc.GetInvitationByFormId(params, c)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	fmt.Printf("Invitations: %v\n", invitations)

	c.JSON(http.StatusOK, gin.H{
		"message": "Invitations retrieved successfully",
		"data":    invitations,
	})
}
