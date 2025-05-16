package service

import (
	"context"
	"fmt"
	"time"

	"github.com/aldotp/employee-attendance-system/internal/adapter/dto"
	"github.com/aldotp/employee-attendance-system/internal/core/domain"
	"github.com/aldotp/employee-attendance-system/internal/core/port"
	"github.com/aldotp/employee-attendance-system/pkg/consts"
	"github.com/aldotp/employee-attendance-system/pkg/util"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserService struct {
	cache          port.CacheInterface
	repo           port.UserRepository
	employeeRepo   port.EmployeeRepository
	departmentRepo port.DepartmentRepository
	token          port.TokenInterface
	log            *zap.Logger
}

func NewUserService(repo port.UserRepository, employeeRepo port.EmployeeRepository, departmentRepo port.DepartmentRepository, cache port.CacheInterface, token port.TokenInterface, log *zap.Logger) *UserService {
	return &UserService{
		cache:          cache,
		repo:           repo,
		employeeRepo:   employeeRepo,
		departmentRepo: departmentRepo,
		token:          token,
		log:            log,
	}
}

// Register creates a new user
func (us *UserService) Register(ctx context.Context, input *dto.RegisterRequest) (*domain.User, error) {
	exist, err := us.repo.ExistEmail(ctx, input.Email)
	if err != nil {
		us.log.Error(err.Error())
		return nil, consts.ErrInternal
	}

	if exist {
		return nil, consts.ErrEmailAlreadyExist
	}

	hashedPassword, err := util.HashPassword(input.Password)
	if err != nil {
		us.log.Error(err.Error())
		return nil, consts.ErrInternal
	}

	tNow := time.Now()

	tx, err := us.repo.BeginTx(ctx)
	if err != nil {
		us.log.Error(err.Error())
		return nil, consts.ErrInternal
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback(ctx)
			panic(p)
		} else if err != nil {
			tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	user := &domain.User{
		ID:              uuid.New().String(),
		Email:           input.Email,
		Password:        hashedPassword,
		EmailVerifiedAt: &tNow,
		Role:            domain.Employees,
		CreatedAt:       tNow,
		UpdatedAt:       tNow,
	}

	user, err = us.repo.CreateUserTx(ctx, tx, user)
	if err != nil {
		us.log.Error(err.Error())
		if err == consts.ErrConflictingData {
			return nil, err
		}
		return nil, consts.ErrInternal
	}

	if input.Department == "" {
		input.Department = "Information Technology"
	}

	department, err := us.departmentRepo.GetDepartmentByName(ctx, input.Department)
	if err != nil {
		us.log.Error(err.Error())
		return nil, consts.ErrInternal
	}

	employee := &domain.Employee{
		ID:           uuid.New().String(),
		UserID:       user.ID,
		Name:         input.FullName,
		PhotoURL:     input.PhotoURL,
		Location:     input.Location,
		Timezone:     "UTC",
		Status:       domain.StatusActive,
		DepartmentID: department.ID,
		JoinDate:     tNow,
		CreatedAt:    tNow,
		UpdatedAt:    tNow,
	}

	_, err = us.employeeRepo.CreateEmployeeTx(ctx, tx, employee)
	if err != nil {
		us.log.Error(err.Error())
		if err == consts.ErrConflictingData {
			return nil, err
		}
		return nil, consts.ErrInternal
	}

	return user, nil
}

// GetUser gets a user by ID
func (us *UserService) GetUser(ctx context.Context, id string) (*domain.User, error) {
	var user *domain.User

	cacheKey := util.GenerateCacheKey("user", id)
	cachedUser, err := us.cache.Get(ctx, cacheKey)
	if err == nil {
		err = util.Deserialize(cachedUser, &user)
		if err != nil {
			return nil, consts.ErrInternal
		}
		return user, nil
	}

	user, err = us.repo.GetUserByID(ctx, id)
	if err != nil {
		if err == consts.ErrDataNotFound {
			return nil, err
		}
		return nil, consts.ErrInternal
	}

	userSerialized, err := util.Serialize(user)
	if err != nil {
		return nil, consts.ErrInternal
	}

	err = us.cache.Set(ctx, cacheKey, userSerialized, 0)
	if err != nil {
		return nil, consts.ErrInternal
	}

	return user, nil
}

func (us *UserService) ListUsers(ctx context.Context, page, limit uint64) ([]domain.User, error) {
	var users []domain.User

	params := util.GenerateCacheKeyParams(page, limit)
	cacheKey := util.GenerateCacheKey("users", params)

	cachedUsers, err := us.cache.Get(ctx, cacheKey)
	if err == nil {
		err = util.Deserialize(cachedUsers, &users)
		if err != nil {
			return nil, consts.ErrInternal
		}

		return users, nil
	}

	users, err = us.repo.ListUsers(ctx, page, limit)
	if err != nil {
		return nil, consts.ErrInternal
	}

	usersSerialized, err := util.Serialize(users)
	if err != nil {
		return nil, consts.ErrInternal
	}

	err = us.cache.Set(ctx, cacheKey, usersSerialized, time.Minute*10)
	if err != nil {
		return nil, consts.ErrInternal
	}

	return users, nil
}

// DeleteUser deletes a user by ID
func (us *UserService) DeleteUser(ctx context.Context, id string) error {
	_, err := us.repo.GetUserByID(ctx, id)
	if err != nil {
		if err == consts.ErrDataNotFound {
			return err
		}
		return consts.ErrInternal
	}

	cacheKey := util.GenerateCacheKey("user", id)

	err = us.cache.Delete(ctx, cacheKey)
	if err != nil {
		return consts.ErrInternal
	}

	err = us.cache.DeleteByPrefix(ctx, "users:*")
	if err != nil {
		return consts.ErrInternal
	}

	return us.repo.DeleteUser(ctx, id)
}

func (s *UserService) UpdateProfile(ctx context.Context, input dto.UpdateUserRequest) (*domain.User, error) {
	userSession := util.GetAuthPayload(ctx, consts.AuthorizationKey)
	if userSession == nil {
		return nil, fmt.Errorf("user session not found")
	}

	var user domain.User
	user.ID = userSession.UserID

	if input.Password != "" {
		hashedPassword, err := util.HashPassword(input.Password)
		if err != nil {
			return nil, consts.ErrInternal
		}

		user.Password = hashedPassword
	}

	return s.repo.UpdateUser(ctx, &user)
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *UserService) CreateUser(ctx context.Context, req dto.CreateUserRequest) (*domain.User, error) {
	exist, err := s.repo.ExistEmail(ctx, req.Email)
	if err != nil {
		if err == consts.ErrDataNotFound {
			return nil, consts.ErrDataNotFound
		}
		s.log.Error("Error checking email existence", zap.Error(err))
		return nil, consts.ErrInternal
	}

	if exist {
		s.log.Warn("Email already exists", zap.String("email", req.Email))
		return nil, consts.ErrEmailAlreadyExist
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		s.log.Error("Error hashing password", zap.Error(err))
		return nil, consts.ErrInternal
	}

	user := &domain.User{
		Email:    req.Email,
		Password: hashedPassword,
		Role:     req.Role,
	}

	user, err = s.repo.Create(ctx, user)
	if err != nil {
		s.log.Error("Error creating user", zap.Error(err))
		return nil, consts.ErrInternal
	}

	return user, nil
}

func (s *UserService) UpdateUserByID(ctx context.Context, id string, req dto.UpdateUserRequest) (*domain.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Email != "" {
		exist, err := s.repo.ExistEmail(ctx, req.Email)
		if err != nil {
			return nil, err
		}
		if exist {
			return nil, consts.ErrEmailAlreadyExist
		}

		user.Email = req.Email
	}

	if req.Password != "" {
		hashedPassword, err := util.HashPassword(req.Password)
		if err != nil {
			return nil, consts.ErrInternal
		}
		user.Password = hashedPassword
	}

	if req.Role != "" {
		user.Role = req.Role
	}

	if req.Password != "" {
		user.Password = req.Password
	}
	user.Role = req.Role

	return s.repo.Update(ctx, user)
}

func (s *UserService) DeleteUserByID(ctx context.Context, id string) error {
	// Delete user from database
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	// Delete individual user cache
	userCacheKey := util.GenerateCacheKey("user", id)
	err = s.cache.Delete(ctx, userCacheKey)
	if err != nil {
		return consts.ErrInternal
	}

	// Delete all users list cache
	err = s.cache.DeleteByPrefix(ctx, "users:*")
	if err != nil {
		return consts.ErrInternal
	}

	return nil
}
