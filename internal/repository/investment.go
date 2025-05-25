package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/CryptoCrowd/internal/db"
	"github.com/CryptoCrowd/internal/model"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"time"
)

var (
	ErrInvestmentNotFound = errors.New("инвестиция не найдена")
	ErrInvestmentTxStart  = errors.New("ошибка начала транзакции")
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
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInvestmentTxStart, err)
	}
	defer tx.Rollback(ctx)

	now := time.Now()

	_, err = tx.Exec(ctx, `
        INSERT INTO investments (user_id, project_id, amount, invested_at)
        VALUES ($1, $2, $3, $4)`,
		investment.UserID,
		investment.ProjectID,
		investment.Amount,
		now,
	)
	if err != nil {
		return fmt.Errorf("ошибка создания инвестиции: %w", err)
	}

	return tx.Commit(ctx)
}

func (r *PostgresInvestment) GetByID(ctx context.Context, id int64) (model.Investment, error) {
	var investment model.Investment
	err := pgxscan.Get(ctx, r.pool, &investment,
		`SELECT id, user_id, project_id, amount, invested_at FROM investments WHERE id = $1`,
		id,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Investment{}, ErrInvestmentNotFound
		}
		return model.Investment{}, fmt.Errorf("ошибка получения инвестиции: %w", err)
	}
	return investment, nil
}

func (r *PostgresInvestment) GetByUserID(ctx context.Context, userID int64) ([]model.Investment, error) {
	var investments []model.Investment
	err := pgxscan.Select(ctx, r.pool, &investments,
		`SELECT id, user_id, project_id, amount, invested_at FROM investments WHERE user_id = $1 ORDER BY invested_at DESC`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения инвестиций пользователя: %w", err)
	}
	return investments, nil
}

func (r *PostgresInvestment) GetByProjectID(ctx context.Context, projectID int64) ([]model.Investment, error) {
	var investments []model.Investment
	err := pgxscan.Select(ctx, r.pool, &investments,
		`SELECT id, user_id, project_id, amount, invested_at FROM investments WHERE project_id = $1 ORDER BY invested_at DESC`,
		projectID,
	)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения инвестиций проекта: %w", err)
	}
	return investments, nil
}
