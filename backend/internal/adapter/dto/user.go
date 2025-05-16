package dto

import (
	"time"

	"github.com/aldotp/employee-attendance-system/internal/core/domain"
)

type UserResponse struct {
	ID        uint64    `json:"id" example:"1"`
	Name      string    `json:"name" example:"John Doe"`
	Email     string    `json:"email" example:"test@example.com"`
	CreatedAt time.Time `json:"created_at" example:"1970-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"1970-01-01T00:00:00Z"`
}

type RegisterRequest struct {
	Name       string `json:"name" binding:"required" example:"John Doe"`
	Email      string `json:"email" binding:"required,email" example:"test@example.com"`
	Password   string `json:"password" binding:"required,min=8" example:"12345678"`
	FullName   string `json:"full_name"`
	Department string `json:"department" binding:"required"`
	Location   string `json:"location"`
	Timezone   string `json:"timezone"`
	PhotoURL   string `json:"photo_url"`
}

type ListUserRequest struct {
	Page  uint64 `form:"page"`
	Limit uint64 `form:"limit"`
}

type GetUserRequest struct {
	ID uint64 `uri:"id" binding:"required,min=1" example:"1"`
}

type DeleteUserRequest struct {
	ID uint64 `uri:"id" binding:"required,min=1" example:"1"`
}

type ActivationRequest struct {
	Token string `uri:"token" binding:"required" form:"token"`
}

type CreateUserRequest struct {
	Name     string            `json:"name" validate:"required"`
	Email    string            `json:"email" validate:"required,email"`
	Password string            `json:"password" validate:"required,min=8"`
	Role     domain.UserRole   `json:"role" validate:"required,user_role"`
	Location string            `json:"location" validate:"required"`
	Timezone string            `json:"timezone" validate:"required"`
	PhotoURL string            `json:"photo_url"`
	Status   domain.UserStatus `json:"status" validate:"required"`
}

type UpdateUserRequest struct {
	Name     string            `json:"name"`
	Email    string            `json:"email"`
	Password string            `json:"password"`
	Role     domain.UserRole   `json:"role"`
	Location string            `json:"location"`
	Timezone string            `json:"timezone"`
	PhotoURL string            `json:"photo_url"`
	Status   domain.UserStatus `json:"status"`
}
