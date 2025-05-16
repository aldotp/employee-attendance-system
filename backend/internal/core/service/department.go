package service

import (
	"context"

	"github.com/aldotp/employee-attendance-system/internal/core/domain"
	"github.com/aldotp/employee-attendance-system/internal/core/port"
)

type DepartmentService struct {
	repo port.DepartmentRepository
}

func NewDepartmentService(repo port.DepartmentRepository) *DepartmentService {
	return &DepartmentService{repo: repo}
}

func (s *DepartmentService) ListDepartments(ctx context.Context, skip, limit uint64) ([]domain.Department, error) {
	return s.repo.ListDepartments(ctx, 0, 0)
}
