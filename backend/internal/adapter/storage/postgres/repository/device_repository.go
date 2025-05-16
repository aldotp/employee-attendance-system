package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/aldotp/employee-attendance-system/internal/adapter/storage/postgres"
	"github.com/aldotp/employee-attendance-system/internal/core/domain"
	"github.com/aldotp/employee-attendance-system/pkg/consts"

	"github.com/jackc/pgx/v5"
)

type DeviceRepository struct {
	db *postgres.DB
}

func NewDeviceRepository(db *postgres.DB) *DeviceRepository {
	return &DeviceRepository{
		db,
	}
}

func (dr *DeviceRepository) CreateDevice(ctx context.Context, device *domain.Device) (*domain.Device, error) {
	query := dr.db.QueryBuilder.Insert("devices").
		Columns("id", "name", "type", "location", "status", "last_check", "created_at", "updated_at").
		Values(device.ID, device.Name, device.Type, device.Location, device.Status, device.LastCheck, device.CreatedAt, device.UpdatedAt).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = dr.db.QueryRow(ctx, sql, args...).Scan(
		&device.ID,
		&device.Name,
		&device.Type,
		&device.Location,
		&device.Status,
		&device.LastCheck,
		&device.CreatedAt,
		&device.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return device, nil
}

func (dr *DeviceRepository) UpdateDevice(ctx context.Context, device *domain.Device) (*domain.Device, error) {
	query := dr.db.QueryBuilder.Update("devices").
		Set("name", device.Name).
		Set("type", device.Type).
		Set("status", device.Status).
		Set("updated_at", device.UpdatedAt).
		Where(sq.Eq{"id": device.ID}).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = dr.db.QueryRow(ctx, sql, args...).Scan(
		&device.ID,
		&device.Name,
		&device.Type,
		&device.Status,
		&device.CreatedAt,
		&device.UpdatedAt,
		// add other fields as needed
	)
	if err != nil {
		return nil, err
	}

	return device, nil
}

func (dr *DeviceRepository) GetDeviceByID(ctx context.Context, id string) (*domain.Device, error) {
	var device domain.Device

	query := dr.db.QueryBuilder.Select("*").
		From("devices").
		Where(sq.Eq{"id": id}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = dr.db.QueryRow(ctx, sql, args...).Scan(
		&device.ID,
		&device.Name,
		&device.Type,
		&device.Location,
		&device.Status,
		&device.LastCheck,
		&device.CreatedAt,
		&device.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, consts.ErrDataNotFound
		}
		return nil, err
	}

	return &device, nil
}

func (dr *DeviceRepository) DeleteDevice(ctx context.Context, id string) error {
	query := dr.db.QueryBuilder.Delete("devices").
		Where(sq.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = dr.db.Exec(ctx, sql, args...)
	return err
}

func (dr *DeviceRepository) ListDevices(ctx context.Context, page, limit uint64) ([]domain.Device, error) {
	var devices []domain.Device

	if limit == 0 {
		limit = 10
	}
	if page == 0 {
		page = 1
	}

	query := dr.db.QueryBuilder.Select("*").
		From("devices").
		OrderBy("id").
		Limit(limit).
		Offset((page - 1) * limit)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := dr.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var device domain.Device
		err := rows.Scan(
			&device.ID,
			&device.Name,
			&device.Type,
			&device.Status,
			&device.CreatedAt,
			&device.UpdatedAt,
			// add other fields as needed
		)
		if err != nil {
			return nil, err
		}
		devices = append(devices, device)
	}

	return devices, nil
}
