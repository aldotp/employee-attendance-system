package service

import (
	"context"
	"fmt"
	"time"

	"github.com/aldotp/employee-attendance-system/internal/adapter/dto"
	"github.com/aldotp/employee-attendance-system/internal/core/domain"
	"github.com/aldotp/employee-attendance-system/internal/core/port"
	"github.com/aldotp/employee-attendance-system/pkg/consts"
	"github.com/aldotp/employee-attendance-system/pkg/util"
	"github.com/google/uuid"
)

type LeaveService struct {
	repo            port.LeaveRequestRepository
	notificationSvc port.NotificationService
}

func NewLeaveService(repo port.LeaveRequestRepository, notificationService port.NotificationService) *LeaveService {
	return &LeaveService{
		repo:            repo,
		notificationSvc: notificationService,
	}
}

func (s *LeaveService) SubmitLeaveRequest(ctx context.Context, req dto.LeaveRequest) (dto.LeaveResponse, error) {

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return dto.LeaveResponse{}, fmt.Errorf("invalid start date format")
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return dto.LeaveResponse{}, fmt.Errorf("invalid end date format")
	}

	if startDate.After(endDate) {
		return dto.LeaveResponse{}, fmt.Errorf("start date must be before end date")
	}

	leave := &domain.LeaveRequest{
		ID:        uuid.New().String(),
		UserID:    req.UserID,
		Type:      domain.LeaveType(req.Type),
		StartDate: startDate,
		EndDate:   endDate,
		Reason:    req.Reason,
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	created, err := s.repo.CreateLeaveRequest(ctx, leave)
	if err != nil {
		return dto.LeaveResponse{}, err
	}

	go s.notificationSvc.CreateNotification(ctx, &domain.Notification{
		ID:        uuid.New().String(),
		UserID:    leave.UserID,
		Type:      domain.NotificationTypeInfo,
		Message:   fmt.Sprintf("Your leave request for %s has been submitted.", leave.Type),
		CreatedAt: time.Now(),
	})

	return dto.LeaveResponse{
		LeaveID:   created.ID,
		Type:      string(created.Type),
		StartDate: created.StartDate,
		EndDate:   created.EndDate,
		Status:    string(created.Status),
		Reason:    created.Reason,
	}, nil
}

func (s *LeaveService) ValidateLeaveBalance(ctx context.Context, userID string, leaveType string) (bool, error) {
	// Implementasi dummy, ganti sesuai kebutuhan
	return true, nil
}

func (s *LeaveService) NotifyApprover(ctx context.Context, leaveID string) error {
	// Implementasi dummy, ganti sesuai kebutuhan
	return nil
}

func (s *LeaveService) ReviewLeaveSubmission(ctx context.Context, leaveID, approverID string, approve bool, note string) error {
	leave, err := s.repo.GetLeaveRequestByID(ctx, leaveID)
	if err != nil {
		return err
	}

	if approve {
		leave.Status = "approved"
		err := s.repo.ApproveLeaveRequest(ctx, leaveID, approverID)
		if err != nil {
			return err
		}
	} else {
		leave.Status = "rejected"
		err := s.repo.RejectLeaveRequest(ctx, leaveID, approverID, note)
		if err != nil {
			return err
		}

	}
	go s.notificationSvc.CreateNotification(ctx, &domain.Notification{
		ID:        uuid.New().String(),
		UserID:    leave.UserID,
		Type:      domain.NotificationTypeInfo,
		Message:   fmt.Sprintf("Your leave request for %s has been %s.", leave.Type, leave.Status),
		CreatedAt: time.Now(),
	})

	return nil

}

func (s *LeaveService) UpdateLeaveStatus(ctx context.Context, leaveID string, status string) error {
	leave, err := s.repo.GetLeaveRequestByID(ctx, leaveID)
	if err != nil {
		return err
	}
	leave.Status = domain.LeaveStatus(status)
	leave.UpdatedAt = time.Now()
	_, err = s.repo.UpdateLeaveRequest(ctx, leave)
	return err
}

func (s *LeaveService) UpdateAttendanceForLeave(ctx context.Context, leaveID string) error {
	// Implementasi dummy, ganti sesuai kebutuhan
	return nil
}

func (s *LeaveService) SendLeaveNotification(ctx context.Context, userID string, leaveType string, status string) error {
	go s.notificationSvc.CreateNotification(ctx, &domain.Notification{
		ID:        uuid.New().String(),
		UserID:    userID,
		Type:      domain.NotificationTypeInfo,
		Message:   fmt.Sprintf("Your leave request for %s has been %s.", leaveType, status),
		CreatedAt: time.Now(),
	})

	return nil
}

