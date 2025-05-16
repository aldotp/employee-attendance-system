package repository

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/aldotp/employee-attendance-system/internal/adapter/storage/postgres"
	"github.com/aldotp/employee-attendance-system/internal/core/domain"
	"github.com/aldotp/employee-attendance-system/pkg/consts"

	"github.com/jackc/pgx/v5"
)

type AttendanceRepository struct {
	db *postgres.DB
}

func NewAttendanceRepository(db *postgres.DB) *AttendanceRepository {
	return &AttendanceRepository{
		db,
	}
}

func (ar *AttendanceRepository) CreateAttendance(ctx context.Context, attendance *domain.Attendance) (*domain.Attendance, error) {
	query := ar.db.QueryBuilder.Insert("attendances").
		Columns(
			"id", "user_id", "time", "type", "status", "notes",
			"latitude", "longitude", "selfie_url", "created_at", "updated_at",
		).
		Values(
			attendance.ID, attendance.UserID, attendance.Time, attendance.Type, attendance.Status, attendance.Notes,
			attendance.Latitude, attendance.Longitude, attendance.SelfieURL, attendance.CreatedAt, attendance.UpdatedAt,
		).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = ar.db.QueryRow(ctx, sql, args...).Scan(
		&attendance.ID,
		&attendance.UserID,
		&attendance.Time,
		&attendance.Type,
		&attendance.Status,
		&attendance.Notes,
		&attendance.Latitude,
		&attendance.Longitude,
		&attendance.SelfieURL,
		&attendance.CreatedAt,
		&attendance.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return attendance, nil
}

func (ar *AttendanceRepository) GetAttendanceHistory(ctx context.Context, userID string, startDate, endDate string) ([]domain.Attendance, error) {
	var attendances []domain.Attendance

	// Convert string dates to time.Time
	startTime, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format: %v", err)
	}
	endTime, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date format: %v", err)
	}

	query := ar.db.QueryBuilder.Select("a.id, a.user_id, a.time, a.type, a.status, a.notes, a.latitude, a.longitude, a.selfie_url, a.created_at, a.updated_at").
		From("attendances a").
		Join("employees e ON e.user_id = a.user_id").
		Where(sq.And{
			sq.Eq{"a.user_id": userID},
			sq.GtOrEq{"a.time": startTime},
			sq.LtOrEq{"a.time": endTime},
		}).
		OrderBy("time DESC")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := ar.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var attendance domain.Attendance
		err := rows.Scan(
			&attendance.ID,
			&attendance.UserID,
			&attendance.Time,
			&attendance.Type,
			&attendance.Status,
			&attendance.Notes,
			&attendance.Latitude,
			&attendance.Longitude,
			&attendance.SelfieURL,
			&attendance.CreatedAt,
			&attendance.UpdatedAt,
		)
		if err != nil {
			return attendances, err
		}
		attendances = append(attendances, attendance)
	}

	return attendances, nil
}

