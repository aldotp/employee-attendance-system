package dto

import (
	"fmt"
	"time"
)

type AttendanceRequest struct {
	TypeAttendance string    `json:"attendance_type"`
	Latitude       float64   `json:"latitude"`
	Longitude      float64   `json:"longitude"`
	SelfieURL      string    `json:"selfie_url"`
	Status         string    `json:"status"`
	Notes          string    `json:"notes"`
	Time           time.Time `json:"time"`
}

func (a *AttendanceRequest) Validate() error {
	if a.TypeAttendance == "" {
		return fmt.Errorf("attendance type is required")
	}

	if a.TypeAttendance != "check_in" && a.TypeAttendance != "check_out" {
		return fmt.Errorf("invalid attendance type")
	}

	if a.Latitude == 0 {
		return fmt.Errorf("latitude is required")
	}

	if a.Longitude == 0 {
		return fmt.Errorf("longitude is required")
	}

	if a.SelfieURL == "" {
		return fmt.Errorf("selfie URL is required")
	}

	if a.Status == "" {
		return fmt.Errorf("status is required")
	}

	return nil

}

type AttendanceResponse struct {
	AttendanceID string    `json:"attendance_id"`
	UserID       string    `json:"user_id"`
	Time         time.Time `json:"time"`
	Status       string    `json:"status"`
}
