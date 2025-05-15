package dto

import "time"

type ScheduleRequest struct {
	UserID    string    `json:"user_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Notes     string    `json:"notes,omitempty"`
}

type ScheduleResponse struct {
	ScheduleID string `json:"schedule_id"`
	Status     string `json:"status"`
}

type GetWorkCalendar struct {
	Year  int `json:"year" form:"year"`
	Month int `json:"month" form:"month"`
}
