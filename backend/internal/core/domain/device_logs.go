package domain

import (
	"time"
)

type DeviceLog struct {
	ID          string    `json:"id"`
	DeviceID    string    `json:"device_id"`
	UserID      string    `json:"user_id"`
	Action      string    `json:"action"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
