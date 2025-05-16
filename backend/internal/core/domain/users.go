package domain

import (
	"time"
)

type UserRole string

const (
	Admin     UserRole = "admin"
	HR        UserRole = "hr"
	Manager   UserRole = "manager"
	Employees UserRole = "employee"
)

var ExistRoleMap = map[UserRole]string{
	Admin:     "Admin",
	HR:        "HR",
	Manager:   "Manager",
	Employees: "Employee",
}

type UserStatus string

const (
	Active     UserStatus = "active"
	Inactive   UserStatus = "inactive"
	Terminated UserStatus = "terminated"
	Suspended  UserStatus = "suspended"
)

type User struct {
	ID              string     `json:"id"`
	Email           string     `json:"email"`
	Password        string     `json:"password"`
	Role            UserRole   `json:"role"`
	EmailVerifiedAt *time.Time `json:"email_verified_at"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at"`
}
