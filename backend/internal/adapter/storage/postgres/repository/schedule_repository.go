package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/aldotp/employee-attendance-system/internal/adapter/storage/postgres"
	"github.com/aldotp/employee-attendance-system/internal/core/domain"
	"github.com/aldotp/employee-attendance-system/pkg/consts"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ScheduleRepository struct {
	db *postgres.DB
}

func NewScheduleRepository(db *postgres.DB) *ScheduleRepository {
	return &ScheduleRepository{
		db: db,
	}
}

func (sr *ScheduleRepository) ListSchedules(ctx context.Context) ([]domain.Schedule, error) {
	var schedules []domain.Schedule

	query := sr.db.QueryBuilder.Select(
		"id",
		"user_id",
		"date",
		"shift_start",
		"shift_end",
		"break_start",
		"break_end",
		"work_location_id",
		"schedule_type",
		"created_at",
		"updated_at",
	).
		From("schedules").
		OrderBy("created_at DESC")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := sr.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var schedule domain.Schedule
		err := rows.Scan(
			&schedule.ID,
			&schedule.UserID,
			&schedule.Date,
			&schedule.ShiftStart,
			&schedule.ShiftEnd,
			&schedule.BreakStart,
			&schedule.BreakEnd,
			&schedule.WorkLocationID,
			&schedule.ScheduleType,
			&schedule.CreatedAt,
			&schedule.UpdatedAt,
		)
		if err != nil {
			return schedules, err
		}
		schedules = append(schedules, schedule)
	}

	return schedules, nil
}

