package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/CryptoCrowd/internal/db"
	"github.com/CryptoCrowd/internal/model"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

var (
	// ErrUserNotFound определяет ошибку, которая возникает, когда пользователь не найден в репозитории
	ErrUserNotFound = errors.New("пользователь не найден")
	// ErrUserAlreadyExists определяет ошибку, которая возникает, когда пользователь с такой почтой уже существует в репозитории
	ErrUserAlreadyExists = errors.New("пользователь с таким email уже существует")
	// ErrTransactionStartError определяет ошибку, которая возникает при ошибке начала транзакции
	ErrTransactionStartError = errors.New("ошибка начала транзакции")
)

type PostgresAccount struct {
	pool *db.Pool
}

func NewPostgresAccount(pool *db.Pool) *PostgresAccount {
	return &PostgresAccount{
		pool: pool,
	}
}

func (r *PostgresAccount) Create(ctx context.Context, acc model.Account, plainPassword string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrTransactionStartError, err)
	}
	defer tx.Rollback(ctx)

	var exists bool
	err = tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 and role = $2 FOR UPDATE)", acc.Email, acc.Role).Scan(&exists)
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

// UpdatePassword обновляет пароль пользователя
func (r *PostgresAccount) UpdatePassword(ctx context.Context, id int64, newPassword string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrTransactionStartError, err)
	}
	defer tx.Rollback(ctx)

	var exists bool
	err = tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1 FOR UPDATE)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("ошибка проверки существования пользователя: %w", err)
	}
	if !exists {
		return ErrUserNotFound
	}

	// Хешируем новый пароль
	passwordHash, err := hashPasswordSHA256(newPassword)
	if err != nil {
		return fmt.Errorf("ошибка хеширования пароля: %w", err)
	}

	now := time.Now()

	commandTag, err := tx.Exec(ctx, `
        UPDATE users
        SET password_hash = $2, updated_at = $3
        WHERE id = $1`,
		id,
		passwordHash,
		now,
	)
	if err != nil {
		return fmt.Errorf("ошибка обновления пароля: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return tx.Commit(ctx)
}

// Delete удаляет пользователя по ID
func (r *PostgresAccount) Delete(ctx context.Context, id int64) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrTransactionStartError, err)
	}
	defer tx.Rollback(ctx)

	var exists bool
	err = tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1 FOR UPDATE)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("ошибка проверки существования пользователя: %w", err)
	}
	if !exists {
		return ErrUserNotFound
	}

	commandTag, err := tx.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("ошибка удаления пользователя: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return tx.Commit(ctx)
}

// GetByEmailAndRole получает пользователя по email
func (r *PostgresAccount) GetByEmailAndRole(ctx context.Context, email string, role string) (model.Account, error) {
	var user model.Account
	err := pgxscan.Get(ctx, r.pool, &user,
		`SELECT id, username, email, role, created_at, updated_at FROM users WHERE email = $1 and role = $2`,
		email,
		role,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Account{}, ErrUserNotFound
		}
		return model.Account{}, fmt.Errorf("ошибка получения пользователя: %w", err)
	}

	return user, nil
}

// List возвращает список всех пользователей
func (r *PostgresAccount) List(ctx context.Context, searchTerm string) ([]model.Account, error) {
	var users []model.Account
	query := "SELECT id, username, email, role, created_at, updated_at FROM users WHERE 1=1"
	var args []any

	if searchTerm != "" {
		query += " AND (username LIKE $1 OR email LIKE $1)"
		args = append(args, "%"+searchTerm+"%")
	}

	query += " ORDER BY username"

	var err error
	if len(args) > 0 {
		err = pgxscan.Select(ctx, r.pool, &users, query, args...)
	} else {
		err = pgxscan.Select(ctx, r.pool, &users, query)
	}

	if err != nil {
		return nil, fmt.Errorf("ошибка получения списка пользователей: %w", err)
	}

	return users, nil
}

// TODO: Add check password function