func (ar *AttendanceRepository) UpdateAttendance(ctx context.Context, attendance *domain.Attendance) (*domain.Attendance, error) {
	query := ar.db.QueryBuilder.Update("attendances").
		Set("user_id", sq.Expr("COALESCE(?, user_id)", attendance.UserID)).
		Set("time", sq.Expr("COALESCE(?, time)", attendance.Time)).
		Set("type", sq.Expr("COALESCE(?, type)", attendance.Type)).
		Set("status", sq.Expr("COALESCE(?, status)", attendance.Status)).
		Set("notes", sq.Expr("COALESCE(?, notes)", attendance.Notes)).
		Set("latitude", sq.Expr("COALESCE(?, latitude)", attendance.Latitude)).
		Set("longitude", sq.Expr("COALESCE(?, longitude)", attendance.Longitude)).
		Set("selfie_url", sq.Expr("COALESCE(?, selfie_url)", attendance.SelfieURL)).
		Set("updated_at", attendance.UpdatedAt).
		Where(sq.Eq{"id": attendance.ID}).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = ar.db.QueryRow(ctx, sql, args...).Scan(
		&attendance.ID,
		&attendance.UserID,
		&attendance.Time,
		&attendance.Type,
		&attendance.Status,
		&attendance.Notes,
		&attendance.Latitude,
		&attendance.Longitude,
		&attendance.SelfieURL,
		&attendance.CreatedAt,
		&attendance.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return attendance, nil
}

func (ar *AttendanceRepository) GetAttendanceByID(ctx context.Context, id string) (*domain.Attendance, error) {
	var attendance domain.Attendance

	query := ar.db.QueryBuilder.Select("*").
		From("attendances").
		Where(sq.Eq{"id": id}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = ar.db.QueryRow(ctx, sql, args...).Scan(
		&attendance.ID,
		&attendance.UserID,
		&attendance.Time,
		&attendance.Type,
		&attendance.Status,
		&attendance.Notes,
		&attendance.Latitude,
		&attendance.Longitude,
		&attendance.SelfieURL,
		&attendance.CreatedAt,
		&attendance.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, consts.ErrDataNotFound
		}
		return nil, err
	}

	return &attendance, nil
}

func (ar *AttendanceRepository) DeleteAttendance(ctx context.Context, id string) error {
	query := ar.db.QueryBuilder.Delete("attendances").
		Where(sq.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = ar.db.Exec(ctx, sql, args...)
	return err
}

func (ar *AttendanceRepository) ListAttendances(ctx context.Context, page, limit uint64, date string, attendanceType string) ([]domain.GetAttendanceResponse, error) {
	var data []domain.GetAttendanceResponse

	if limit == 0 {
		limit = 10
	}
	if page == 0 {
		page = 1
	}

	query := ar.db.QueryBuilder.Select(
		"COALESCE(e.name, '') AS employee_name",
		"COALESCE(u.email, '') AS email",
		"COALESCE(d.name, '') AS department_name",
		"a.time",
		"a.latitude",
		"a.longitude",
		"a.selfie_url",
		"a.type",
		"a.notes",
		"a.status",
		"a.created_at",
		"a.updated_at",
	).
		From("attendances a").
		Join("users u ON u.id::uuid = a.user_id::uuid").
		Join("employees e ON u.id::uuid = e.user_id::uuid").
		LeftJoin("departments d ON d.id::uuid = e.department_id::uuid").
		OrderBy("a.time DESC").
		Limit(limit).
		Offset((page - 1) * limit)

	if date != "" {
		query = query.Where(sq.Eq{"DATE(a.time)": date})
	}

	if attendanceType != "" {
		query = query.Where(sq.Eq{"a.type": attendanceType})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := ar.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var attendance domain.GetAttendanceResponse
		err := rows.Scan(
			&attendance.Name,
			&attendance.Email,
			&attendance.Department,
			&attendance.Time,
			&attendance.Latitude,
			&attendance.Longitude,
			&attendance.SelfieURL,
			&attendance.Type,
			&attendance.Notes,
			&attendance.Status,
			&attendance.CreatedAt,
			&attendance.UpdatedAt,
		)
		if err != nil {
			return data, err
		}
		data = append(data, attendance)
	}

	return data, nil
}

func (r *AttendanceRepository) GetUsersAttendanceStatus(ctx context.Context, date string) (map[string]bool, error) {
	query := `
        SELECT 
            e.name,
            CASE WHEN a.id IS NOT NULL THEN true ELSE false END AS attendance_status
        FROM users u
		LEFT JOIN employees e ON u.id = e.user_id
        LEFT JOIN attendances a 
            ON u.id = a.user_id 
            AND DATE(a.time) = $1
    `

	rows, err := r.db.Query(ctx, query, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	statusMap := make(map[string]bool)
	for rows.Next() {
		var name string
		var status bool
		if err := rows.Scan(&name, &status); err != nil {
			return nil, err
		}
		statusMap[name] = status
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return statusMap, nil
}
