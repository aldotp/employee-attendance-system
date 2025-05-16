package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/aldotp/employee-attendance-system/internal/adapter/storage/postgres"
	"github.com/aldotp/employee-attendance-system/internal/core/domain"
	"github.com/aldotp/employee-attendance-system/pkg/consts"
	"github.com/jackc/pgx/v5"
)

type NotificationRepository struct {
	db *postgres.DB
}

func NewNotificationRepository(db *postgres.DB) *NotificationRepository {
	return &NotificationRepository{
		db,
	}
}

func (nr *NotificationRepository) CreateNotification(ctx context.Context, notif *domain.Notification) (string, error) {
	query := nr.db.QueryBuilder.Insert("notifications").
		Columns("id", "user_id", "type", "message", "send_at", "created_at").
		Values(notif.ID, notif.UserID, notif.Type, notif.Message, notif.SendAt, notif.CreatedAt).
		Suffix("RETURNING id")

	sql, args, err := query.ToSql()
	if err != nil {
		return "", err
	}

	var id string
	err = nr.db.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (nr *NotificationRepository) GetNotificationByID(ctx context.Context, id string) (*domain.Notification, error) {
	var notif domain.Notification

	query := nr.db.QueryBuilder.Select("*").
		From("notifications").
		Where(sq.Eq{"id": id}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = nr.db.QueryRow(ctx, sql, args...).Scan(
		&notif.ID,
		&notif.UserID,
		&notif.Type,
		&notif.Message,
		&notif.SendAt,
		&notif.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, consts.ErrDataNotFound
		}
		return nil, err
	}

	return &notif, nil
}

func (nr *NotificationRepository) ListNotifications(ctx context.Context, userID string, skip, limit uint64) ([]domain.Notification, error) {
	var notif domain.Notification
	var notifs []domain.Notification

	if limit == 0 {
		limit = 10
	}
	if skip == 0 {
		skip = 1
	}

	query := nr.db.QueryBuilder.Select("*").
		From("notifications").
		Where(sq.Eq{"user_id": userID}).
		OrderBy("created_at DESC").
		Limit(limit).
		Offset((skip - 1) * limit)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := nr.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&notif.ID,
			&notif.UserID,
			&notif.Type,
			&notif.Message,
			&notif.SendAt,
			&notif.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		notifs = append(notifs, notif)
	}

	return notifs, nil
}

func (nr *NotificationRepository) DeleteNotification(ctx context.Context, notificationID string) error {
	query := nr.db.QueryBuilder.Delete("notifications").
		Where(sq.Eq{"id": notificationID})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = nr.db.Exec(ctx, sql, args...)
	return err
}

func (nr *NotificationRepository) UpdateNotificationStatus(ctx context.Context, notificationID string, status string) error {
	query := nr.db.QueryBuilder.Update("notifications").
		Set("status", status).
		Where(sq.Eq{"id": notificationID})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = nr.db.Exec(ctx, sql, args...)
	return err
}
