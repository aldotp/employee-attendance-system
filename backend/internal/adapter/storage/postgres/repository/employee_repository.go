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

type EmployeeRepository struct {
	db *postgres.DB
}

func NewEmployeeRepository(db *postgres.DB) *EmployeeRepository {
	return &EmployeeRepository{
		db,
	}
}

func (er *EmployeeRepository) CreateEmployee(ctx context.Context, employee *domain.Employee) (*domain.Employee, error) {
	query := er.db.QueryBuilder.Insert("employees").
		Columns("id", "user_id", "department_id", "name", "location", "timezone", "photo_url", "status", "join_date", "reporting_to", "created_at", "updated_at").
		Values(employee.ID, employee.UserID, employee.DepartmentID, employee.Name, employee.Location, employee.Timezone, employee.PhotoURL, employee.Status, employee.JoinDate, employee.ReportingTo, employee.CreatedAt, employee.UpdatedAt).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = er.db.QueryRow(ctx, sql, args...).Scan(
		&employee.ID,
		&employee.UserID,
		&employee.DepartmentID,
		&employee.Name,
		&employee.Location,
		&employee.Timezone,
		&employee.PhotoURL,
		&employee.Status,
		&employee.JoinDate,
		&employee.ReportingTo,
		&employee.CreatedAt,
		&employee.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return employee, nil
}

func (er *EmployeeRepository) GetEmployeeByID(ctx context.Context, id string) (*domain.Employee, error) {
	var employee domain.Employee

	query := er.db.QueryBuilder.Select("*").
		From("employees").
		Where(sq.Eq{"id": id}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = er.db.QueryRow(ctx, sql, args...).Scan(
		&employee.ID,
		&employee.UserID,
		&employee.DepartmentID,
		&employee.Name,
		&employee.Location,
		&employee.Timezone,
		&employee.PhotoURL,
		&employee.Status,
		&employee.JoinDate,
		&employee.ReportingTo,
		&employee.CreatedAt,
		&employee.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, consts.ErrDataNotFound
		}
		return nil, err
	}

	return &employee, nil
}

func (er *EmployeeRepository) ListEmployees(ctx context.Context, skip, limit uint64) ([]domain.Employee, error) {
	var employee domain.Employee
	var employees []domain.Employee

	if limit == 0 {
		limit = 10
	}

	if skip == 0 {
		skip = 1
	}

	query := er.db.QueryBuilder.Select("*").
		From("employees").
		OrderBy("id").
		Limit(limit).
		Offset((skip - 1) * limit)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := er.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&employee.ID,
			&employee.UserID,
			&employee.DepartmentID,
			&employee.Name,
			&employee.Location,
			&employee.Timezone,
			&employee.PhotoURL,
			&employee.Status,
			&employee.JoinDate,
			&employee.ReportingTo,
			&employee.CreatedAt,
			&employee.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		employees = append(employees, employee)
	}

	return employees, nil
}

func (er *EmployeeRepository) UpdateEmployee(ctx context.Context, employee *domain.Employee) (*domain.Employee, error) {
	query := er.db.QueryBuilder.Update("employees").
		Set("user_id", sq.Expr("COALESCE(?, user_id)", nullString(employee.UserID))).
		Set("department_id", sq.Expr("COALESCE(?, department_id)", nullString(employee.DepartmentID))).
		Set("name", sq.Expr("COALESCE(?, name)", nullString(employee.Name))).
		Set("location", sq.Expr("COALESCE(?, location)", nullString(employee.Location))).
		Set("timezone", sq.Expr("COALESCE(?, timezone)", nullString(employee.Timezone))).
		Set("photo_url", sq.Expr("COALESCE(?, photo_url)", nullString(employee.PhotoURL))).
		Set("status", sq.Expr("COALESCE(?, status)", nullString(string(employee.Status)))).
		Set("join_date", sq.Expr("COALESCE(?, join_date)", employee.JoinDate)).
		Set("reporting_to", sq.Expr("COALESCE(?, reporting_to)", nullString(employee.ReportingTo))).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": employee.ID}).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = er.db.QueryRow(ctx, sql, args...).Scan(
		&employee.ID,
		&employee.UserID,
		&employee.DepartmentID,
		&employee.Name,
		&employee.Location,
		&employee.Timezone,
		&employee.PhotoURL,
		&employee.Status,
		&employee.JoinDate,
		&employee.ReportingTo,
		&employee.CreatedAt,
		&employee.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return employee, nil
}

func (er *EmployeeRepository) DeleteEmployee(ctx context.Context, id string) error {
	query := er.db.QueryBuilder.Delete("employees").
		Where(sq.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = er.db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func (er *EmployeeRepository) FindOneByFilters(ctx context.Context, filter map[string]interface{}) (*domain.Employee, error) {
	query := er.db.QueryBuilder.Select(
		"id", "user_id", "department_id", "name", "email", "role", "location", "timezone", "photo_url", "status", "join_date", "reporting_to", "created_at", "updated_at",
	).From("employees")

	for key, value := range filter {
		query = query.Where(sq.Eq{key: value})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var employee domain.Employee
	err = er.db.QueryRow(ctx, sql, args...).Scan(
		&employee.ID,
		&employee.UserID,
		&nullStringWrapper{&employee.DepartmentID},
		&employee.Name,
		&employee.Location,
		&employee.Timezone,
		&employee.PhotoURL,
		&employee.Status,
		&employee.JoinDate,
		&nullStringWrapper{&employee.ReportingTo},
		&employee.CreatedAt,
		&employee.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &employee, nil
}

func (er *EmployeeRepository) CreateEmployeeTx(ctx context.Context, tx pgx.Tx, employee *domain.Employee) (*domain.Employee, error) {

	query := `INSERT INTO employees (id, user_id, department_id, name,  location, timezone, photo_url, status, join_date, reporting_to, created_at, updated_at) 
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) 
        RETURNING id, user_id, department_id, name, location, timezone, photo_url, status, join_date, reporting_to, created_at, updated_at`

	row := tx.QueryRow(ctx, query,
		employee.ID, employee.UserID, nullString(employee.DepartmentID), employee.Name, employee.Location, employee.Timezone, employee.PhotoURL, employee.Status, employee.JoinDate, nullString(employee.ReportingTo), employee.CreatedAt, employee.UpdatedAt,
	)
	err := row.Scan(
		&employee.ID, &employee.UserID, &nullStringWrapper{&employee.DepartmentID}, &employee.Name, &employee.Location, &employee.Timezone, &employee.PhotoURL, &employee.Status, &employee.JoinDate, &nullStringWrapper{&employee.ReportingTo}, &employee.CreatedAt, &employee.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return employee, nil
}

type nullStringWrapper struct {
	str *string
}

func (w *nullStringWrapper) Scan(value interface{}) error {
	if value == nil {
		*w.str = ""
		return nil
	}
	s, ok := value.(string)
	if !ok {
		return fmt.Errorf("expected string, got %T", value)
	}
	*w.str = s
	return nil
}
