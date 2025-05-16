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
	ID        string           `json:"id"`
	UserID    string           `json:"user_id"`
	Type      NotificationType `json:"type"`
	Message   string           `json:"message"`
	SendAt    time.Time        `json:"send_at"`
	CreatedAt time.Time        `json:"created_at"`
}

func NewNotification(userID string, notificationType NotificationType, message string, sendAt time.Time) *Notification {
	return &Notification{
		ID:        uuid.New().String(),
		UserID:    userID,
		Type:      notificationType,
		Message:   message,
		SendAt:    sendAt,
		CreatedAt: time.Now(),
	}
}
