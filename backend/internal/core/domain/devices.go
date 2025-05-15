package domain

import (
	"time"
)

type DeviceType string

const (
	Biometric DeviceType = "biometric"
	RFID      DeviceType = "rfid"
)

type Device struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Type      DeviceType `json:"type"`
	Location  string     `json:"location"`
	Status    string     `json:"status"`
	LastCheck time.Time  `json:"last_check"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
