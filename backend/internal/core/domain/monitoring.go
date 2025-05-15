package domain

import "time"

type MonitoringReport struct {
	ID          string    `json:"id"`
	ReportType  string    `json:"report_type"`
	Data        string    `json:"data"`
	GeneratedAt time.Time `json:"generated_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type MonitoringSummary struct {
	TotalUsers      int       `json:"total_users"`
	ActiveUsers     int       `json:"active_users"`
	TotalAttendance int       `json:"total_attendance"`
	PendingLeaves   int       `json:"pending_leaves"`
	ApprovedLeaves  int       `json:"approved_leaves"`
	RejectedLeaves  int       `json:"rejected_leaves"`
	GeneratedAt     time.Time `json:"generated_at"`
}

type DashboardAnalytics struct {
	WeeklyAttendance  []int          `json:"weekly_attendance"`
	LeaveDistribution map[string]int `json:"leave_distribution"`
	UserActivity      []UserActivity `json:"user_activity"`
	GeneratedAt       time.Time      `json:"generated_at"`
}

type AttendanceReport struct {
	UserID      string    `json:"user_id"`
	TotalHours  float64   `json:"total_hours"`
	LateCount   int       `json:"late_count"`
	AbsentCount int       `json:"absent_count"`
	PeriodStart time.Time `json:"period_start"`
	PeriodEnd   time.Time `json:"period_end"`
}

type Anomaly struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	DetectedAt  time.Time `json:"detected_at"`
	Status      string    `json:"status"`
}

type ExportRequest struct {
	Format     string    `json:"format"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	ReportType string    `json:"report_type"`
}

type ExportResponse struct {
	FileURL   string    `json:"file_url"`
	ExpiresAt time.Time `json:"expires_at"`
	FileSize  int64     `json:"file_size"`
	Format    string    `json:"format"`
}

type UserActivity struct {
	UserID     string    `json:"user_id"`
	LoginCount int       `json:"login_count"`
	LastActive time.Time `json:"last_active"`
}
