package port

import (
	"context"

	"github.com/aldotp/employee-attendance-system/internal/adapter/dto"
)

type ReportService interface {
	SelectReportType(ctx context.Context, reportType string) error
	GatherReportData(ctx context.Context, reportType string, params map[string]interface{}) (interface{}, error)
	ProcessReportData(ctx context.Context, rawData interface{}) (interface{}, error)
	GenerateReport(ctx context.Context, processedData interface{}) (dto.ReportResponse, error)
	StoreReportHistory(ctx context.Context, reportID string, data interface{}) error
}
