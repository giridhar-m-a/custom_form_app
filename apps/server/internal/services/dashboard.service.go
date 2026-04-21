package services

import (
	"context"

	"github.com/giridhar-m-a/custom_form_app/internal/dto"
	"github.com/giridhar-m-a/custom_form_app/internal/repositories"
	"github.com/google/uuid"
)

type DashboardService interface {
	GetDashboardData(ctx context.Context, userID uuid.NullUUID) (dto.DashboardResponse, error)
}

type dashboardService struct {
	dashboardRepo repositories.DashboardRepository
}

func NewDashboardService(dashboardRepo repositories.DashboardRepository) DashboardService {
	return &dashboardService{
		dashboardRepo: dashboardRepo,
	}
}

func (s *dashboardService) GetDashboardData(ctx context.Context, userID uuid.NullUUID) (dto.DashboardResponse, error) {
	totalForms, err := s.dashboardRepo.GetTotalForms(ctx, userID)
	if err != nil {
		return dto.DashboardResponse{}, err
	}
	totalSubmissions, err := s.dashboardRepo.GetTotalSubmissions(ctx, userID)
	if err != nil {
		return dto.DashboardResponse{}, err
	}
	totalActiveForms, err := s.dashboardRepo.GetTotalActiveForms(ctx, userID)
	if err != nil {
		return dto.DashboardResponse{}, err
	}
	totalClosedForms, err := s.dashboardRepo.GetTotalClosedForms(ctx, userID)
	if err != nil {
		return dto.DashboardResponse{}, err
	}
	totalInvitations, err := s.dashboardRepo.GetTotalInvitations(ctx, userID)
	if err != nil {
		return dto.DashboardResponse{}, err
	}
	formSubmissionsByMonth, err := s.dashboardRepo.GetFormSubmissionsByMonth(ctx, userID)
	if err != nil {
		return dto.DashboardResponse{}, err
	}
	submissionsByMonth := make([]dto.SubmissionsByMonth,0)
	for _, submission := range formSubmissionsByMonth {
		submissionsByMonth = append(submissionsByMonth, dto.SubmissionsByMonth{
			Month: submission.Month,
			TotalSubmissions: submission.TotalSubmissions,
		})
	}
	
	
	return dto.DashboardResponse{
		TotalForms: totalForms,
		TotalSubmissions: totalSubmissions,
		TotalActiveForms: totalActiveForms,
		TotalClosedForms: totalClosedForms,
		TotalInvitations: totalInvitations,
		SubmissionsByMonth: submissionsByMonth,
	}, nil
}