package domain

import "time"

type AttendanceStatus string

const (
	AttendanceStatusPresent AttendanceStatus = "present"
	AttendanceStatusLate    AttendanceStatus = "late"
	AttendanceStatusAbsent  AttendanceStatus = "absent"
)

type Attendance struct {
	ID        string           `json:"id"`
	UserID    string           `json:"user_id"`
	Time      time.Time        `json:"time"`
	Latitude  float64          `json:"latitude"`
	Longitude float64          `json:"longitude"`
	SelfieURL string           `json:"selfie_url"`
	Type      string           `json:"type"`
	Notes     string           `json:"notes"`
	Status    AttendanceStatus `json:"status"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

type GetAttendanceResponse struct {
	Name       string           `json:"name"`
	Email      string           `json:"email"`
	Department string           `json:"department"`
	Time       time.Time        `json:"time"`
	Latitude   float64          `json:"latitude"`
	Longitude  float64          `json:"longitude"`
	SelfieURL  string           `json:"selfie_url"`
	Type       string           `json:"type"`
	Notes      string           `json:"notes"`
	Status     AttendanceStatus `json:"status"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
}
