package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"time"

	"github.com/CryptoCrowd/internal/logger"
	"github.com/CryptoCrowd/internal/model"
)

var (
	ErrInvalidProjectName        = errors.New("invalid project name")
	ErrInvalidProjectDescription = errors.New("invalid project description")
	ErrInvalidProjectOwner       = errors.New("invalid project owner")
	ErrInvalidProjectStatus      = errors.New("invalid project status")
	ErrInvalidProjectAmount      = errors.New("invalid project amount")
	ErrInvalidProjectDeadline    = errors.New("invalid project deadline")
)

// ProjectRepository defines the interface for project repository operations
type ProjectRepository interface {
	Create(ctx context.Context, project model.Project) error
	Update(ctx context.Context, project model.Project) error
	Delete(ctx context.Context, id int) error
	GetByID(ctx context.Context, id int) (model.Project, error)
	List(ctx context.Context, searchTerm string) ([]model.Project, error)
	ListByOwnerID(ctx context.Context, id int64, searchTerm string) ([]model.Project, error)
	GetPhotosByProjectID(ctx context.Context, projectID int) ([]model.ProjectImage, error)
}

// Project service implements business logic for project operations
type Project struct {
	repo ProjectRepository
}

// NewProject creates a new project service
func NewProject(repo ProjectRepository) *Project {
	return &Project{
		repo: repo,
	}
}

// validateProject validates project data
func (p *Project) validateProject(project model.Project) error {
	// Validate project name
	if project.Name == "" {
		logger.Error("Invalid project name")
		return fmt.Errorf("%w", ErrInvalidProjectName)
	}

	// Validate project description
	if project.Description == "" {
		logger.Error("Invalid project description")
		return fmt.Errorf("%w", ErrInvalidProjectDescription)
	}

	// Validate project owner
	if project.OwnerID <= 0 {
		logger.Error("Invalid project owner")
		return fmt.Errorf("%w", ErrInvalidProjectOwner)
	}

	// Validate project status
	if project.Status == "" {
		logger.Error("Invalid project status")
		return fmt.Errorf("%w", ErrInvalidProjectStatus)
	}

	// Validate project amount
	if project.AmountRequested.LessThanOrEqual(decimal.NewFromInt(0)) {
		logger.Error("Invalid project amount requested")
		return fmt.Errorf("%w", ErrInvalidProjectAmount)
	}

	// Validate project deadline
	if project.DeadlineAt != nil {
		now := time.Now()
		if project.DeadlineAt.Before(now) {
			logger.Error("Invalid project deadline: deadline is in the past")
			return fmt.Errorf("%w: deadline is in the past", ErrInvalidProjectDeadline)
		}
	}

	return nil
}

// Create создает новый проект (заглушка)
func (p *Project) Create(ctx context.Context, project model.Project) error {
	return nil
}

// Update обновляет существующий проект (заглушка)
func (p *Project) Update(ctx context.Context, project model.Project, userID int) error {
	return nil
}

// Delete удаляет проект (заглушка)
func (p *Project) Delete(ctx context.Context, id int, userID int) error {
	return nil
}

// GetByID возвращает проект по ID (заглушка)
func (p *Project) GetByID(ctx context.Context, id int) (model.Project, error) {
	return model.Project{}, nil
}

// List возвращает список проектов (заглушка)
func (p *Project) List(ctx context.Context, searchTerm string) ([]model.Project, error) {
	return nil, nil
}

// ListByOwnerID возвращает список проектов по ownerID (заглушка)
func (p *Project) ListByOwnerID(ctx context.Context, ownerID int64, searchTerm string) ([]model.Project, error) {
	return nil, nil
}

// GetPhotosByProjectID возвращает фото проекта (заглушка)
func (p *Project) GetPhotosByProjectID(ctx context.Context, projectID int) ([]model.ProjectImage, error) {
	return nil, nil
}
