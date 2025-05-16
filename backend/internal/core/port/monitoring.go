package port

import (
	"context"

	"github.com/aldotp/employee-attendance-system/internal/core/domain"
)

type MonitoringService interface {
	GetReports(ctx context.Context) ([]domain.MonitoringReport, error)
	GetSummary(ctx context.Context, date string) (*domain.MonitoringSummary, error)
	GetDashboardAnalytics(context.Context, string) (*domain.DashboardAnalytics, error)
	GenerateAttendanceReport(ctx context.Context) (*domain.AttendanceReport, error)
	DetectAnomalies(ctx context.Context) ([]domain.Anomaly, error)
	ExportData(ctx context.Context, req domain.ExportRequest) (*domain.ExportResponse, error)
}

type MonitoringRepository interface {
	GetReports(ctx context.Context) ([]domain.MonitoringReport, error)
	GetSummary(ctx context.Context, date string) (*domain.MonitoringSummary, error)
	GetDashboardAnalytics(ctx context.Context) (*domain.DashboardAnalytics, error)
	GenerateAttendanceReport(ctx context.Context) (*domain.AttendanceReport, error)
	DetectAnomalies(ctx context.Context) ([]domain.Anomaly, error)
	ExportData(ctx context.Context, req domain.ExportRequest) (*domain.ExportResponse, error)
	InsertReport(ctx context.Context, report *domain.MonitoringReport) error
	InsertReports(ctx context.Context, reports []domain.MonitoringReport) error
	GetReportByDate(ctx context.Context, input domain.ExportRequest) ([]domain.MonitoringReport, error)
}
