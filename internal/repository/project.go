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
	ErrProjectNotFound = errors.New("проект не найден")
	ErrProjectTxStart  = errors.New("ошибка начала транзакции")
)

type PostgresProject struct {
	pool *db.Pool
}

func NewPostgresProject(pool *db.Pool) *PostgresProject {
	return &PostgresProject{
		pool: pool,
	}
}

func (r *PostgresProject) Create(ctx context.Context, project model.Project) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrProjectTxStart, err)
	}
	defer tx.Rollback(ctx)

	now := time.Now()
	_, err = tx.Exec(ctx, `
        INSERT INTO projects (owner_id, status, name, description, amount_requested, amount_raised, deadline_at, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		project.OwnerID,
		project.Status,
		project.Name,
		project.Description,
		project.AmountRequested,
		project.AmountRaised,
		project.DeadlineAt,
		now,
	)
	if err != nil {
		return fmt.Errorf("ошибка создания проекта: %w", err)
	}

	return tx.Commit(ctx)
}

func (r *PostgresProject) Update(ctx context.Context, project model.Project) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrProjectTxStart, err)
	}
	defer tx.Rollback(ctx)

	var exists bool
	err = tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM projects WHERE id = $1 FOR UPDATE)", project.ID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("ошибка проверки существования проекта: %w", err)
	}
	if !exists {
		return ErrProjectNotFound
	}

	_, err = tx.Exec(ctx, `
        UPDATE projects
        SET status = $2, name = $3, description = $4, amount_requested = $5, deadline_at = $6
        WHERE id = $1`,
		project.ID,
		project.Status,
		project.Name,
		project.Description,
		project.AmountRequested,
		project.DeadlineAt,
	)
	if err != nil {
		return fmt.Errorf("ошибка обновления проекта: %w", err)
	}

	return tx.Commit(ctx)
}

func (r *PostgresProject) Delete(ctx context.Context, id int64) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrProjectTxStart, err)
	}
	defer tx.Rollback(ctx)

	var exists bool
	err = tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM projects WHERE id = $1 FOR UPDATE)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("ошибка проверки существования проекта: %w", err)
	}
	if !exists {
		return ErrProjectNotFound
	}

	_, err = tx.Exec(ctx, "DELETE FROM projects WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("ошибка удаления проекта: %w", err)
	}

	return tx.Commit(ctx)
}

func (r *PostgresProject) GetByID(ctx context.Context, id int64) (model.Project, error) {
	var project model.Project
	err := pgxscan.Get(ctx, r.pool, &project,
		`SELECT id, owner_id, status, name, description, amount_requested, amount_raised, deadline_at, created_at FROM projects WHERE id = $1`,
		id,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Project{}, ErrProjectNotFound
		}
		return model.Project{}, fmt.Errorf("ошибка получения проекта: %w", err)
	}
	return project, nil
}

// TODO: Implement GetByIDs method

func (r *PostgresProject) List(ctx context.Context, searchTerm string) ([]model.Project, error) {
	var projects []model.Project
	query := `SELECT id, owner_id, status, name, description, amount_requested, amount_raised, deadline_at, created_at FROM projects`
	var args []any

	if searchTerm != "" {
		query += " WHERE (name ILIKE $1 OR description ILIKE $1)"
		args = append(args, "%"+searchTerm+"%")
	}
	query += " ORDER BY created_at DESC"

	var err error
	if len(args) > 0 {
		err = pgxscan.Select(ctx, r.pool, &projects, query, args...)
	} else {
		err = pgxscan.Select(ctx, r.pool, &projects, query)
	}
	if err != nil {
		return nil, fmt.Errorf("ошибка получения списка проектов: %w", err)
	}
	return projects, nil
}

func (r *PostgresProject) ListByOwnerID(ctx context.Context, id int64, searchTerm string) ([]model.Project, error) {
	var projects []model.Project
	query := `SELECT id, owner_id, status, name, description, amount_requested, amount_raised, deadline_at, created_at FROM projects WHERE owner_id = $1`
	var args []any

	if searchTerm != "" {
		query += " AND (name ILIKE $2 OR description ILIKE $2)"
		args = append(args, "%"+searchTerm+"%")
	}
	query += " ORDER BY created_at DESC"

	var err error
	if len(args) > 0 {
		err = pgxscan.Select(ctx, r.pool, &projects, query, append([]any{id}, args...)...)
	} else {
		err = pgxscan.Select(ctx, r.pool, &projects, query, id)
	}
	if err != nil {
		return nil, fmt.Errorf("ошибка получения списка проектов: %w", err)
	}
	return projects, nil
}

func (r *PostgresProject) GetPhotosByProjectID(ctx context.Context, projectID int) ([]model.ProjectImage, error) {
	var photos []model.ProjectImage
	err := pgxscan.Select(ctx, r.pool, &photos,
		`SELECT id, project_id, url, created_at FROM project_images WHERE project_id = $1 ORDER BY created_at`, projectID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения фото проекта: %w", err)
	}
	return photos, nil
}

//TODO: Создать метод для создания фото проекта
