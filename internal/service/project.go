package service

import (
	"context"
	"errors"
	"fmt"
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
	GetPhotosByProjectID(ctx context.Context, projectID int) ([]model.ProjectPhoto, error)
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
	if project.AmountRequested <= 0 {
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

// Create creates a new project with validation
func (p *Project) Create(ctx context.Context, project model.Project) error {
	// Validate project data
	if err := p.validateProject(project); err != nil {
		return err
	}

	// Set initial values for new project
	project.AmountRaised = 0
	if project.Status == "" {
		project.Status = "active"
	}

	// Create project in repository
	return p.repo.Create(ctx, project)
}

// Update updates an existing project with validation
func (p *Project) Update(ctx context.Context, project model.Project, userID int) error {
	// Get existing project to check ownership
	existingProject, err := p.repo.GetByID(ctx, project.ID)
	if err != nil {
		return err
	}

	// Check if user is the owner of the project
	if existingProject.OwnerID != userID {
		logger.Errorf("User %d attempted to update project %d owned by %d", userID, project.ID, existingProject.OwnerID)
		return fmt.Errorf("unauthorized: only the project owner can update the project")
	}

	// Validate project data
	if err := p.validateProject(project); err != nil {
		return err
	}

	// Preserve the original amount raised
	project.AmountRaised = existingProject.AmountRaised

	// Update project in repository
	return p.repo.Update(ctx, project)
}

// Delete deletes a project
func (p *Project) Delete(ctx context.Context, id int, userID int) error {
	// Get existing project to check ownership
	existingProject, err := p.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Check if user is the owner of the project
	if existingProject.OwnerID != userID {
		logger.Errorf("User %d attempted to delete project %d owned by %d", userID, id, existingProject.OwnerID)
		return fmt.Errorf("unauthorized: only the project owner can delete the project")
	}

	// Delete project from repository
	return p.repo.Delete(ctx, id)
}

// GetByID retrieves a project by ID
func (p *Project) GetByID(ctx context.Context, id int) (model.Project, error) {
	return p.repo.GetByID(ctx, id)
}

// List lists all projects with optional search
func (p *Project) List(ctx context.Context, searchTerm string) ([]model.Project, error) {
	return p.repo.List(ctx, searchTerm)
}

// ListByOwnerID lists projects by owner ID with optional search
func (p *Project) ListByOwnerID(ctx context.Context, ownerID int64, searchTerm string) ([]model.Project, error) {
	return p.repo.ListByOwnerID(ctx, ownerID, searchTerm)
}

// GetPhotosByProjectID gets photos for a project
func (p *Project) GetPhotosByProjectID(ctx context.Context, projectID int) ([]model.ProjectPhoto, error) {
	return p.repo.GetPhotosByProjectID(ctx, projectID)
}
