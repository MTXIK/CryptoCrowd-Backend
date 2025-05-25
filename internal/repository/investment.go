package repository

import (
	"context"
	"github.com/CryptoCrowd/internal/db"
	"github.com/CryptoCrowd/internal/model"
)

type PostgresInvestment struct {
	pool *db.Pool
}

func NewPostgresInvestment(pool *db.Pool) *PostgresInvestment {
	return &PostgresInvestment{
		pool: pool,
	}
}

func (r *PostgresInvestment) Create(ctx context.Context, investment model.Investment) error {
	return nil
}

func (r *PostgresInvestment) Update(ctx context.Context, investment model.Investment) error {
	return nil
}

func (r *PostgresInvestment) Delete(ctx context.Context, id int64) error {
	return nil
}

func (r *PostgresInvestment) GetByID(ctx context.Context, id int64) (model.Investment, error) {
	return model.Investment{}, nil
}

func (r *PostgresInvestment) GetByUserID(ctx context.Context, userID int64) ([]model.Investment, error) {
	return nil, nil
}

func (r *PostgresInvestment) GetByProjectID(ctx context.Context, projectID int64) ([]model.Investment, error) {
	return nil, nil
}
