package service

import (
	"context"

	"github.com/aldotp/employee-attendance-system/internal/core/domain"
	"github.com/aldotp/employee-attendance-system/internal/core/port"
	"github.com/google/uuid"
)

type ScheduleService struct {
	repo port.ScheduleRepository
}

func NewScheduleService(repo port.ScheduleRepository) *ScheduleService {
	return &ScheduleService{
		repo: repo,
	}
}

func (s *ScheduleService) ListSchedules(ctx context.Context) ([]domain.Schedule, error) {
	return s.repo.ListSchedules(ctx)
}

func (s *ScheduleService) CreateSchedule(ctx context.Context, schedule *domain.Schedule) (*domain.Schedule, error) {
	if schedule.ID == "" {
		schedule.ID = uuid.New().String()
	}

	return s.repo.CreateSchedule(ctx, schedule)
}

func (s *ScheduleService) GetSchedule(ctx context.Context, id string) (*domain.Schedule, error) {
	return s.repo.GetSchedule(ctx, id)
}

func (s *ScheduleService) UpdateSchedule(ctx context.Context, id string, schedule *domain.Schedule) (*domain.Schedule, error) {
	return s.repo.UpdateSchedule(ctx, id, schedule)
}

func (s *ScheduleService) DeleteSchedule(ctx context.Context, id string) error {
	return s.repo.DeleteSchedule(ctx, id)
}
func (s *ScheduleService) GetWorkRotation(ctx context.Context, employeeID string) (*domain.Schedule, error) {
	return s.repo.GetWorkRotation(ctx, employeeID)
}

func (s *ScheduleService) GetWorkCalendar(ctx context.Context, userID string, month int, year int) ([]domain.Schedule, error) {
	return s.repo.GetWorkCalendar(ctx, userID, month, year)
}

func (s *ScheduleService) RequestScheduleSwap(ctx context.Context, requestorID string, targetScheduleID string, proposedScheduleID string) error {
	return s.repo.RequestScheduleSwap(ctx, requestorID, targetScheduleID, proposedScheduleID)
}
