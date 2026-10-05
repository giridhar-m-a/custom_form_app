package dto

type SubmissionsByMonth struct {
	Month string `json:"month"`
	TotalSubmissions int64 `json:"totalSubmissions"`
}

type DashboardResponse struct {
	TotalForms int64 `json:"totalForms"`
	TotalSubmissions int64 `json:"totalSubmissions"`
	TotalActiveForms int64 `json:"totalActiveForms"`
	TotalClosedForms int64 `json:"totalClosedForms"`
	TotalInvitations int64 `json:"totalInvitations"`
	SubmissionsByMonth []SubmissionsByMonth `json:"submissionsByMonth"`
}
