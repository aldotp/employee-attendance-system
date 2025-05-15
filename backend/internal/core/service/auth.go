package service

import (
	"context"

	"github.com/aldotp/employee-attendance-system/internal/adapter/dto"
	"github.com/aldotp/employee-attendance-system/internal/core/port"
	"github.com/aldotp/employee-attendance-system/pkg/consts"
	"github.com/aldotp/employee-attendance-system/pkg/util"
	"go.uber.org/zap"
)

type AuthService struct {
	repo         port.UserRepository
	employeeRepo port.EmployeeRepository
	ts           port.TokenInterface
	log          *zap.Logger
}

func NewAuthService(repo port.UserRepository, employeeRepo port.EmployeeRepository, ts port.TokenInterface, log *zap.Logger) *AuthService {
	return &AuthService{
		repo,
		employeeRepo,
		ts,
		log,
	}
}

func (as *AuthService) Login(ctx context.Context, email, password string) (dto.LoginResponse, error) {
	user, err := as.repo.GetUserByEmail(ctx, email)
	if err != nil {
		as.log.Error("failed to get user by email", zap.Error(err))
		if err == consts.ErrDataNotFound {
			return dto.LoginResponse{}, consts.ErrInvalidCredentials
		}
		return dto.LoginResponse{}, consts.ErrInternal
	}

	if err = util.ComparePassword(password, user.Password); err != nil {
		as.log.Error("password comparison failed", zap.Error(err))
		return dto.LoginResponse{}, consts.ErrInvalidCredentials
	}

	employee, err := as.employeeRepo.FindOneByFilters(ctx, map[string]interface{}{"user_id": user.ID})
	if err != nil {
		as.log.Error("failed to get employee by ID", zap.Error(err))
		if err == consts.ErrDataNotFound {
			return dto.LoginResponse{}, consts.ErrInvalidCredentials
		}
		return dto.LoginResponse{}, consts.ErrInternal
	}

	accessToken, err := as.ts.GenerateAccessToken(user, employee)
	if err != nil {
		as.log.Error("failed to generate access token", zap.Error(err))
		return dto.LoginResponse{}, consts.ErrTokenCreation
	}

	refreshToken, err := as.ts.GenerateRefreshToken(user, employee)
	if err != nil {
		as.log.Error("failed to generate refresh token", zap.Error(err))
		return dto.LoginResponse{}, consts.ErrTokenCreation
	}

	return dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (as *AuthService) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	payload, _, err := as.ts.VerifyRefreshToken(ctx, refreshToken)
	if err != nil {
		as.log.Error("failed to verify refresh token", zap.Error(err))
		return "", consts.ErrInvalidSignature
	}

	user, err := as.repo.GetUserByID(ctx, payload.UserID)
	if err != nil {
		as.log.Error("failed to get user by ID", zap.Error(err))
		if err == consts.ErrDataNotFound {
			return "", consts.ErrInvalidToken
		}
		return "", consts.ErrInternal
	}

	employee, err := as.employeeRepo.FindOneByFilters(ctx, map[string]interface{}{"user_id": user.ID})
	if err != nil {
		as.log.Error("failed to get employee by ID", zap.Error(err))
		if err == consts.ErrDataNotFound {
			return "", consts.ErrInvalidToken
		}
		return "", consts.ErrInternal
	}

	accessToken, err := as.ts.GenerateAccessToken(user, employee)
	if err != nil {
		as.log.Error("failed to generate access token", zap.Error(err))
		return "", consts.ErrTokenCreation
	}

	return accessToken, nil
}
