package dto

type AnomalyRequest struct {
	UserID       string `json:"user_id"`
	AttendanceID string `json:"attendance_id"`
	Type         string `json:"type"`
	Description  string `json:"description"`
}

type AnomalyResponse struct {
	AnomalyID string `json:"anomaly_id"`
	Status    string `json:"status"`
}
