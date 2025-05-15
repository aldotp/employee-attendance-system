package service

import (
	"context"
	"time"

	"github.com/aldotp/employee-attendance-system/internal/core/domain"
	"github.com/aldotp/employee-attendance-system/internal/core/port"
)

type MonitoringService struct {
	repo port.MonitoringRepository
}

func NewMonitoringService(repo port.MonitoringRepository) *MonitoringService {
	return &MonitoringService{
		repo: repo,
	}
}

func (ms *MonitoringService) GetReports(ctx context.Context) ([]domain.MonitoringReport, error) {
	return ms.repo.GetReports(ctx)
}

func (ms *MonitoringService) GetSummary(ctx context.Context) (*domain.MonitoringSummary, error) {
	return ms.repo.GetSummary(ctx)
}

func (ms *MonitoringService) GetDashboardAnalytics(ctx context.Context) (*domain.DashboardAnalytics, error) {
	// Get basic analytics data
	summary, err := ms.repo.GetSummary(ctx)
	if err != nil {
		return nil, err
	}

	// Get recent reports
	// reports, err := ms.repo.GetReports(ctx)
	// if err != nil {
	// 	return nil, err
	// }

	// Create dashboard analytics
	analytics := &domain.DashboardAnalytics{
		WeeklyAttendance: []int{summary.TotalAttendance, 0, 0, 0, 0, 0, 0}, // Placeholder
		LeaveDistribution: map[string]int{
			"pending":  summary.PendingLeaves,
			"approved": summary.ApprovedLeaves,
			"rejected": summary.RejectedLeaves,
		},
		GeneratedAt: time.Now(),
	}

	return analytics, nil
}

func (ms *MonitoringService) GenerateAttendanceReport(ctx context.Context) (*domain.AttendanceReport, error) {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

	report := &domain.AttendanceReport{
		PeriodStart: startOfMonth,
		PeriodEnd:   endOfMonth,
	}

	return report, nil
}

func (ms *MonitoringService) DetectAnomalies(ctx context.Context) ([]domain.Anomaly, error) {
	return []domain.Anomaly{}, nil
}

func (ms *MonitoringService) ExportData(ctx context.Context, req domain.ExportRequest) (*domain.ExportResponse, error) {
	exportResponse := &domain.ExportResponse{
		FileURL:   "https://example.com/exports/" + req.ReportType + "." + req.Format,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		FileSize:  1024,
		Format:    req.Format,
	}

	return exportResponse, nil
}
