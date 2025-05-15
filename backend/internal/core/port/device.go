package port

import (
	"context"

	"github.com/aldotp/employee-attendance-system/internal/core/domain"
)

type DeviceRepository interface {
	CreateDevice(ctx context.Context, device *domain.Device) (*domain.Device, error)
	GetDeviceByID(ctx context.Context, id string) (*domain.Device, error)
	ListDevices(ctx context.Context, skip, limit uint64) ([]domain.Device, error)
	UpdateDevice(ctx context.Context, device *domain.Device) (*domain.Device, error)
	DeleteDevice(ctx context.Context, id string) error
}
