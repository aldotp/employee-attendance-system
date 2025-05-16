package repository

import (
	"context"

	"github.com/aldotp/employee-attendance-system/internal/adapter/storage/postgres"
	"github.com/aldotp/employee-attendance-system/internal/core/domain"
	"github.com/jackc/pgx/v5"
)

type DepartmentRepository struct {
	db *postgres.DB
}

func NewDepartmentRepository(db *postgres.DB) *DepartmentRepository {
	return &DepartmentRepository{db: db}
}

func (r *DepartmentRepository) CreateDepartment(ctx context.Context, department *domain.Department) (*domain.Department, error) {
	query := r.db.QueryBuilder.Insert("departments").
		Columns("id", "name", "location", "timezone", "wfa_policy").
		Values(department.ID, department.Name, department.Location, department.Timezone, department.WFA_Policy)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	_, err = r.db.Exec(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return department, nil
}

func (r *DepartmentRepository) GetDepartmentByID(ctx context.Context, id string) (*domain.Department, error) {
	query := r.db.QueryBuilder.Select("id", "name", "location", "timezone", "wfa_policy").
		From("departments").
		Where("id = ?", id)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var dept domain.Department
	err = r.db.QueryRow(ctx, sql, args...).Scan(&dept.ID, &dept.Name, &dept.Location, &dept.Timezone, &dept.WFA_Policy)
	if err != nil {
		return nil, err
	}
	return &dept, nil
}

func (r *DepartmentRepository) ListDepartments(ctx context.Context, skip, limit uint64) ([]domain.Department, error) {
	var departments []domain.Department

	query := r.db.QueryBuilder.Select(
		"id",
		"COALESCE(name, '')",
		"COALESCE(location, '')",
		"COALESCE(timezone, '')",
		"COALESCE(wfa_policy, '{}'::json)",
	).From("departments")
	if skip > 0 {
		query = query.Offset(skip)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var dept domain.Department
		if err := rows.Scan(&dept.ID, &dept.Name, &dept.Location, &dept.Timezone, &dept.WFA_Policy); err != nil {
			return nil, err
		}
		departments = append(departments, dept)
	}

	return departments, nil
}

func (r *DepartmentRepository) UpdateDepartment(ctx context.Context, department *domain.Department) (*domain.Department, error) {
	query := r.db.QueryBuilder.Update("departments").
		Set("name", department.Name).
		Set("location", department.Location).
		Set("timezone", department.Timezone).
		Set("wfa_policy", department.WFA_Policy).
		Where("id = ?", department.ID)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	result, err := r.db.Exec(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	if result.RowsAffected() == 0 {
		return nil, pgx.ErrNoRows
	}
	return department, nil
}

func (r *DepartmentRepository) DeleteDepartment(ctx context.Context, id string) error {
	query := r.db.QueryBuilder.Delete("departments").Where("id = ?", id)
	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	result, err := r.db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
func (r *DepartmentRepository) GetDepartmentByName(ctx context.Context, name string) (*domain.Department, error) {
	query := r.db.QueryBuilder.Select(
		"id",
		"COALESCE(name, '')",
		"COALESCE(location, '')",
		"COALESCE(timezone, '')",
		"COALESCE(wfa_policy, '{}'::json)",
	).
		From("departments").
		Where("name = ?", name)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var dept domain.Department
	err = r.db.QueryRow(ctx, sql, args...).Scan(&dept.ID, &dept.Name, &dept.Location, &dept.Timezone, &dept.WFA_Policy)
	if err != nil {
		return nil, err
	}
	return &dept, nil
}
