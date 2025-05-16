package port

import (
	"context"

	"github.com/aldotp/employee-attendance-system/internal/adapter/dto"
	"github.com/aldotp/employee-attendance-system/internal/core/domain"
)

type AttendanceRepository interface {
	CreateAttendance(ctx context.Context, attendance *domain.Attendance) (*domain.Attendance, error)
	GetAttendanceByID(ctx context.Context, id string) (*domain.Attendance, error)
	ListAttendances(ctx context.Context, page, limit uint64, date string, attendanceType string) ([]domain.GetAttendanceResponse, error)
	UpdateAttendance(ctx context.Context, attendance *domain.Attendance) (*domain.Attendance, error)
	DeleteAttendance(ctx context.Context, id string) error
	GetAttendanceHistory(ctx context.Context, employeeID string, startDate, endDate string) ([]domain.Attendance, error)
	GetUsersAttendanceStatus(ctx context.Context, date string) (map[string]bool, error)
}

type AttendanceService interface {
	OpenAttendance(ctx context.Context, req dto.AttendanceRequest, userID string) (dto.AttendanceResponse, error)
	ValidateSchedule(ctx context.Context, userID, scheduleID string) (bool, error)
	CheckGPSLocation(ctx context.Context, userID string, lat, lng float64) (bool, error)
	ValidateRadius(ctx context.Context, userID string, lat, lng float64) (bool, error)
	RecordAttendance(ctx context.Context, req dto.AttendanceRequest, userID string) error
	SendAttendanceNotification(ctx context.Context, userID string) error

	ListAttendances(context.Context, domain.ListAttendanceRequest) ([]domain.GetAttendanceResponse, error)
	GetAttendanceByID(ctx context.Context, id string) (*domain.Attendance, error)
	UpdateAttendance(ctx context.Context, id string, req dto.AttendanceRequest) (*domain.Attendance, error)
	DeleteAttendance(ctx context.Context, id string) error
}
