package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/giridhar-m-a/custom_form_app/internal/db"
	"github.com/giridhar-m-a/custom_form_app/internal/dto"
	"github.com/giridhar-m-a/custom_form_app/internal/repositories"
	serializers "github.com/giridhar-m-a/custom_form_app/internal/serialisers"
	"github.com/giridhar-m-a/custom_form_app/internal/services"
	"github.com/giridhar-m-a/custom_form_app/internal/utils"
)

type FormsHandler interface {
	CreateForm(ctx *gin.Context)
	CreateFormFields(ctx *gin.Context)
	GetForms(ctx *gin.Context)
	GetSingleForm(ctx *gin.Context)
	UpdateForm(ctx *gin.Context)
	DeleteForm(ctx *gin.Context)
	GetFormFields(ctx *gin.Context)
	UpdateFormFields(ctx *gin.Context)
	GetFormForResponse(ctx *gin.Context)
	GetFormFieldsForResponse(ctx *gin.Context)
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
		utils.HandleError(ctx, err)
		return
	}

	userID, exists := ctx.Get("userID")
	if !exists {
		utils.HandleError(ctx, errors.New("user ID not found in context"))
		return
	}
	createdForm, err := r.formService.CreateForm(ctx, form, userID.(string))
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	response := dto.FormResponse{
		ID:                  createdForm.FormID.String(),
		Title:               createdForm.FormTitle,
		Description:         utils.NullStringToString(createdForm.FormDescription),
		CreatedBy:           utils.NullUUIDToString(createdForm.CreatedBy),
		Status:              string(createdForm.FormStatus.FormStatus),
		CreatedAt:           utils.NullTimeToString(createdForm.FormCreatedAt),
		UpdatedAt:           utils.NullTimeToString(createdForm.FormUpdatedAt),
		Access:              string(createdForm.FormAccess.FormAccess),
		SchedulingID:        utils.NullUUIDToString(createdForm.SchedulingID),
		ScheduledTime:       utils.NullTimeToString(createdForm.ScheduledTime),
		ClosingTime:         utils.NullTimeToString(createdForm.ClosingTime),
		IsScheduleCompleted: utils.NullBoolToBool(createdForm.IsScheduleCompleted, false),
		IsScheduled:         utils.NullBoolToBool(createdForm.IsScheduled, false),
	}

	ctx.JSON(201, gin.H{
		"status":  201,
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
// @Success      201   {object}  object{status=string,message=string,data=[]dto.FormFieldResponseDto}  "Form fields created successfully"
// @Failure      400   {object}  object{status=string,message=string}  "Invalid request payload"
// @Failure      401   {object}  object{status=string,message=string}  "Unauthorized"
// @Failure      500   {object}  object{status=string,message=string}  "Internal server error"
// @Router       /forms/fields [post]
// @Security     BearerAuth
func (r *formHandler) CreateFormFields(ctx *gin.Context) {
	var form dto.CreateFormFieldsDTO
	if err := ctx.ShouldBind(&form); err != nil {
		utils.HandleError(ctx, err)
		return
	}

	userID, exists := ctx.Get("userID")
	if !exists {
		utils.HandleError(ctx, errors.New("user ID not found in context"))
		return
	}
	createdFormFields, err := r.formService.CreateFormFields(ctx, form, userID.(string))
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	serializedFields := []dto.FormFieldResponseDto{}
	for _, field := range createdFormFields {
		serializedFields = append(serializedFields, serializers.MapCreatedFormFieldToResponse(field))
	}

	ctx.JSON(201, gin.H{
		"status":  201,
		"message": "Form fields created successfully",
		"data":    serializedFields,
	})
}

// @Summary      Update form fields
// @Description  Creates form fields for the authenticated user
// @Tags         Forms
// @Accept       json
// @Produce      json
// @Param        form  body      dto.UpdateFormFieldsDTO  true  "Form fields data"
// @Success      201   {object}  object{status=string,message=string,data=[]dto.FormFieldResponseDto}  "Form fields Updated successfully"
// @Failure      400   {object}  object{status=string,message=string}  "Invalid request payload"
// @Failure      401   {object}  object{status=string,message=string}  "Unauthorized"
// @Failure      500   {object}  object{status=string,message=string}  "Internal server error"
// @Router       /forms/fields [Patch]
// @Security     BearerAuth
func (r *formHandler) UpdateFormFields(ctx *gin.Context) {
	var fields dto.UpdateFormFieldsDTO
	if err := ctx.ShouldBind(&fields); err != nil {
		utils.HandleError(ctx, err)
		return
	}

	updatedFormFields, err := r.formService.UpdateFormFields(ctx, fields)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	serializedFields := []dto.FormFieldResponseDto{}
	for _, field := range updatedFormFields {
		serializedFields = append(serializedFields, serializers.MapCreatedFormFieldToResponse(field))
	}

	ctx.JSON(201, gin.H{
		"status":  201,
		"message": "Form fields updated successfully",
		"data":    serializedFields,
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
		utils.HandleError(ctx, err)
		return
	}
	userID, exists := ctx.Get("userID")
	if !exists {
		utils.HandleError(ctx, errors.New("user ID not found in context"))
		return
	}

	forms, err := r.formService.GetForms(ctx, userID.(string), query)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	var formResponse []dto.FormResponse
	for _, form := range forms.Forms {
		var invitationScheduleGap *int32
		if form.InvitationScheduleGap.Valid {
			invitationScheduleGap = &form.InvitationScheduleGap.Int32
		}

		formResponse = append(formResponse, dto.FormResponse{
			ID:                    form.FormID.String(),
			Title:                 form.FormTitle,
			CreatedAt:             utils.NullTimeToString(form.FormCreatedAt),
			UpdatedAt:             utils.NullTimeToString(form.FormUpdatedAt),
			Description:           utils.NullStringToString(form.FormDescription),
			CreatedBy:             utils.NullUUIDToString(form.CreatedBy),
			Status:                string(form.FormStatus.FormStatus),
			Access:                string(form.FormAccess.FormAccess),
			SchedulingID:          utils.NullUUIDToString(form.SchedulingID),
			ScheduledTime:         utils.NullTimeToString(form.ScheduledTime),
			ClosingTime:           utils.NullTimeToString(form.ClosingTime),
			IsScheduleCompleted:   utils.NullBoolToBool(form.IsScheduleCompleted, false),
			IsScheduled:           utils.NullBoolToBool(form.IsScheduled, false),
			InvitationScheduleGap: invitationScheduleGap,
		})
	}

	ctx.JSON(200, gin.H{
		"status":  200,
		"message": "Forms retrieved successfully",
		"data":    formResponse,
		"pagination": gin.H{
			"totalRecords": forms.Total,
			"page":         forms.Page,
			"limit":        forms.Limit,
			"totalPages":   forms.Pages,
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
	if formID == "" {
		utils.HandleError(ctx, errors.New("form id is required"))
		return
	}
	form, err := r.formService.GetSingleForm(ctx, formID)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}
	var invitationScheduleGap *int32
	if form.InvitationScheduleGap.Valid {
		invitationScheduleGap = &form.InvitationScheduleGap.Int32
	}

	ctx.JSON(200, gin.H{
		"status":  200,
		"message": "Form retrieved successfully",
		"data": dto.FormResponse{
			ID:                    form.FormID.String(),
			Title:                 form.FormTitle,
			Description:           utils.NullStringToString(form.FormDescription),
			CreatedBy:             utils.NullUUIDToString(form.CreatedBy),
			Status:                string(form.FormStatus.FormStatus),
			CreatedAt:             utils.NullTimeToString(form.FormCreatedAt),
			UpdatedAt:             utils.NullTimeToString(form.FormUpdatedAt),
			Access:                string(form.FormAccess.FormAccess),
			SchedulingID:          utils.NullUUIDToString(form.SchedulingID),
			ScheduledTime:         utils.NullTimeToString(form.ScheduledTime),
			ClosingTime:           utils.NullTimeToString(form.ClosingTime),
			IsScheduleCompleted:   utils.NullBoolToBool(form.IsScheduleCompleted, false),
			IsScheduled:           utils.NullBoolToBool(form.IsScheduled, false),
			InvitationScheduleGap: invitationScheduleGap,
		},
	})
}

// @Summary      Update form
// @Description  Updates a form for the authenticated user
// @Tags         Forms
// @Accept       json
// @Produce      json
// @Param        formID  path      string  true  "Form ID"
// @Param        form    body      dto.UpdateFormDTO  true  "Form fields data"
// @Success      200     {object}  object{status=string,message=string,data=dto.FormResponse}  "Form updated successfully"
// @Failure      400     {object}  object{status=string,message=string}  "Invalid request payload"
// @Failure      401     {object}  object{status=string,message=string}  "Unauthorized"
// @Failure      500     {object}  object{status=string,message=string}  "Internal server error"
// @Router       /forms/{formID} [patch]
// @Security     BearerAuth
func (r *formHandler) UpdateForm(ctx *gin.Context) {
	formID := ctx.Param("formID")

	if formID == "" {
		utils.HandleError(ctx, errors.New("form id is required"))
		return
	}

	var form dto.UpdateFormDTO
	if err := ctx.ShouldBind(&form); err != nil {
		utils.HandleError(ctx, err)
		return
	}
	updatedForm, err := r.formService.UpdateForm(ctx, form, formID)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	ctx.JSON(200, gin.H{
		"status":  200,
		"message": "Form updated successfully",
		"data": dto.FormResponse{
			ID:                  updatedForm.FormID.String(),
			Title:               updatedForm.FormTitle,
			Description:         utils.NullStringToString(updatedForm.FormDescription),
			CreatedBy:           utils.NullUUIDToString(updatedForm.CreatedBy),
			Status:              string(updatedForm.FormStatus.FormStatus),
			CreatedAt:           utils.NullTimeToString(updatedForm.FormCreatedAt),
			UpdatedAt:           utils.NullTimeToString(updatedForm.FormUpdatedAt),
			Access:              string(updatedForm.FormAccess.FormAccess),
			SchedulingID:        utils.NullUUIDToString(updatedForm.SchedulingID),
			ScheduledTime:       utils.NullTimeToString(updatedForm.ScheduledTime),
			ClosingTime:         utils.NullTimeToString(updatedForm.ClosingTime),
			IsScheduleCompleted: utils.NullBoolToBool(updatedForm.IsScheduleCompleted, false),
			IsScheduled:         utils.NullBoolToBool(updatedForm.IsScheduled, false),
		},
	})
}

// @Summary      Delete form
// @Description  Deletes a form for the authenticated user
// @Tags         Forms
// @Accept       json
// @Produce      json
// @Param        formID  path      string  true  "Form ID"
// @Success      200     {object}  object{status=string,message=string,data=dto.FormResponse}  "Form deleted successfully"
// @Failure      400     {object}  object{status=string,message=string}  "Invalid request payload"
// @Failure      401     {object}  object{status=string,message=string}  "Unauthorized"
// @Failure      500     {object}  object{status=string,message=string}  "Internal server error"
// @Router       /forms/{formID} [delete]
// @Security     BearerAuth
func (r *formHandler) DeleteForm(ctx *gin.Context) {
	formID := ctx.Param("formID")
	if formID == "" {
		utils.HandleError(ctx, errors.New("form id is required"))
		return
	}
	deletedForm, err := r.formService.DeleteForm(ctx, formID)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	ctx.JSON(200, gin.H{
		"status":  200,
		"message": "Form deleted successfully",
		"data": dto.FormResponse{
			ID:          deletedForm.FormID.String(),
			Title:       deletedForm.FormTitle,
			Description: deletedForm.FormDescription.String,
			CreatedBy:   deletedForm.CreatedBy.UUID.String(),
			Status:      string(deletedForm.FormStatus.FormStatus),
			CreatedAt:   deletedForm.FormCreatedAt.Time.String(),
			UpdatedAt:   deletedForm.FormUpdatedAt.Time.String(),
			Access:      string(deletedForm.FormAccess.FormAccess),
		},
	})
}

// @Summary      Get form fields
// @Description  Gets form fields for the authenticated user
// @Tags         Forms
// @Accept       json
// @Produce      json
// @Param        formID  path      string  true  "Form ID"
// @Success      200     {object}  object{status=string,message=string,data=[]dto.FormFieldResponseDto}  "Form fields retrieved successfully"
// @Failure      400     {object}  object{status=string,message=string}  "Invalid request payload"
// @Failure      401     {object}  object{status=string,message=string}  "Unauthorized"
// @Failure      500     {object}  object{status=string,message=string}  "Internal server error"
// @Router       /forms/fields/{formID} [get]
// @Security     BearerAuth
func (r *formHandler) GetFormFields(ctx *gin.Context) {
	formID := ctx.Param("formID")
	if formID == "" {
		utils.HandleError(ctx, errors.New("form id is required"))
		return
	}
	formFields, err := r.formService.GetFormFieldsByFormId(ctx, formID)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	serializedFields := []dto.FormFieldResponseDto{}
	for _, field := range formFields {
		serializedFields = append(serializedFields, serializers.MapCreatedFormFieldToResponse(field))
	}

	ctx.JSON(200, gin.H{
		"status":  200,
		"message": "Form fields retrieved successfully",
		"data":    serializedFields,
	})
}

// @Summary      Get form by ID for response
// @Description  Gets a single form for response submission
// @Tags         Forms
// @Accept       json
// @Produce      json
// @Success      200     {object}  object{status=string,message=string,data=dto.FormResponse}  "Form retrieved successfully"
// @Failure      400     {object}  object{status=string,message=string}  "Invalid request payload"
// @Failure      401     {object}  object{status=string,message=string}  "Unauthorized"
// @Failure      500     {object}  object{status=string,message=string}  "Internal server error"
// @Router       /forms/response [get]
// @Security     BearerAuth
func (r *formHandler) GetFormForResponse(ctx *gin.Context) {
	formID, exists := ctx.Get("formID")
	if !exists {
		utils.HandleError(ctx, errors.New("form id is required"))
		return
	}
	form, err := r.formService.GetSingleForm(ctx, formID.(string))
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	ctx.JSON(200, gin.H{
		"status":  200,
		"message": "Form retrieved successfully",
		"data": dto.FormResponse{
			ID:                  form.FormID.String(),
			Title:               form.FormTitle,
			Description:         utils.NullStringToString(form.FormDescription),
			CreatedBy:           utils.NullUUIDToString(form.CreatedBy),
			Status:              string(form.FormStatus.FormStatus),
			CreatedAt:           utils.NullTimeToString(form.FormCreatedAt),
			UpdatedAt:           utils.NullTimeToString(form.FormUpdatedAt),
			Access:              string(form.FormAccess.FormAccess),
			SchedulingID:        utils.NullUUIDToString(form.SchedulingID),
			ScheduledTime:       utils.NullTimeToString(form.ScheduledTime),
			ClosingTime:         utils.NullTimeToString(form.ClosingTime),
			IsScheduleCompleted: utils.NullBoolToBool(form.IsScheduleCompleted, false),
			IsScheduled:         utils.NullBoolToBool(form.IsScheduled, false),
		},
	})
}

// @Summary      Get form fields for response
// @Description  Gets form fields for the response submission
// @Tags         Forms
// @Accept       json
// @Produce      json
// @Success      200     {object}  object{status=string,message=string,data=[]dto.FormFieldResponseDto}  "Form fields retrieved successfully"
// @Failure      400     {object}  object{status=string,message=string}  "Invalid request payload"
// @Failure      401     {object}  object{status=string,message=string}  "Unauthorized"
// @Failure      500     {object}  object{status=string,message=string}  "Internal server error"
// @Router       /forms/fields/response [get]
// @Security     BearerAuth
func (r *formHandler) GetFormFieldsForResponse(ctx *gin.Context) {
	formID, exists := ctx.Get("formID")
	if !exists {
		utils.HandleError(ctx, errors.New("form id is required"))
		return
	}
	formFields, err := r.formService.GetFormFieldsByFormId(ctx, formID.(string))
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	serializedFields := []dto.FormFieldResponseDto{}
	for _, field := range formFields {
		serializedFields = append(serializedFields, serializers.MapCreatedFormFieldToResponse(field))
	}

	ctx.JSON(200, gin.H{
		"status":  200,
		"message": "Form fields retrieved successfully",
		"data":    serializedFields,
	})
}
