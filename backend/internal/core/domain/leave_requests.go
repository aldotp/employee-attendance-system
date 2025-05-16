package domain

import (
	"time"
)

type LeaveType string

const (
	Annual    LeaveType = "annual"
	Sick      LeaveType = "sick"
	Unpaid    LeaveType = "unpaid"
	Maternity LeaveType = "maternity"
	Paternity LeaveType = "paternity"
)

type LeaveStatus string

const (
	Pending  LeaveStatus = "pending"
	Approved LeaveStatus = "approved"
	Rejected LeaveStatus = "rejected"
)

type LeaveRequest struct {
	ID         string      `json:"id"`
	UserID     string      `json:"user_id"`
	StartDate  time.Time   `json:"start_date"`
	EndDate    time.Time   `json:"end_date"`
	Type       LeaveType   `json:"type"`
	Reason     string      `json:"reason"`
	Status     LeaveStatus `json:"status"`
	ReviewedBy string      `json:"reviewed_by"`
	ReviewedAt *time.Time  `json:"reviewed_at"`
	Note       string      `json:"note"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}