// CRUD untuk handler
func (s *LeaveService) ListLeaves(ctx context.Context) ([]domain.LeaveRequest, error) {
	return s.repo.ListLeaveRequests(ctx, 0, 100)
}

func (s *LeaveService) CreateLeave(ctx context.Context, req dto.LeaveRequest) (*domain.LeaveRequest, error) {

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format")
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date format")
	}

	if startDate.After(endDate) {
		return nil, fmt.Errorf("start date must be before end date")
	}

	leave := &domain.LeaveRequest{
		ID:        uuid.New().String(),
		UserID:    req.UserID,
		Type:      domain.LeaveType(req.Type),
		StartDate: startDate,
		EndDate:   endDate,
		Reason:    req.Reason,
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return s.repo.CreateLeaveRequest(ctx, leave)
}

func (s *LeaveService) GetLeaveByID(ctx context.Context, id string) (*domain.LeaveRequest, error) {
	return s.repo.GetLeaveRequestByID(ctx, id)
}

func (s *LeaveService) UpdateLeave(ctx context.Context, id string, req dto.LeaveRequest) (*domain.LeaveRequest, error) {
	leave, err := s.repo.GetLeaveRequestByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Type != "" {
		leave.Type = domain.LeaveType(req.Type)
	}

	if req.StartDate != "" {
		startDate, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return nil, fmt.Errorf("invalid start date format")
		}

		leave.StartDate = startDate
	}

	if req.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end date format")
		}
		leave.EndDate = endDate
	}

	if req.Reason != "" {
		leave.Reason = req.Reason
	}
	leave.UpdatedAt = time.Now()
	return s.repo.UpdateLeaveRequest(ctx, leave)
}

func (s *LeaveService) DeleteLeave(ctx context.Context, id string) error {
	return s.repo.DeleteLeaveRequest(ctx, id)
}
func (s *LeaveService) ApproveLeave(ctx context.Context, leaveID string) error {
	userSession := util.GetAuthPayload(ctx, consts.AuthorizationKey)
	if userSession == nil {
		return fmt.Errorf("user session not found")
	}

	leave, err := s.repo.GetLeaveRequestByID(ctx, leaveID)
	if err != nil {
		return fmt.Errorf("failed to get leave request: %w", err)
	}

	if leave.Status != "pending" {
		return fmt.Errorf("leave request is not in pending status")
	}

	err = s.repo.ApproveLeaveRequest(ctx, leaveID, userSession.UserID)
	if err != nil {
		return fmt.Errorf("failed to approve leave request: %w", err)
	}

	if err := s.SendLeaveNotification(ctx, leave.UserID, string(leave.Type), "approve"); err != nil {
		fmt.Printf("failed to send notification: %v", err)
	}

	return nil
}

func (s *LeaveService) RejectLeave(ctx context.Context, leaveID string, reason string) error {
	userSession := util.GetAuthPayload(ctx, consts.AuthorizationKey)
	if userSession == nil {
		return fmt.Errorf("user session not found")
	}

	leave, err := s.repo.GetLeaveRequestByID(ctx, leaveID)
	if err != nil {
		return fmt.Errorf("failed to get leave request: %w", err)
	}

	if leave.Status != "pending" {
		return fmt.Errorf("leave request is not in pending status")
	}

	err = s.repo.RejectLeaveRequest(ctx, leaveID, userSession.UserID, reason)
	if err != nil {
		return fmt.Errorf("failed to reject leave request: %w", err)
	}

	if err := s.SendLeaveNotification(ctx, leave.UserID, string(leave.Type), "reject"); err != nil {
		fmt.Printf("failed to send notification: %v", err)
	}

	return nil
}

func (s *LeaveService) GetLeaveBalance(ctx context.Context, leaveType string) (float64, error) {
	userSession := util.GetAuthPayload(ctx, consts.AuthorizationKey)
	if userSession == nil {
		return 0, fmt.Errorf("user session not found")
	}

	leaves, err := s.repo.ListLeaveRequests(ctx, 0, 100)
	if err != nil {
		return 0, fmt.Errorf("failed to get leave requests: %w", err)
	}

	var usedLeaves float64
	currentYear := time.Now().Year()

	for _, leave := range leaves {
		if leave.UserID == userSession.UserID &&
			string(leave.Type) == leaveType &&
			leave.Status == "approved" &&
			leave.StartDate.Year() == currentYear {

			duration := leave.EndDate.Sub(leave.StartDate)
			usedLeaves += duration.Hours() / 24
		}
	}

	const defaultAnnualBalance = 12.0

	remainingBalance := defaultAnnualBalance - usedLeaves
	if remainingBalance < 0 {
		remainingBalance = 0
	}

	return remainingBalance, nil
}
