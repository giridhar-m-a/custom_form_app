package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/giridhar-m-a/custom_form_app/internal/db"
	"github.com/giridhar-m-a/custom_form_app/internal/dto"
	"github.com/giridhar-m-a/custom_form_app/internal/repositories"
	"github.com/giridhar-m-a/custom_form_app/internal/services"
	"github.com/giridhar-m-a/custom_form_app/internal/utils"
)

type FormsHandler interface {
	CreateForm(ctx *gin.Context)
	CreateFormFields(ctx *gin.Context)
	GetForms(ctx *gin.Context)
	GetSingleForm(ctx *gin.Context)
}

type formHandler struct {
	formService services.FormService
}

func NewFormsHandler() FormsHandler {
	conn := db.Connection
	queries := db.Queries
	formsRepo := repositories.NewFormsRepository(queries)
	fieldsRepo := repositories.NewFormFieldsRepository(queries)
	fieldOptionsRepo := repositories.NewFormFieldOptionsRepository(queries)
	return &formHandler{formService: services.NewFormService(formsRepo, fieldsRepo, fieldOptionsRepo, conn)}
}

// CreateForm creates a new form
// @Summary      Create a new form
// @Description  Creates a new form for the authenticated user
// @Tags         Forms
// @Accept       json
// @Produce      json
// @Param        form  body      dto.CreateFormDTO  true  "Form data"
// @Success      201   {object}  object{status=string,message=string,data=dto.FormResponse}  "Form created successfully"
// @Failure      400   {object}  object{status=string,message=string}  "Invalid request payload"
// @Failure      401   {object}  object{status=string,message=string}  "Unauthorized"
// @Failure      500   {object}  object{status=string,message=string}  "Internal server error"
// @Router       /forms [post]
// @Security     BearerAuth
func (r *formHandler) CreateForm(ctx *gin.Context) {
	var form dto.CreateFormDTO
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(400, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(401, gin.H{
			"status":  "error",
			"message": "Unauthorized: user ID not found in context",
		})
		return
	}
	createdForm, err := r.formService.CreateForm(ctx, form, userID.(string))
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	response := dto.FormResponse{
		ID:          createdForm.FormID.String(),
		Title:       createdForm.FormTitle,
		Description: createdForm.FormDescription.String,
		CreatedBy:   createdForm.CreatedBy.UUID.String(),
		Status:      string(createdForm.FormStatus.FormStatus),
		CreatedAt:   createdForm.FormCreatedAt.Time.String(),
		UpdatedAt:   createdForm.FormUpdatedAt.Time.String(),
		Access:      string(createdForm.FormAccess.FormAccess),
	}

	ctx.JSON(201, gin.H{
		"status":  "success",
		"message": "Form created successfully",
		"data":    response,
	})
}

// @Summary      Create form fields
// @Description  Creates form fields for the authenticated user
// @Tags         Forms
// @Accept       json
// @Produce      json
// @Param        form  body      dto.CreateFormFieldsDTO  true  "Form fields data"
// @Success      201   {object}  object{status=string,message=string,data=dto.CreatedFormFieldDTO}  "Form fields created successfully"
// @Failure      400   {object}  object{status=string,message=string}  "Invalid request payload"
// @Failure      401   {object}  object{status=string,message=string}  "Unauthorized"
// @Failure      500   {object}  object{status=string,message=string}  "Internal server error"
// @Router       /forms/fields [post]
// @Security     BearerAuth
func (r *formHandler) CreateFormFields(ctx *gin.Context) {
	var form dto.CreateFormFieldsDTO
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(400, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(401, gin.H{
			"status":  "error",
			"message": "Unauthorized: user ID not found in context",
		})
		return
	}
	createdFormFields, err := r.formService.CreateFormFields(ctx, form, userID.(string))
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	ctx.JSON(201, gin.H{
		"status":  "success",
		"message": "Form fields created successfully",
		"data":    createdFormFields,
	})
}

// @Summary      Get forms
// @Description  Gets forms for the authenticated user
// @Tags         Forms
// @Accept       json
// @Produce      json
// @Param        query  query      dto.ListFormQuery  true  "Form query"
// @Success      200    {object}  object{status=string,message=string,data=[]dto.FormResponse, pagination=object{totalRecords=int, page=int, limit=int, totalPages=int}}  "Forms retrieved successfully"
// @Failure      400    {object}  object{status=string,message=string}  "Invalid request payload"
// @Failure      401    {object}  object{status=string,message=string}  "Unauthorized"
// @Failure      500    {object}  object{status=string,message=string}  "Internal server error"
// @Router       /forms [get]
// @Security     BearerAuth
func (r *formHandler) GetForms(ctx *gin.Context) {
	query := dto.ListFormQuery{}
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(400, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(401, gin.H{
			"status":  "error",
			"message": "Unauthorized: user ID not found in context",
		})
		return
	}

	forms, err := r.formService.GetForms(ctx, userID.(string), query)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	ctx.JSON(200, gin.H{
		"status":  "success",
		"message": "Forms retrieved successfully",
		"data":    forms.Forms,
		"pagination": gin.H{
			"totalRecords":  forms.Total,
			"page":   forms.Page,
			"limit":  forms.Limit,
			"totalPages":  forms.Pages,
		},
	})
}

// @Summary      Get single form
// @Description  Gets a single form for the authenticated user
// @Tags         Forms
// @Accept       json
// @Produce      json
// @Param        formID  path      string  true  "Form ID"
// @Success      200     {object}  object{status=string,message=string,data=dto.FormResponse}  "Form retrieved successfully"
// @Failure      400     {object}  object{status=string,message=string}  "Invalid request payload"
// @Failure      401     {object}  object{status=string,message=string}  "Unauthorized"
// @Failure      500     {object}  object{status=string,message=string}  "Internal server error"
// @Router       /forms/{formID} [get]
// @Security     BearerAuth
func (r *formHandler) GetSingleForm(ctx *gin.Context) {
	formID := ctx.Param("formID")
	form, err := r.formService.GetSingleForm(ctx, formID)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	ctx.JSON(200, gin.H{
		"status":  "success",
		"message": "Form retrieved successfully",
		"data":    form,
	})
}
