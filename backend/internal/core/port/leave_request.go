package port

import (
	"context"

	"github.com/aldotp/employee-attendance-system/internal/adapter/dto"
	"github.com/aldotp/employee-attendance-system/internal/core/domain"
)

type LeaveRequestRepository interface {
	CreateLeaveRequest(ctx context.Context, request *domain.LeaveRequest) (*domain.LeaveRequest, error)
	GetLeaveRequestByID(ctx context.Context, id string) (*domain.LeaveRequest, error)
	ListLeaveRequests(ctx context.Context, skip, limit uint64) ([]domain.LeaveRequest, error)
	UpdateLeaveRequest(ctx context.Context, request *domain.LeaveRequest) (*domain.LeaveRequest, error)
	DeleteLeaveRequest(ctx context.Context, id string) error
	ApproveLeaveRequest(ctx context.Context, id string, reviewedBy string) error
	RejectLeaveRequest(ctx context.Context, id string, reviewedBy string, note string) error
}

type LeaveService interface {
	SubmitLeaveRequest(ctx context.Context, req dto.LeaveRequest) (dto.LeaveResponse, error)
	ValidateLeaveBalance(ctx context.Context, userID string, leaveType string) (bool, error)
	NotifyApprover(ctx context.Context, leaveID string) error
	ReviewLeaveSubmission(ctx context.Context, leaveID, approverID string, approve bool, note string) error
	UpdateLeaveStatus(ctx context.Context, leaveID string, status string) error
	UpdateAttendanceForLeave(ctx context.Context, leaveID string) error
	SendLeaveNotification(ctx context.Context, userID string) error
	ApproveLeave(ctx context.Context, leaveID string) error
	RejectLeave(ctx context.Context, leaveID string, reason string) error
	GetLeaveBalance(ctx context.Context, leaveType string) (float64, error)

	ListLeaves(ctx context.Context) ([]domain.LeaveRequest, error)
	CreateLeave(ctx context.Context, req dto.LeaveRequest) (*domain.LeaveRequest, error)
	GetLeaveByID(ctx context.Context, id string) (*domain.LeaveRequest, error)
	UpdateLeave(ctx context.Context, id string, req dto.LeaveRequest) (*domain.LeaveRequest, error)
	DeleteLeave(ctx context.Context, id string) error
}
