package domain

import (
	"time"
)

type Department struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Location   string    `json:"location"`
	Timezone   string    `json:"timezone"`
	WFA_Policy string    `json:"wfa_policy"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
