package domain

type TokenPayload struct {
	Email      string   `json:"email"`
	UserID     string   `json:"user_id"`
	Role       UserRole `json:"role"`
	EmployeeID string   `json:"employee_id"`
}

type RefreshTokenPayload struct {
	UserID string `json:"user_id"`
}
