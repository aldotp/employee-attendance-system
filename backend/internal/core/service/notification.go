package service

import (
	"context"
	"errors"
	"time"

	"github.com/aldotp/employee-attendance-system/internal/core/domain"
	"github.com/aldotp/employee-attendance-system/internal/core/port"
	"github.com/google/uuid"
)

type NotificationService struct {
	repo     port.NotificationRepository
	UserRepo port.UserRepository
}

func NewNotificationService(repo port.NotificationRepository, UserRepo port.UserRepository) *NotificationService {
	return &NotificationService{
		repo:     repo,
		UserRepo: UserRepo,
	}
}

func (ns *NotificationService) CreateNotification(ctx context.Context, notif *domain.Notification) (string, error) {

	_, err := ns.UserRepo.GetUserByID(ctx, notif.UserID)
	if err != nil {
		if err.Error() == "data not found" {
			return "", errors.New("employee not found")
		}
		return "", err
	}

	if notif.ID == "" {
		notif.ID = uuid.New().String()
	}

	if notif.CreatedAt.IsZero() {
		notif.CreatedAt = time.Now()
	}

	if notif.SendAt.IsZero() {
		notif.SendAt = time.Now()
	}

	return ns.repo.CreateNotification(ctx, notif)
}

func (ns *NotificationService) GetNotificationByID(ctx context.Context, id string) (*domain.Notification, error) {
	return ns.repo.GetNotificationByID(ctx, id)
}

func (ns *NotificationService) ListNotifications(ctx context.Context, employeeID string, skip, limit uint64) ([]domain.Notification, error) {
	if skip == 0 {
		skip = 0
	}

	if limit == 0 {
		limit = 10
	}

	return ns.repo.ListNotifications(ctx, employeeID, skip, limit)
}

func (ns *NotificationService) UpdateNotificationStatus(ctx context.Context, notificationID string, status string) error {
	return ns.repo.UpdateNotificationStatus(ctx, notificationID, status)
}

func (ns *NotificationService) DeleteNotification(ctx context.Context, notificationID string) error {
	return ns.repo.DeleteNotification(ctx, notificationID)
}
