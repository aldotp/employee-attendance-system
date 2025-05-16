package port

import (
	"context"

	"github.com/aldotp/employee-attendance-system/internal/core/domain"
)

type DepartmentRepository interface {
	CreateDepartment(ctx context.Context, department *domain.Department) (*domain.Department, error)
	GetDepartmentByID(ctx context.Context, id string) (*domain.Department, error)
	ListDepartments(ctx context.Context, skip, limit uint64) ([]domain.Department, error)
	UpdateDepartment(ctx context.Context, department *domain.Department) (*domain.Department, error)
	DeleteDepartment(ctx context.Context, id string) error
	GetDepartmentByName(ctx context.Context, name string) (*domain.Department, error)
}
