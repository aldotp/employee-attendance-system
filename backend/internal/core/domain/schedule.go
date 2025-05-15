package domain

import "time"

type Schedule struct {
	ID             string    `json:"id"`
	EmployeeID     string    `json:"employee_id"`
	Date           time.Time `json:"date"`
	ShiftStart     string    `json:"shift_start"`
	ShiftEnd       string    `json:"shift_end"`
	BreakStart     string    `json:"break_start,omitempty"`
	BreakEnd       string    `json:"break_end,omitempty"`
	WorkLocationID string    `json:"work_location_id,omitempty"`
	ScheduleType   string    `json:"schedule_type"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type ScheduleSwapRequest struct {
	ScheduleID1 string `json:"schedule_id_1"`
	ScheduleID2 string `json:"schedule_id_2"`
}