func (sr *ScheduleRepository) CreateSchedule(ctx context.Context, schedule *domain.Schedule) (*domain.Schedule, error) {
	if schedule.CreatedAt.IsZero() {
		schedule.CreatedAt = time.Now()
	}
	if schedule.UpdatedAt.IsZero() {
		schedule.UpdatedAt = time.Now()
	}

	if schedule.ID == "" {
		schedule.ID = uuid.NewString()
	}

	query := sr.db.QueryBuilder.Insert("schedules").
		Columns(
			"id", "user_id", "date", "shift_start", "shift_end",
			"break_start", "break_end", "work_location_id", "schedule_type",
			"created_at", "updated_at",
		).
		Values(
			schedule.ID, schedule.UserID, schedule.Date, schedule.ShiftStart,
			schedule.ShiftEnd, schedule.BreakStart, schedule.BreakEnd,
			schedule.WorkLocationID, schedule.ScheduleType,
			schedule.CreatedAt, schedule.UpdatedAt,
		).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = sr.db.QueryRow(ctx, sql, args...).Scan(
		&schedule.ID,
		&schedule.UserID,
		&schedule.Date,
		&schedule.ShiftStart,
		&schedule.ShiftEnd,
		&schedule.BreakStart,
		&schedule.BreakEnd,
		&schedule.WorkLocationID,
		&schedule.ScheduleType,
		&schedule.CreatedAt,
		&schedule.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return schedule, nil
}

func (sr *ScheduleRepository) GetSchedule(ctx context.Context, id string) (*domain.Schedule, error) {
	var schedule domain.Schedule

	query := sr.db.QueryBuilder.Select(
		"id", "user_id", "date", "shift_start", "shift_end",
		"break_start", "break_end", "work_location_id", "schedule_type",
		"created_at", "updated_at",
	).
		From("schedules").
		Where(sq.Eq{"id": id}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = sr.db.QueryRow(ctx, sql, args...).Scan(
		&schedule.ID,
		&schedule.UserID,
		&schedule.Date,
		&schedule.ShiftStart,
		&schedule.ShiftEnd,
		&schedule.BreakStart,
		&schedule.BreakEnd,
		&schedule.WorkLocationID,
		&schedule.ScheduleType,
		&schedule.CreatedAt,
		&schedule.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, consts.ErrDataNotFound
		}
		return nil, err
	}

	return &schedule, nil
}

func (sr *ScheduleRepository) UpdateSchedule(ctx context.Context, id string, schedule *domain.Schedule) (*domain.Schedule, error) {
	schedule.UpdatedAt = time.Now()

	query := sr.db.QueryBuilder.Update("schedules").
		Set("user_id", sq.Expr("COALESCE(?, user_id)", schedule.UserID)).
		Set("date", sq.Expr("COALESCE(?, date)", schedule.Date)).
		Set("shift_start", sq.Expr("COALESCE(?, shift_start)", schedule.ShiftStart)).
		Set("shift_end", sq.Expr("COALESCE(?, shift_end)", schedule.ShiftEnd)).
		Set("break_start", sq.Expr("COALESCE(?, break_start)", schedule.BreakStart)).
		Set("break_end", sq.Expr("COALESCE(?, break_end)", schedule.BreakEnd)).
		Set("work_location_id", sq.Expr("COALESCE(?, work_location_id)", schedule.WorkLocationID)).
		Set("schedule_type", sq.Expr("COALESCE(?, schedule_type)", schedule.ScheduleType)).
		Set("updated_at", schedule.UpdatedAt).
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = sr.db.QueryRow(ctx, sql, args...).Scan(
		&schedule.ID,
		&schedule.UserID,
		&schedule.Date,
		&schedule.ShiftStart,
		&schedule.ShiftEnd,
		&schedule.BreakStart,
		&schedule.BreakEnd,
		&schedule.WorkLocationID,
		&schedule.ScheduleType,
		&schedule.CreatedAt,
		&schedule.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return schedule, nil
}

func (sr *ScheduleRepository) DeleteSchedule(ctx context.Context, id string) error {
	query := sr.db.QueryBuilder.Delete("schedules").
		Where(sq.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = sr.db.Exec(ctx, sql, args...)
	return err
}

func (sr *ScheduleRepository) RequestScheduleSwap(ctx context.Context, requestorID string, targetScheduleID string, proposedScheduleID string) error {
	query := sr.db.QueryBuilder.Insert("schedule_swap_requests").
		Columns("requestor_id", "target_schedule_id", "proposed_schedule_id", "status", "created_at").
		Values(requestorID, targetScheduleID, proposedScheduleID, "pending", time.Now())

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = sr.db.Exec(ctx, sql, args...)
	return err
}

// GetWorkCalendar retrieves all schedules for an employee in a specific month and year
func (sr *ScheduleRepository) GetWorkCalendar(ctx context.Context, userID string, month int, year int) ([]domain.Schedule, error) {
	var schedules []domain.Schedule

	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)

	query := sr.db.QueryBuilder.Select("*").
		From("schedules").
		Where(sq.And{
			sq.Eq{"user_id": userID},
			sq.GtOrEq{"start_time": startDate},
			sq.LtOrEq{"end_time": endDate},
		}).
		OrderBy("start_time ASC")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := sr.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var schedule domain.Schedule
		err := rows.Scan(
			&schedule.ID,
			&schedule.UserID,
			&schedule.Date,
			&schedule.ShiftStart,
			&schedule.ShiftEnd,
			&schedule.BreakStart,
			&schedule.BreakEnd,
			&schedule.WorkLocationID,
			&schedule.ScheduleType,
			&schedule.CreatedAt,
			&schedule.UpdatedAt,
		)
		if err != nil {
			return schedules, err
		}
		schedules = append(schedules, schedule)
	}

	return schedules, nil
}

// GetWorkRotation retrieves the current work rotation schedule for an employee
func (sr *ScheduleRepository) GetWorkRotation(ctx context.Context, employeeID string) (*domain.Schedule, error) {
	var schedule domain.Schedule

	query := sr.db.QueryBuilder.Select("*").
		From("schedules").
		Where(sq.And{
			sq.Eq{"user_id": employeeID},
			sq.GtOrEq{"end_time": time.Now()},
		}).
		OrderBy("start_time ASC").
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = sr.db.QueryRow(ctx, sql, args...).Scan(
		&schedule.ID,
		&schedule.UserID,
		&schedule.Date,
		&schedule.ShiftStart,
		&schedule.ShiftEnd,
		&schedule.BreakStart,
		&schedule.BreakEnd,
		&schedule.WorkLocationID,
		&schedule.ScheduleType,
		&schedule.CreatedAt,
		&schedule.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, consts.ErrDataNotFound
		}
		return nil, err
	}

	return &schedule, nil
}
