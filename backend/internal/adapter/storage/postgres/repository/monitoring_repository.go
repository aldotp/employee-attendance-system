package repository

import (
	"context"
	"time"

	"github.com/aldotp/employee-attendance-system/internal/adapter/storage/postgres"
	"github.com/aldotp/employee-attendance-system/internal/core/domain"
)

type MonitoringRepository struct {
	db *postgres.DB
}

func NewMonitoringRepository(db *postgres.DB) *MonitoringRepository {
	return &MonitoringRepository{
		db: db,
	}
}

func (mr *MonitoringRepository) GetReports(ctx context.Context) ([]domain.MonitoringReport, error) {
	query := mr.db.QueryBuilder.Select("*").From("monitoring_reports")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := mr.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []domain.MonitoringReport
	for rows.Next() {
		var report domain.MonitoringReport
		err := rows.Scan(
			&report.ID,
			&report.ReportType,
			&report.Data,
			&report.GeneratedAt,
			&report.CreatedAt,
			&report.UpdatedAt,
		)
		if err != nil {
			return reports, err
		}
		reports = append(reports, report)
	}

	return reports, nil
}

func (mr *MonitoringRepository) GetSummary(ctx context.Context, date string) (*domain.MonitoringSummary, error) {
	var summary domain.MonitoringSummary

	query := mr.db.QueryBuilder.Select("COUNT(DISTINCT id)").From("users")
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	err = mr.db.QueryRow(ctx, sql, args...).Scan(&summary.TotalUsers)
	if err != nil {
		return nil, err
	}

	query = mr.db.QueryBuilder.Select("COUNT(DISTINCT users.id)").
		From("users").
		LeftJoin("employees ON users.id = employees.user_id").
		Where("employees.status = 'active'")
	sql, args, err = query.ToSql()
	if err != nil {
		return nil, err
	}
	err = mr.db.QueryRow(ctx, sql, args...).Scan(&summary.ActiveUsers)
	if err != nil {
		return nil, err
	}

	query = mr.db.QueryBuilder.Select("COUNT(DISTINCT id)").From("attendances").Where("type = ?", "check_in")
	if date != "" {
		query = query.Where("DATE(time) = ?", date)
	}
	sql, args, err = query.ToSql()
	if err != nil {
		return nil, err
	}
	err = mr.db.QueryRow(ctx, sql, args...).Scan(&summary.TotalCheckin)
	if err != nil {
		return nil, err
	}

	query = mr.db.QueryBuilder.Select("COUNT(DISTINCT id)").From("attendances").Where("type = ?", "check_out")
	if date != "" {
		query = query.Where("DATE(time) = ?", date)
	}
	sql, args, err = query.ToSql()
	if err != nil {
		return nil, err
	}
	err = mr.db.QueryRow(ctx, sql, args...).Scan(&summary.TotalCheckOut)
	if err != nil {
		return nil, err
	}

	query = mr.db.QueryBuilder.Select("COUNT(DISTINCT id)").From("leave_requests").Where("status = 'pending'")
	sql, args, err = query.ToSql()
	if err != nil {
		return nil, err
	}
	err = mr.db.QueryRow(ctx, sql, args...).Scan(&summary.PendingLeaves)
	if err != nil {
		return nil, err
	}

	query = mr.db.QueryBuilder.Select("COUNT(DISTINCT id)").From("leave_requests").Where("status = 'approved'")
	sql, args, err = query.ToSql()
	if err != nil {
		return nil, err
	}
	err = mr.db.QueryRow(ctx, sql, args...).Scan(&summary.ApprovedLeaves)
	if err != nil {
		return nil, err
	}

	query = mr.db.QueryBuilder.Select("COUNT(DISTINCT id)").From("leave_requests").Where("status = 'rejected'")
	sql, args, err = query.ToSql()
	if err != nil {
		return nil, err
	}
	err = mr.db.QueryRow(ctx, sql, args...).Scan(&summary.RejectedLeaves)
	if err != nil {
		return nil, err
	}

	summary.GeneratedAt = time.Now()
	return &summary, nil
}

func (mr *MonitoringRepository) GetDashboardAnalytics(ctx context.Context) (*domain.DashboardAnalytics, error) {
	var analytics domain.DashboardAnalytics
	return &analytics, nil
}

func (mr *MonitoringRepository) GenerateAttendanceReport(ctx context.Context) (*domain.AttendanceReport, error) {
	var report domain.AttendanceReport
	return &report, nil
}

func (mr *MonitoringRepository) DetectAnomalies(ctx context.Context) ([]domain.Anomaly, error) {
	var anomalies []domain.Anomaly
	return anomalies, nil
}

func (mr *MonitoringRepository) ExportData(ctx context.Context, req domain.ExportRequest) (*domain.ExportResponse, error) {
	var response domain.ExportResponse
	return &response, nil
}

func (mr *MonitoringRepository) InsertReport(ctx context.Context, report *domain.MonitoringReport) error {
	query := mr.db.QueryBuilder.Insert("monitoring_reports").
		Columns("id, report_type", "data", "generated_at").
		Values(report.ID, report.ReportType, report.Data, report.GeneratedAt)

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = mr.db.Exec(ctx, sql, args...)
	return err
}

func (mr *MonitoringRepository) InsertReports(ctx context.Context, reports []domain.MonitoringReport) error {
	query := mr.db.QueryBuilder.Insert("monitoring_reports").
		Columns("report_type", "data", "generated_at")

	for _, report := range reports {
		query = query.Values(report.ReportType, report.Data, report.GeneratedAt)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = mr.db.Exec(ctx, sql, args...)
	return err
}

func (mr *MonitoringRepository) GetReportByDate(ctx context.Context, input domain.ExportRequest) ([]domain.MonitoringReport, error) {
	query := mr.db.QueryBuilder.Select("*").
		From("monitoring_reports").
		Where("date(generated_at) BETWEEN ? AND ?", input.StartDate, input.EndDate)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := mr.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []domain.MonitoringReport
	for rows.Next() {
		var report domain.MonitoringReport
		err := rows.Scan(
			&report.ID,
			&report.ReportType,
			&report.Data,
			&report.GeneratedAt,
			&report.CreatedAt,
			&report.UpdatedAt,
		)
		if err != nil {
			return reports, err
		}
		reports = append(reports, report)
	}

	return reports, nil
}
