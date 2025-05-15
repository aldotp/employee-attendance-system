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
	FullName        string     `json:"full_name"`
	Role            UserRole   `json:"role"`
	Location        string     `json:"location"`
	Timezone        string     `json:"timezone"`
	PhotoURL        string     `json:"photo_url"`
	Status          UserStatus `json:"status"`
	EmailVerifiedAt *time.Time `json:"email_verified_at"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at"`
}
