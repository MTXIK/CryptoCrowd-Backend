package model

import (
	"github.com/shopspring/decimal"
	"time"
)

type Investment struct {
	ID         int64           `json:"id"`
	UserID     int64           `json:"user_id"`
	ProjectID  int64           `json:"project_id"`
	Amount     decimal.Decimal `json:"amount"`
	InvestedAt *time.Time      `json:"invested_at"`
}
