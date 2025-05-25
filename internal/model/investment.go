package model

import "time"

type Investment struct {
	ID         int64      `json:"id"`
	UserID     int64      `json:"user_id"`
	ProjectID  int64      `json:"project_id"`
	Amount     float64    `json:"amount"` // TODO: Не использовать float64, использовать github.com/Rhymond/go-money
	InvestedAt *time.Time `json:"invested_at"`
}
