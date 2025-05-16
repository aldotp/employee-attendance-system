package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/aldotp/employee-attendance-system/internal/adapter/storage/postgres"
	"github.com/aldotp/employee-attendance-system/internal/core/domain"
	"github.com/aldotp/employee-attendance-system/pkg/consts"
	"github.com/jackc/pgx/v5"
)

type WorkLocationRepository struct {
	db *postgres.DB
}

func NewWorkLocationRepository(db *postgres.DB) *WorkLocationRepository {
	return &WorkLocationRepository{
		db,
	}
}

func (wlr *WorkLocationRepository) CreateWorkLocation(ctx context.Context, location *domain.WorkLocation) (*domain.WorkLocation, error) {
	query := wlr.db.QueryBuilder.Insert("work_locations").
		Columns("id", "name", "address", "city", "state", "country", "postal_code", "timezone", "created_at", "updated_at").
		Values(location.ID, location.Name, location.Address, location.City, location.State, location.Country, location.PostalCode, location.Timezone, location.CreatedAt, location.UpdatedAt).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = wlr.db.QueryRow(ctx, sql, args...).Scan(
		&location.ID,
		&location.Name,
		&location.Address,
		&location.City,
		&location.State,
		&location.Country,
		&location.PostalCode,
		&location.Timezone,
		&location.CreatedAt,
		&location.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return location, nil
}

func (wlr *WorkLocationRepository) DeleteWorkLocation(ctx context.Context, id string) error {
	query := wlr.db.QueryBuilder.Delete("work_locations").
		Where(sq.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = wlr.db.Exec(ctx, sql, args...)
	return err
}

func (wlr *WorkLocationRepository) GetWorkLocationByID(ctx context.Context, id string) (*domain.WorkLocation, error) {
	var location domain.WorkLocation

	query := wlr.db.QueryBuilder.Select("*").
		From("work_locations").
		Where(sq.Eq{"id": id}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = wlr.db.QueryRow(ctx, sql, args...).Scan(
		&location.ID,
		&location.Name,
		&location.Address,
		&location.City,
		&location.State,
		&location.Country,
		&location.PostalCode,
		&location.Timezone,
		&location.CreatedAt,
		&location.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, consts.ErrDataNotFound
		}
		return nil, err
	}

	return &location, nil
}

func (wlr *WorkLocationRepository) ListWorkLocations(ctx context.Context, skip, limit uint64) ([]domain.WorkLocation, error) {
	var location domain.WorkLocation
	var locations []domain.WorkLocation

	if limit == 0 {
		limit = 10
	}

	if skip == 0 {
		skip = 1
	}

	query := wlr.db.QueryBuilder.Select("*").
		From("work_locations").
		OrderBy("id").
		Limit(limit).
		Offset((skip - 1) * limit)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := wlr.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&location.ID,
			&location.Name,
			&location.Address,
			&location.City,
			&location.State,
			&location.Country,
			&location.PostalCode,
			&location.Timezone,
			&location.CreatedAt,
			&location.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		locations = append(locations, location)
	}

	return locations, nil
}

func (wlr *WorkLocationRepository) UpdateWorkLocation(ctx context.Context, location *domain.WorkLocation) (*domain.WorkLocation, error) {
	query := wlr.db.QueryBuilder.Update("work_locations").
		Set("name", sq.Expr("COALESCE(?, name)", nullString(location.Name))).
		Set("address", sq.Expr("COALESCE(?, address)", nullString(location.Address))).
		Set("city", sq.Expr("COALESCE(?, city)", nullString(location.City))).
		Set("state", sq.Expr("COALESCE(?, state)", nullString(location.State))).
		Set("country", sq.Expr("COALESCE(?, country)", nullString(location.Country))).
		Set("postal_code", sq.Expr("COALESCE(?, postal_code)", nullString(location.PostalCode))).
		Set("timezone", sq.Expr("COALESCE(?, timezone)", nullString(location.Timezone))).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": location.ID}).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = wlr.db.QueryRow(ctx, sql, args...).Scan(
		&location.ID,
		&location.Name,
		&location.Address,
		&location.City,
		&location.State,
		&location.Country,
		&location.PostalCode,
		&location.Timezone,
		&location.CreatedAt,
		&location.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return location, nil
}
