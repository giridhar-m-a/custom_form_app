package dto

type ErrorResponse struct {
	Status  int    `json:"status" example:"400"`
	Message string `json:"message" example:"Invalid request payload"`
}

type PaginationResponse struct {
	TotalRecords int `json:"totalRecords"`
	Page         int `json:"page"`
	Limit        int `json:"limit"`
	TotalPages   int `json:"totalPages"`
}

type ApiResponse[T any] struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}
