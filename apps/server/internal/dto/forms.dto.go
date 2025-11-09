package dto

type CreateFormDTO struct {
	Title       string `json:"title" binding:"required" form:"title" message:"title is required"`
	Description string `json:"description" form:"description"`
}

type FormResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatedBy   string `json:"createdBy"`
	Status      string `json:"status"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
	Access      string `json:"access"`
}
