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

type LeaveRequestRepository struct {
	db *postgres.DB
}

func NewLeaveRequestRepository(db *postgres.DB) *LeaveRequestRepository {
	return &LeaveRequestRepository{
		db,
	}
}

func (lr *LeaveRequestRepository) CreateLeaveRequest(ctx context.Context, request *domain.LeaveRequest) (*domain.LeaveRequest, error) {
	query := lr.db.QueryBuilder.Insert("leave_requests").
		Columns("id", "user_id", "start_date", "end_date", "type", "reason", "status", "reviewed_by", "reviewed_at", "note", "created_at", "updated_at").
		Values(request.ID, request.UserID, request.StartDate, request.EndDate, request.Type, request.Reason, request.Status, request.ReviewedBy, request.ReviewedAt, request.Note, request.CreatedAt, request.UpdatedAt).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = lr.db.QueryRow(ctx, sql, args...).Scan(
		&request.ID,
		&request.UserID,
		&request.StartDate,
		&request.EndDate,
		&request.Type,
		&request.Reason,
		&request.Status,
		&request.ReviewedBy,
		&request.ReviewedAt,
		&request.Note,
		&request.CreatedAt,
		&request.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return request, nil
}

func (lr *LeaveRequestRepository) ApproveLeaveRequest(ctx context.Context, id string, reviewedBy string) error {
	query := lr.db.QueryBuilder.Update("leave_requests").
		Set("status", "approved").
		Set("reviewed_by", reviewedBy).
		Set("reviewed_at", time.Now()).
		Where(sq.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = lr.db.Exec(ctx, sql, args...)
	return err
}

func (lr *LeaveRequestRepository) RejectLeaveRequest(ctx context.Context, id string, reviewedBy string, note string) error {
	query := lr.db.QueryBuilder.Update("leave_requests").
		Set("status", "rejected").
		Set("reviewed_by", reviewedBy).
		Set("reviewed_at", time.Now()).
		Set("note", note).
		Where(sq.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = lr.db.Exec(ctx, sql, args...)
	return err
}

func (lr *LeaveRequestRepository) GetLeaveRequestByID(ctx context.Context, id string) (*domain.LeaveRequest, error) {
	var request domain.LeaveRequest

	query := lr.db.QueryBuilder.Select("*").
		From("leave_requests").
		Where(sq.Eq{"id": id}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = lr.db.QueryRow(ctx, sql, args...).Scan(
		&request.ID,
		&request.UserID,
		&request.StartDate,
		&request.EndDate,
		&request.Type,
		&request.Reason,
		&request.Status,
		&request.ReviewedBy,
		&request.ReviewedAt,
		&request.Note,
		&request.CreatedAt,
		&request.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, consts.ErrDataNotFound
		}
		return nil, err
	}

	return &request, nil
}

func (lr *LeaveRequestRepository) ListLeaveRequests(ctx context.Context, skip, limit uint64) ([]domain.LeaveRequest, error) {
	var request domain.LeaveRequest
	var requests []domain.LeaveRequest

	if limit == 0 {
		limit = 10
	}

	if skip == 0 {
		skip = 1
	}

	query := lr.db.QueryBuilder.Select("*").
		From("leave_requests").
		OrderBy("id").
		Limit(limit).
		Offset((skip - 1) * limit)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := lr.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&request.ID,
			&request.UserID,
			&request.StartDate,
			&request.EndDate,
			&request.Type,
			&request.Reason,
			&request.Status,
			&request.ReviewedBy,
			&request.ReviewedAt,
			&request.Note,
			&request.CreatedAt,
			&request.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		requests = append(requests, request)
	}

	return requests, nil
}

func (lr *LeaveRequestRepository) UpdateLeaveRequest(ctx context.Context, request *domain.LeaveRequest) (*domain.LeaveRequest, error) {
	query := lr.db.QueryBuilder.Update("leave_requests").
		Set("user_id", sq.Expr("COALESCE(?, employee_id)", nullString(request.UserID))).
		Set("start_date", sq.Expr("COALESCE(?, start_date)", request.StartDate)).
		Set("end_date", sq.Expr("COALESCE(?, end_date)", request.EndDate)).
		Set("type", sq.Expr("COALESCE(?, type)", nullString(string(request.Type)))).
		Set("reason", sq.Expr("COALESCE(?, reason)", nullString(request.Reason))).
		Set("status", sq.Expr("COALESCE(?, status)", nullString(string(request.Status)))).
		Set("reviewed_by", sq.Expr("COALESCE(?, reviewed_by)", nullString(request.ReviewedBy))).
		Set("reviewed_at", sq.Expr("COALESCE(?, reviewed_at)", request.ReviewedAt)).
		Set("note", sq.Expr("COALESCE(?, note)", nullString(request.Note))).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": request.ID}).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = lr.db.QueryRow(ctx, sql, args...).Scan(
		&request.ID,
		&request.UserID,
		&request.StartDate,
		&request.EndDate,
		&request.Type,
		&request.Reason,
		&request.Status,
		&request.ReviewedBy,
		&request.ReviewedAt,
		&request.Note,
		&request.CreatedAt,
		&request.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return request, nil
}

func (lr *LeaveRequestRepository) DeleteLeaveRequest(ctx context.Context, id string) error {
	query := lr.db.QueryBuilder.Delete("leave_requests").
		Where(sq.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = lr.db.Exec(ctx, sql, args...)
	return err
}
