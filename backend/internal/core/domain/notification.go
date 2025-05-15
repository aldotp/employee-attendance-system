package domain

import (
	"time"

	"github.com/google/uuid"
)

type NotificationType string

const (
	NotificationTypeReminder NotificationType = "reminder"
	NotificationTypeWarning  NotificationType = "warning"
	NotificationTypeInfo     NotificationType = "info"
)

type Notification struct {
	ID         string           `json:"id"`
	EmployeeID string           `json:"employee_id"`
	Type       NotificationType `json:"type"`
	Message    string           `json:"message"`
	SendAt     time.Time        `json:"send_at"`
	CreatedAt  time.Time        `json:"created_at"`
}

func NewNotification(employeeID string, notificationType NotificationType, message string, sendAt time.Time) *Notification {
	return &Notification{
		ID:         uuid.New().String(),
		EmployeeID: employeeID,
		Type:       notificationType,
		Message:    message,
		SendAt:     sendAt,
		CreatedAt:  time.Now(),
	}
}
