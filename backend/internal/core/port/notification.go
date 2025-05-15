package port

import (
	"context"

	"github.com/aldotp/employee-attendance-system/internal/core/domain"
)

type NotificationRepository interface {
	CreateNotification(ctx context.Context, notif *domain.Notification) (string, error)
	GetNotificationByID(ctx context.Context, id string) (*domain.Notification, error)
	ListNotifications(ctx context.Context, employeeID string, skip, limit uint64) ([]domain.Notification, error)
	UpdateNotificationStatus(ctx context.Context, notificationID string, status string) error
	DeleteNotification(ctx context.Context, notificationID string) error
}

type NotificationService interface {
	CreateNotification(ctx context.Context, notif *domain.Notification) (string, error)
	GetNotificationByID(ctx context.Context, id string) (*domain.Notification, error)
	ListNotifications(ctx context.Context, employeeID string, skip, limit uint64) ([]domain.Notification, error)
	UpdateNotificationStatus(ctx context.Context, notificationID string, status string) error
	DeleteNotification(ctx context.Context, notificationID string) error
}
