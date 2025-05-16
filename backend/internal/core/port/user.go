package port

import (
	"context"

	"github.com/aldotp/employee-attendance-system/internal/adapter/dto"
	"github.com/aldotp/employee-attendance-system/internal/core/domain"
	"github.com/jackc/pgx/v5"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	ListUsers(ctx context.Context, skip, limit uint64) ([]domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	DeleteUser(ctx context.Context, id string) error
	ExistEmail(ctx context.Context, email string) (bool, error)
	CreateUserTx(ctx context.Context, tx pgx.Tx, user *domain.User) (*domain.User, error)
	BeginTx(ctx context.Context) (pgx.Tx, error)
	FindAll(ctx context.Context) ([]domain.User, error)
	FindByID(ctx context.Context, id string) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) (*domain.User, error)
	Delete(ctx context.Context, id string) error
	FindAllWithDetails(ctx context.Context) ([]domain.UserWithEmployee, error)
}

type UserService interface {
	Register(ctx context.Context, input *dto.RegisterRequest) (*domain.User, error)
	GetUser(ctx context.Context, id string) (*domain.User, error)
	ListUsers(ctx context.Context, skip, limit uint64) ([]domain.User, error)
	DeleteUser(ctx context.Context, id string) error
	UpdateProfile(ctx context.Context, input dto.UpdateUserRequest) (*domain.User, error)
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	CreateUser(ctx context.Context, req dto.CreateUserRequest) (*domain.User, error)
	UpdateUserByID(ctx context.Context, id string, req dto.UpdateUserRequest) (*domain.User, error)
	DeleteUserByID(ctx context.Context, id string) error
}
