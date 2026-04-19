package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giridhar-m-a/custom_form_app/internal/db"
	"github.com/giridhar-m-a/custom_form_app/internal/dto"
	"github.com/giridhar-m-a/custom_form_app/internal/repositories"
	"github.com/giridhar-m-a/custom_form_app/internal/services"
	"github.com/giridhar-m-a/custom_form_app/internal/utils"
)

type ResponseHandler interface {
	CreateSubmission(ctx *gin.Context)
	GetSubmissions(ctx *gin.Context)
	GetSingleSubmission(ctx *gin.Context)
}

type responseHandler struct {
	responseService services.ResponseService
}

func NewResponseHandler() ResponseHandler {
	conn := db.Connection
	queries := db.Queries
	repo := repositories.NewResponseRepository(queries)
	service := services.NewResponseService(repo, conn)
	return &responseHandler{responseService: service}
}

// CreateSubmission creates a new submission
// @Summary CreateSubmission creates a new submission
// @Description CreateSubmission creates a new submission
// @Tags Response
// @Accept json
// @Produce json
// @Param data body dto.CreateSubmissionRequest true "Submission data"
// @Success 200 {object} dto.ApiResponse[dto.SubmissionResponse] "Submission created successfully"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /response [post]
// @Security BearerAuth
// @type http
// @scheme bearer
func (h *responseHandler) CreateSubmission(ctx *gin.Context) {
	var req dto.CreateSubmissionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.HandleError(ctx, err)
		return
	}

	response, err := h.responseService.SubmitForm(ctx, req)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}
	apiResponse := dto.ApiResponse[dto.SubmissionResponse]{
		Message: "Form submitted successfully",
		Status:  http.StatusOK,
		Data:    response,
	}
	ctx.JSON(http.StatusOK, apiResponse)
}

// GetSubmissions retrieves submissions by form ID
// @Summary GetSubmissions retrieves submissions by form ID
// @Description GetSubmissions retrieves submissions by form ID
// @Tags Response
// @Accept json
// @Produce json
// @Param formId path string true "Form ID"
// @Param   query query dto.ResponseQuery true "Invitation List Query"
// @Success 200 {object} dto.ApiResponse[[]dto.SubmissionList] "Submissions retrieved successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid request payload"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /response/{formId} [get]
// @Security BearerAuth
// @type http
// @scheme bearer
func (h *responseHandler) GetSubmissions(ctx *gin.Context) {
	var req dto.ResponseQuery
	formId := ctx.Param("formId")
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.HandleError(ctx, err)
		return
	}
	response, err := h.responseService.GetSubmissions(ctx, dto.GetSubmissionsRequest{
		FormID: formId,
		ResponseQuery: req,
	})
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}
	apiResponse := dto.ApiResponse[[]dto.SubmissionList]{
		Message: "Submissions retrieved successfully",
		Status:  http.StatusOK,
		Data:    response.Submissions,
		Pagination: dto.PaginationResponse{
			Page:         response.Page,
			Limit:        response.Limit,
			TotalPages:   response.Pages,
			TotalRecords: response.TotalRecords,
		}}
	ctx.JSON(http.StatusOK, apiResponse)
}

// @Summary GetSingleSubmission retrieves a single submission by ID
// @Description GetSingleSubmission retrieves a single submission by ID
// @Tags Response
// @Accept json
// @Produce json
// @Param submissionId path string true "Submission ID"
// @Success 200 {object} dto.ApiResponse[dto.SubmissionResponse] "Submission retrieved successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid request payload"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /response/submission/{submissionId} [get]
// @Security BearerAuth
// @type http
// @scheme bearer
func (h *responseHandler) GetSingleSubmission(ctx *gin.Context) {
	submissionId := ctx.Param("submissionId")

	response, err := h.responseService.GetSingleSubmission(ctx, submissionId)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}
	apiResponse := dto.ApiResponse[any]{
		Message: "Submission retrieved successfully",
		Status:  http.StatusOK,
		Data:    response,
	}
	ctx.JSON(http.StatusOK, apiResponse)
}
