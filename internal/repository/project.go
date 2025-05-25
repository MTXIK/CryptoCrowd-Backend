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
	ErrProjectNotFound      = errors.New("проект не найден")
	ErrProjectAlreadyExists = errors.New("проект с таким именем уже существует")
	ErrProjectTxStart       = errors.New("ошибка начала транзакции")
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
        INSERT INTO projects (ownerID, status, name, description, amountRequested, amountRaised, deadlineAt, createdAt)
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

	var existing model.Project
	err = pgxscan.Get(ctx, tx, &existing, "SELECT id FROM projects WHERE id = $1 FOR UPDATE", project.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrProjectNotFound
		}
		return fmt.Errorf("ошибка получения проекта для обновления: %w", err)
	}

	_, err = tx.Exec(ctx, `
        UPDATE projects
        SET status = $2, name = $3, description = $4, amountRequested = $5, amountRaised = $6, deadlineAt = $7
        WHERE id = $1`,
		project.ID,
		project.Status,
		project.Name,
		project.Description,
		project.AmountRequested,
		project.AmountRaised,
		project.DeadlineAt,
	)
	if err != nil {
		return fmt.Errorf("ошибка обновления проекта: %w", err)
	}

	return tx.Commit(ctx)
}

func (r *PostgresProject) Delete(ctx context.Context, id int) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrProjectTxStart, err)
	}
	defer tx.Rollback(ctx)

	var existing model.Project
	err = pgxscan.Get(ctx, tx, &existing, "SELECT id FROM projects WHERE id = $1 FOR UPDATE", id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrProjectNotFound
		}
		return fmt.Errorf("ошибка получения проекта для удаления: %w", err)
	}

	_, err = tx.Exec(ctx, "DELETE FROM projects WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("ошибка удаления проекта: %w", err)
	}

	return tx.Commit(ctx)
}

func (r *PostgresProject) GetByID(ctx context.Context, id int) (model.Project, error) {
	var project model.Project
	err := pgxscan.Get(ctx, r.pool, &project,
		`SELECT id, ownerID, status, name, description, amountRequested, amountRaised, deadlineAt, createdAt FROM projects WHERE id = $1`,
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
	query := `SELECT id, ownerID, status, name, description, amountRequested, amountRaised, deadlineAt, createdAt FROM projects WHERE 1=1`
	var args []any

	if searchTerm != "" {
		query += " AND (name ILIKE $1 OR description ILIKE $1)"
		args = append(args, "%"+searchTerm+"%")
	}
	query += " ORDER BY createdAt DESC"

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
	query := `SELECT id, ownerID, status, name, description, amountRequested, amountRaised, deadlineAt, createdAt FROM projects WHERE ownerID = $1`
	var args []any

	if searchTerm != "" {
		query += " AND (name ILIKE $2 OR description ILIKE $2)"
		args = append(args, "%"+searchTerm+"%")
	}
	query += " ORDER BY createdAt DESC"

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

func (r *PostgresProject) GetPhotosByProjectID(ctx context.Context, projectID int) ([]model.ProjectPhoto, error) {
	var photos []model.ProjectPhoto
	err := pgxscan.Select(ctx, r.pool, &photos,
		`SELECT id, project_id, url, created_at FROM project_photos WHERE project_id = $1 ORDER BY created_at`, projectID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения фото проекта: %w", err)
	}
	return photos, nil
}

//TODO: Создать метод для создания фото проекта
//TODO: Создать метод для обновления собранной суммы проекта
