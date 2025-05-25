package model

import "time"

type ProjectPhoto struct {
	ID        int       `db:"id" json:"id"`
	ProjectID int       `db:"project_id" json:"project_id"`
	URL       string    `db:"url" json:"url"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type Project struct {
	ID              int        `db:"id" json:"id,omitempty"`
	OwnerID         int        `db:"ownerid" json:"owner_id"`
	Status          string     `db:"status" json:"status"`
	Name            string     `db:"name" json:"name"`
	Description     string     `db:"description" json:"description"`
	AmountRequested float64    `db:"amountrequested" json:"amount_requested"`
	AmountRaised    float64    `db:"amountraised" json:"amount_raised"`
	DeadlineAt      *time.Time `db:"deadlineat" json:"deadline_at,omitempty"`
	CreatedAt       *time.Time `db:"createdat" json:"created_at,omitempty"`
}
