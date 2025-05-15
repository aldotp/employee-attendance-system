package port

import (
	"context"

	"github.com/aldotp/employee-attendance-system/internal/adapter/dto"
	"github.com/aldotp/employee-attendance-system/internal/core/domain"
)

//go:generate mockgen -source=auth.go -destination=mock/auth.go -package=mock

type TokenInterface interface {
	GenerateAccessToken(user *domain.User, employee *domain.Employee) (string, error)
	GenerateRefreshToken(user *domain.User, employee *domain.Employee) (string, error)
	VerifyAccessToken(encodedToken string) (*domain.TokenPayload, error)
	VerifyRefreshToken(ctx context.Context, refreshToken string) (*domain.RefreshTokenPayload, bool, error)
	GenerateActivationToken(userID uint64) (string, error)
}

type AuthService interface {
	Login(ctx context.Context, email, password string) (dto.LoginResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, error)
}
