package port

import (
	"context"

	"github.com/aldotp/employee-attendance-system/internal/core/domain"
)

type WorkLocationRepository interface {
	CreateWorkLocation(ctx context.Context, location *domain.WorkLocation) (*domain.WorkLocation, error)
	GetWorkLocationByID(ctx context.Context, id string) (*domain.WorkLocation, error)
	ListWorkLocations(ctx context.Context, skip, limit uint64) ([]domain.WorkLocation, error)
	UpdateWorkLocation(ctx context.Context, location *domain.WorkLocation) (*domain.WorkLocation, error)
	DeleteWorkLocation(ctx context.Context, id string) error
}
