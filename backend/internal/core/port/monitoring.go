package port

import (
	"context"

	"github.com/aldotp/employee-attendance-system/internal/core/domain"
)

type MonitoringService interface {
	GetReports(ctx context.Context) ([]domain.MonitoringReport, error)
	GetSummary(ctx context.Context) (*domain.MonitoringSummary, error)
	GetDashboardAnalytics(ctx context.Context) (*domain.DashboardAnalytics, error)
	GenerateAttendanceReport(ctx context.Context) (*domain.AttendanceReport, error)
	DetectAnomalies(ctx context.Context) ([]domain.Anomaly, error)
	ExportData(ctx context.Context, req domain.ExportRequest) (*domain.ExportResponse, error)
}

type MonitoringRepository interface {
	GetReports(ctx context.Context) ([]domain.MonitoringReport, error)
	GetSummary(ctx context.Context) (*domain.MonitoringSummary, error)
	GetDashboardAnalytics(ctx context.Context) (*domain.DashboardAnalytics, error)
	GenerateAttendanceReport(ctx context.Context) (*domain.AttendanceReport, error)
	DetectAnomalies(ctx context.Context) ([]domain.Anomaly, error)
	ExportData(ctx context.Context, req domain.ExportRequest) (*domain.ExportResponse, error)
}
