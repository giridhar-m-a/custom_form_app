package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/giridhar-m-a/custom_form_app/internal/db"
	"github.com/giridhar-m-a/custom_form_app/internal/db/sqlc"
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
	VerifyInvitation(ctx *gin.Context)
	GenerateAnonymousInvitationToken(c *gin.Context)
}

type invitationHandler struct {
	svc         services.InvitationService
	jwtService  services.JWTService
	formService services.FormService
}

func NewInvitationHandler() InvitationHandler {
	queries := db.Queries
	repo := repositories.NewInvitationRepository(queries)
	conn := db.Connection
	formService := services.NewFormService(
		repositories.NewFormsRepository(queries),
		repositories.NewFormFieldsRepository(queries),
		repositories.NewFormFieldOptionsRepository(queries),
		conn,
	)
	return &invitationHandler{svc: services.NewInvitationService(repo, formService, conn), jwtService: services.NewJWTService(), formService: formService}
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
// @Success 200 {object} dto.InvitationListResponse "Invitations retrieved successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid request payload"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
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

	var page int
	if params.Page == 0 {
		page = 1
	} else {
		page = params.Page
	}

	var limit int
	if params.Limit == 0 {
		limit = 10
	} else {
		limit = params.Limit
	}

	pagination := dto.PaginationResponse{
		TotalRecords: invitations.Total,
		Page:         page,
		Limit:        limit,
		TotalPages:   invitations.Pages,
	}

	response := []dto.InvitationResponseDto{}
	for _, invitation := range invitations.Invitations {
		response = append(response, dto.InvitationResponseDto{
			InvitationID: invitation.InvitationID.String(),
			FormID:       invitation.FormID.String(),
			InvitedEmail: invitation.InvitedEmail,
			Status:       string(invitation.Status.InvitationStatus),
			InvitedBy:    invitation.InvitedBy.UUID.String(),
			InvitedName:  invitation.InvitedName,
		})
	}

	c.JSON(http.StatusOK, dto.InvitationListResponse{
		Status:     http.StatusOK,
		Message:    "Invitations retrieved successfully",
		Data:       response,
		Pagination: pagination,
	})
}

// VerifyInvitation verifies an invitation
// @Summary VerifyInvitation verifies an invitation
// @Description VerifyInvitation verifies an invitation
// @Tags Invitation
// @Accept json
// @Produce json
// @Param data body dto.VerifyInvitationParams true "Invitation data"
// @Success 200 {object} dto.ApiResponse[services.InvitationClaims] "Invitation verified successfully"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /invitations/verify [post]
// @type http
func (h *invitationHandler) VerifyInvitation(c *gin.Context) {

	var req dto.VerifyInvitationParams
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleError(c, err)
		return
	}

	claims, err := h.jwtService.ValidateInvitationToken(req.Token)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	res := dto.ApiResponse[services.InvitationClaims]{
		Message: "Invitation verified successfully",
		Status:  http.StatusOK,
		Data:    *claims,
	}

	c.JSON(http.StatusOK, res)
}

// GenerateAnonymousInvitationToken generates an anonymous invitation token
// @Summary GenerateAnonymousInvitationToken generates an anonymous invitation token
// @Description GenerateAnonymousInvitationToken generates an anonymous invitation token
// @Tags Invitation
// @Accept json
// @Produce json
// @Param data body dto.GenerateAnonymousInvitationTokenParams true "Invitation data"
// @Success 200 {object} dto.ApiResponse[dto.GenerateAnonymousInvitationTokenResponse] "Invitation token generated successfully"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /invitations/anonymous [post]
// @Security BearerAuth
// @type http
func (h *invitationHandler) GenerateAnonymousInvitationToken(c *gin.Context) {

	var req dto.GenerateAnonymousInvitationTokenParams
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleError(c, err)
		return
	}

	form, err := h.formService.GetSingleForm(c, req.FormID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	if form.FormAccess.FormAccess == sqlc.FormAccessRestricted {
		utils.HandleError(c, errors.New("Form is restricted"))
		return
	}
	var expiresIn time.Duration

	if form.ClosingTime.Valid {
		expiresIn = time.Until(form.ClosingTime.Time)
	} else {
		expiresIn = 24 * time.Hour
	}

	token, err := h.jwtService.GenerateAnonymousInvitationToken(req.FormID, expiresIn)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	res := dto.ApiResponse[dto.GenerateAnonymousInvitationTokenResponse]{
		Message: "Invitation token generated successfully",
		Status:  http.StatusOK,
		Data: dto.GenerateAnonymousInvitationTokenResponse{
			Token:     token,
			ExpiresIn: time.Now().Add(expiresIn).String(),
		},
	}

	c.JSON(http.StatusOK, res)
}
