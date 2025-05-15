package dto

import "time"

type LeaveRequest struct {
	EmployeeID string `json:"employee_id"`
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	Type       string `json:"type"`
	Reason     string `json:"reason"`
}

type LeaveResponse struct {
	LeaveID   string    `json:"leave_id"`
	Status    string    `json:"status"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Type      string    `json:"type"`
	Reason    string    `json:"reason"`
}

type RejectLeaveRequest struct {
	Reason string `json:"reason"`
}

type LeaveBalanceRequest struct {
	Type string `json:"type"`
}
