package service

import (
	"context"
	"errors"
	"time"

	"github.com/aldotp/employee-attendance-system/internal/adapter/dto"
	"github.com/aldotp/employee-attendance-system/internal/core/domain"
	"github.com/aldotp/employee-attendance-system/internal/core/port"
	"github.com/google/uuid"
)

type AttendanceService struct {
	repo port.AttendanceRepository
	// add other dependencies as needed (e.g., notification, logger)
}

func NewAttendanceService(repo port.AttendanceRepository) *AttendanceService {
	return &AttendanceService{
		repo: repo,
	}
}

// OpenAttendance handles check-in/check-out logic
func (s *AttendanceService) OpenAttendance(ctx context.Context, req dto.AttendanceRequest, userID string) (dto.AttendanceResponse, error) {
	var typeAttendance string

	if req.TypeAttendance == "check_in" {
		typeAttendance = "check_in"
	} else if req.TypeAttendance == "check_out" {
		typeAttendance = "check_out"
	} else {
		return dto.AttendanceResponse{}, errors.New("invalid attendance status")
	}

	attendance := &domain.Attendance{
		ID:        uuid.New().String(),
		Type:      typeAttendance,
		UserID:    userID,
		Time:      req.Time,
		SelfieURL: req.SelfieURL,
		Status:    domain.AttendanceStatusPresent,
		Notes:     req.Notes,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	created, err := s.repo.CreateAttendance(ctx, attendance)
	if err != nil {
		return dto.AttendanceResponse{}, err
	}

	return dto.AttendanceResponse{
		AttendanceID: created.ID,
		UserID:       created.UserID,
		Time:         created.Time,
		Status:       string(created.Status),
	}, nil
}

// ValidateSchedule checks if the user is scheduled for a given scheduleID
func (s *AttendanceService) ValidateSchedule(ctx context.Context, userID, scheduleID string) (bool, error) {
	// Example logic: always return true (replace with real validation)
	return true, nil
}

// CheckGPSLocation validates if the user's location is within allowed bounds
func (s *AttendanceService) CheckGPSLocation(ctx context.Context, userID string, lat, lng float64) (bool, error) {
	// Example logic: always return true (replace with real validation)
	return true, nil
}

// ValidateRadius checks if the user is within a certain radius for attendance
func (s *AttendanceService) ValidateRadius(ctx context.Context, userID string, lat, lng float64) (bool, error) {
	// Example logic: always return true (replace with real validation)
	return true, nil
}

// RecordAttendance records an attendance event
func (s *AttendanceService) RecordAttendance(ctx context.Context, req dto.AttendanceRequest, userID string) error {
	attendance := &domain.Attendance{
		UserID:    userID,
		Type:      req.TypeAttendance,
		Time:      req.Time,
		SelfieURL: req.SelfieURL,
		Status:    domain.AttendanceStatusPresent,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err := s.repo.CreateAttendance(ctx, attendance)
	return err
}

// SendAttendanceNotification sends a notification to the user
func (s *AttendanceService) SendAttendanceNotification(ctx context.Context, userID string) error {
	// Example logic: not implemented
	return errors.New("SendAttendanceNotification not implemented")
}

func (s *AttendanceService) ListAttendances(ctx context.Context, req domain.ListAttendanceRequest) ([]domain.GetAttendanceResponse, error) {
	return s.repo.ListAttendances(ctx, req.Page, req.Limit, req.Date, req.Type)
}

func (s *AttendanceService) GetAttendanceByID(ctx context.Context, id string) (*domain.Attendance, error) {
	return s.repo.GetAttendanceByID(ctx, id)
}

func (s *AttendanceService) UpdateAttendance(ctx context.Context, id string, req dto.AttendanceRequest) (*domain.Attendance, error) {
	attendance, err := s.repo.GetAttendanceByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided in the request
	if req.Time != (time.Time{}) {
		attendance.Time = req.Time
	}
	if req.TypeAttendance != "" {
		attendance.Type = req.TypeAttendance
	}
	if req.Status != "" {
		attendance.Status = domain.AttendanceStatus(req.Status)
	}
	if req.Notes != "" {
		attendance.Notes = req.Notes
	}
	if req.Latitude != 0 {
		attendance.Latitude = req.Latitude
	}
	if req.Longitude != 0 {
		attendance.Longitude = req.Longitude
	}
	if req.SelfieURL != "" {
		attendance.SelfieURL = req.SelfieURL
	}
	attendance.UpdatedAt = time.Now()

	return s.repo.UpdateAttendance(ctx, attendance)
}

func (s *AttendanceService) DeleteAttendance(ctx context.Context, id string) error {
	return s.repo.DeleteAttendance(ctx, id)
}

func (s *AttendanceService) GetAttendanceHistory(ctx context.Context, employeeID, startDate, endDate string) ([]domain.Attendance, error) {
	return s.repo.GetAttendanceHistory(ctx, employeeID, startDate, endDate)
}

func (s *AttendanceService) GetUsersAttendanceStatus(ctx context.Context, date string) (map[string]bool, error) {
	return s.repo.GetUsersAttendanceStatus(ctx, date)
}
