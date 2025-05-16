package repository

import (
	"context"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/aldotp/employee-attendance-system/internal/adapter/storage/postgres"
	"github.com/aldotp/employee-attendance-system/internal/core/domain"
	"github.com/aldotp/employee-attendance-system/pkg/consts"

	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	db *postgres.DB
}

func NewUserRepository(db *postgres.DB) *UserRepository {
	return &UserRepository{
		db,
	}
}

func (ur *UserRepository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	query := ur.db.QueryBuilder.Insert("users").
		Columns("id", "email", "password", "full_name", "role", "location", "timezone", "photo_url", "status", "email_verified_at", "created_at", "updated_at").
		Values(user.ID, user.Email, user.Password, user.Role, user.EmailVerifiedAt, user.CreatedAt, user.UpdatedAt).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = ur.db.QueryRow(ctx, sql, args...).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.EmailVerifiedAt,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *UserRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User

	query := ur.db.QueryBuilder.Select("*").
		From("users").
		Where(sq.Eq{"id": id}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = ur.db.QueryRow(ctx, sql, args...).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.EmailVerifiedAt,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, consts.ErrDataNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User

	query := ur.db.QueryBuilder.Select("*").
		From("users").
		Where(sq.Eq{"email": email}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = ur.db.QueryRow(ctx, sql, args...).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.EmailVerifiedAt,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, consts.ErrDataNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) ListUsers(ctx context.Context, offset, limit uint64) ([]domain.User, error) {
	var user domain.User
	var users []domain.User

	if limit == 0 {
		limit = 10
	}

	if offset == 0 {
		offset = 1
	}

	query := ur.db.QueryBuilder.Select("*").
		From("users").
		OrderBy("id").
		Limit(limit).
		Offset((offset - 1) * limit)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := ur.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Password,
			&user.Role,
			&user.EmailVerifiedAt,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (ur *UserRepository) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	query := ur.db.QueryBuilder.Update("users").
		Set("email", sq.Expr("COALESCE(?, email)", nullString(user.Email))).
		Set("password", sq.Expr("COALESCE(?, password)", nullString(user.Password))).
		Set("role", sq.Expr("COALESCE(?, role)", nullString(string(user.Role)))).
		Set("email_verified_at", sq.Expr("COALESCE(?, email_verified_at)", user.EmailVerifiedAt)).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": user.ID}).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = ur.db.QueryRow(ctx, sql, args...).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.EmailVerifiedAt,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser deletes a user by ID from the database
func (ur *UserRepository) DeleteUser(ctx context.Context, id string) error {
	query := ur.db.QueryBuilder.Delete("users").
		Where(sq.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = ur.db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) ExistEmail(ctx context.Context, email string) (bool, error) {
	query := ur.db.QueryBuilder.Select("id").
		From("users").
		Where(sq.Eq{"email": email}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return false, err
	}

	var id string

	err = ur.db.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (ur *UserRepository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	tx, err := ur.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (ur *UserRepository) CreateUserTx(ctx context.Context, tx pgx.Tx, user *domain.User) (*domain.User, error) {
	query := `INSERT INTO users (id, email, password, role, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6)
		RETURNING id, email, password, role, created_at, updated_at`
	row := tx.QueryRow(ctx,
		query,
		user.ID, user.Email, user.Password, user.Role, user.CreatedAt, user.UpdatedAt,
	)
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *UserRepository) FindAll(ctx context.Context) ([]domain.User, error) {
	var users []domain.User

	query := ur.db.QueryBuilder.Select("*").
		From("users")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := ur.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user domain.User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Password,
			&user.Role,
			&user.EmailVerifiedAt,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (ur *UserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User

	query := ur.db.QueryBuilder.Select("*").
		From("users").
		Where(sq.Eq{"id": id}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = ur.db.QueryRow(ctx, sql, args...).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.EmailVerifiedAt,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, consts.ErrDataNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	query := ur.db.QueryBuilder.Insert("users").
		Columns("id", "email", "password", "role", "email_verified_at", "created_at", "updated_at", "deleted_at").
		Values(user.ID, user.Email, user.Password, user.Role, user.EmailVerifiedAt, user.CreatedAt, user.UpdatedAt, user.DeletedAt).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = ur.db.QueryRow(ctx, sql, args...).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.EmailVerifiedAt,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return nil, consts.ErrConflictingData
		}
		return nil, err
	}

	return user, nil
}

func (ur *UserRepository) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	query := ur.db.QueryBuilder.Update("users").
		Set("email", user.Email).
		Set("password", user.Password).
		Set("role", user.Role).
		Set("email_verified_at", user.EmailVerifiedAt).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": user.ID}).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = ur.db.QueryRow(ctx, sql, args...).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.EmailVerifiedAt,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, consts.ErrDataNotFound
		}
		return nil, err
	}

	return user, nil
}

func (ur *UserRepository) Delete(ctx context.Context, id string) error {
	query := ur.db.QueryBuilder.Delete("users").
		Where(sq.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	result, err := ur.db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return consts.ErrDataNotFound
	}

	return nil
}

func (ur *UserRepository) FindAllWithDetails(ctx context.Context) ([]domain.UserWithEmployee, error) {
	var users []domain.UserWithEmployee

	query := ur.db.QueryBuilder.
		Select(
			"users.id",
			"users.email",
			"users.password",
			"users.role",
			"users.email_verified_at",
			"users.created_at",
			"users.updated_at",
			"users.deleted_at",
			"employees.id",
			"employees.location",
			"employees.timezone",
			"employees.photo_url",
			"employees.status",
			"employees.reporting_to",
			"employees.name",
			"employees.department_id",
			"employees.join_date",
		).
		From("users").
		LeftJoin("employees ON users.id = employees.user_id")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := ur.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user domain.UserWithEmployee
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Password,
			&user.Role,
			&user.EmailVerifiedAt,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
			&user.EmployeeID,
			&user.Location,
			&user.Timezone,
			&user.PhotoURL,
			&user.Status,
			&user.ReportingTo,
			&user.Name,
			&user.DepartmentID,
			&user.JoinedDate,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
