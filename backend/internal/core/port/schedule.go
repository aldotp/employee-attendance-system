package port

import (
	"context"

	"github.com/aldotp/employee-attendance-system/internal/core/domain"
)

type ScheduleRepository interface {
	ListSchedules(ctx context.Context) ([]domain.Schedule, error)
	CreateSchedule(ctx context.Context, schedule *domain.Schedule) (*domain.Schedule, error)
	GetSchedule(ctx context.Context, id string) (*domain.Schedule, error)
	UpdateSchedule(ctx context.Context, id string, schedule *domain.Schedule) (*domain.Schedule, error)
	DeleteSchedule(ctx context.Context, id string) error
	RequestScheduleSwap(ctx context.Context, requestorID string, targetScheduleID string, proposedScheduleID string) error
	GetWorkCalendar(ctx context.Context, employeeID string, month int, year int) ([]domain.Schedule, error)
	GetWorkRotation(ctx context.Context, employeeID string) (*domain.Schedule, error)
}

type ScheduleService interface {
	GetWorkRotation(ctx context.Context, employeeID string) (*domain.Schedule, error)
	GetWorkCalendar(ctx context.Context, employeeID string, month int, year int) ([]domain.Schedule, error)
	RequestScheduleSwap(ctx context.Context, requestorID string, targetScheduleID string, proposedScheduleID string) error

	ListSchedules(ctx context.Context) ([]domain.Schedule, error)
	CreateSchedule(ctx context.Context, schedule *domain.Schedule) (*domain.Schedule, error)
	GetSchedule(ctx context.Context, id string) (*domain.Schedule, error)
	UpdateSchedule(ctx context.Context, id string, schedule *domain.Schedule) (*domain.Schedule, error)
	DeleteSchedule(ctx context.Context, id string) error
}
