package domain

import (
	"time"
)

type EmployeeStatus string

const (
	StatusActive     EmployeeStatus = "active"
	StatusInactive   EmployeeStatus = "inactive"
	StatusTerminated EmployeeStatus = "terminated"
	StatusSuspended  EmployeeStatus = "suspended"
)

type Employee struct {
	ID           string         `json:"id"`
	UserID       string         `json:"user_id"`
	DepartmentID string         `json:"department_id"`
	Name         string         `json:"name"`
	Location     string         `json:"location"`
	Timezone     string         `json:"timezone"`
	PhotoURL     string         `json:"photo_url"`
	Status       EmployeeStatus `json:"status"`
	JoinDate     time.Time      `json:"join_date"`
	ReportingTo  string         `json:"reporting_to"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}
