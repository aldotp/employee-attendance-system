package port

import (
	"context"

	"github.com/aldotp/employee-attendance-system/internal/core/domain"
	"github.com/jackc/pgx/v5"
)

type EmployeeRepository interface {
	CreateEmployee(ctx context.Context, employee *domain.Employee) (*domain.Employee, error)
	GetEmployeeByID(ctx context.Context, id string) (*domain.Employee, error)
	ListEmployees(ctx context.Context, skip, limit uint64) ([]domain.Employee, error)
	UpdateEmployee(ctx context.Context, employee *domain.Employee) (*domain.Employee, error)
	DeleteEmployee(ctx context.Context, id string) error
	FindOneByFilters(ctx context.Context, filter map[string]interface{}) (*domain.Employee, error)
	CreateEmployeeTx(ctx context.Context, tx pgx.Tx, employee *domain.Employee) (*domain.Employee, error)
}
