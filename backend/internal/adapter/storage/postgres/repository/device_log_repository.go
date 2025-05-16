package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/aldotp/employee-attendance-system/internal/adapter/storage/postgres"
	"github.com/aldotp/employee-attendance-system/internal/core/domain"
	"github.com/aldotp/employee-attendance-system/pkg/consts"

	"github.com/jackc/pgx/v5"
)

type DeviceLogRepository struct {
	db *postgres.DB
}

func NewDeviceLogRepository(db *postgres.DB) *DeviceLogRepository {
	return &DeviceLogRepository{
		db,
	}
}

func (dlr *DeviceLogRepository) CreateDeviceLog(ctx context.Context, deviceLog *domain.DeviceLog) (string, error) {
	query := dlr.db.QueryBuilder.Insert("device_logs").
		Columns("id", "device_id", "user_id", "action", "description", "created_at").
		Values(deviceLog.ID, deviceLog.DeviceID, deviceLog.UserID, deviceLog.Action, deviceLog.Description, deviceLog.CreatedAt).
		Suffix("RETURNING id")

	sql, args, err := query.ToSql()
	if err != nil {
		return "", err
	}

	var id string
	err = dlr.db.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (dlr *DeviceLogRepository) GetDeviceLogByID(ctx context.Context, id string) (*domain.DeviceLog, error) {
	var deviceLog domain.DeviceLog

	query := dlr.db.QueryBuilder.Select("*").
		From("device_logs").
		Where(sq.Eq{"id": id}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = dlr.db.QueryRow(ctx, sql, args...).Scan(
		&deviceLog.ID,
		&deviceLog.DeviceID,
		&deviceLog.UserID,
		&deviceLog.Action,
		&deviceLog.Description,
		&deviceLog.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, consts.ErrDataNotFound
		}
		return nil, err
	}

	return &deviceLog, nil
}

func (dlr *DeviceLogRepository) ListDeviceLogs(ctx context.Context, deviceID string, skip, limit uint64) ([]domain.DeviceLog, error) {
	var log domain.DeviceLog
	var logs []domain.DeviceLog

	if limit == 0 {
		limit = 10
	}
	if skip == 0 {
		skip = 1
	}

	query := dlr.db.QueryBuilder.Select("*").
		From("device_logs").
		Where(sq.Eq{"device_id": deviceID}).
		OrderBy("created_at DESC").
		Limit(limit).
		Offset((skip - 1) * limit)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := dlr.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&log.ID,
			&log.DeviceID,
			&log.UserID,
			&log.Action,
			&log.Description,
			&log.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, nil
}

func (dlr *DeviceLogRepository) DeleteDeviceLog(ctx context.Context, id string) error {
	query := dlr.db.QueryBuilder.Delete("device_logs").
		Where(sq.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = dlr.db.Exec(ctx, sql, args...)
	return err
}
