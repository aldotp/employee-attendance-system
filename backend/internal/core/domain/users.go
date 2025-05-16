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

type UserWithEmployee struct {
	ID              string     `json:"id"`
	Email           string     `json:"email"`
	Password        string     `json:"-"`
	Role            string     `json:"role"`
	EmailVerifiedAt *time.Time `json:"email_verified_at"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at"`
	EmployeeID      *string    `json:"employee_id"`
	Location        *string    `json:"location"`
	Timezone        *string    `json:"timezone"`
	PhotoURL        *string    `json:"photo_url"`
	Status          *string    `json:"status"`
	ReportingTo     *string    `json:"reporting_to"`
	Name            string     `json:"name"`
	DepartmentID    *string    `json:"department_id"`
	JoinedDate      *time.Time `json:"join_date"`
}
