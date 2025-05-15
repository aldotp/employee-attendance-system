package port

import (
	"context"

	"github.com/aldotp/employee-attendance-system/internal/adapter/dto"
)

type AnomalyService interface {
	DetectAnomaly(ctx context.Context, req dto.AnomalyRequest) (dto.AnomalyResponse, error)
	RecordAnomalyType(ctx context.Context, anomalyID string, anomalyType string) error
	NotifyAdminAnomaly(ctx context.Context, anomalyID string) error
	VerifyAnomaly(ctx context.Context, anomalyID string, verified bool) error
	UpdateAnomalyStatus(ctx context.Context, anomalyID string, status string) error
}
