package port

import (
	"context"

	"github.com/aldotp/employee-attendance-system/internal/core/domain"
)

type DeviceLogRepository interface {
	CreateDeviceLog(ctx context.Context, log *domain.DeviceLog) (string, error)
	GetDeviceLogByID(ctx context.Context, id string) (*domain.DeviceLog, error)
	ListDeviceLogs(ctx context.Context, deviceID string, skip, limit uint64) ([]domain.DeviceLog, error)
	DeleteDeviceLog(ctx context.Context, id string) error
}
