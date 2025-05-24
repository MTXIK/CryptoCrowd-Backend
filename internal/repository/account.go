package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/CryptoCrowd/internal/db"
	"github.com/CryptoCrowd/internal/model"
	"time"
)

var (
	// ErrUserNotFound определяет ошибку, которая возникает, когда пользователь не найден в репозитории
	ErrUserNotFound = errors.New("пользователь не найден")
	// ErrUserAlreadyExists определяет ошибку, которая возникает, когда пользователь с такой почтой уже существует в репозитории
	ErrUserAlreadyExists = errors.New("пользователь с таким email уже существует")
	// ErrTransactionStartError определяет ошибку, которая возникает при ошибке начала транзакции
	ErrTransactionStartError = errors.New("ошибка начала транзакции")
)

type PostgresAccountRepository struct {
	pool *db.Pool
}

func NewPostgresAccountRepository(pool *db.Pool) *PostgresAccountRepository {
	return &PostgresAccountRepository{
		pool: pool,
	}
}

func (r *PostgresAccountRepository) Create(ctx context.Context, acc model.Account, plainPassword string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrTransactionStartError, err)
	}
	defer tx.Rollback(ctx)

	var exists bool
	err = tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 FOR UPDATE)", acc.Email).Scan(&exists)
	if err != nil {
		return fmt.Errorf("ошибка проверки существования пользователя: %w", err)
	}
	if exists {
		return ErrUserAlreadyExists
	}

	passwordHash, err := hashPasswordSHA256(plainPassword)
	if err != nil {
		return fmt.Errorf("ошибка хеширования пароля: %w", err)
	}

	now := time.Now()

	_, err = tx.Exec(ctx, `
        INSERT INTO users (username, email, password_hash, role, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)`,
		acc.Username,
		acc.Email,
		passwordHash,
		acc.Role,
		now,
		now,
	)
	if err != nil {
		return fmt.Errorf("ошибка создания пользователя: %w", err)
	}

	return tx.Commit(ctx)
}
