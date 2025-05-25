package model

import (
	"github.com/shopspring/decimal"
	"time"
)

type ProjectImage struct {
	ID        int        `db:"id" json:"id"`
	ProjectID int        `db:"project_id" json:"project_id"`
	URL       string     `db:"url" json:"url"`
	CreatedAt *time.Time `db:"created_at" json:"created_at"`
}

type Project struct {
	ID              int64           `db:"id" json:"id,omitempty"`
	OwnerID         int64           `db:"owner_id" json:"owner_id"`
	Status          string          `db:"status" json:"status"`
	Name            string          `db:"name" json:"name"`
	Description     string          `db:"description" json:"description"`
	AmountRequested decimal.Decimal `db:"amount_requested" json:"amount_requested"`
	AmountRaised    decimal.Decimal `db:"amount_raised" json:"amount_raised"`
	DeadlineAt      *time.Time      `db:"deadline_at" json:"deadline_at"`
	CreatedAt       *time.Time      `db:"created_at" json:"created_at,omitempty"`
}
