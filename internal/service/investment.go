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
	ErrInvalidInvestmentUser    = errors.New("invalid investment user")
	ErrInvalidInvestmentProject = errors.New("invalid investment project")
	ErrInvalidInvestmentAmount  = errors.New("invalid investment amount")
)

// InvestmentRepository defines the interface for investment repository operations
type InvestmentRepository interface {
	Create(ctx context.Context, investment model.Investment) error
	Update(ctx context.Context, investment model.Investment) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (model.Investment, error)
	GetByUserID(ctx context.Context, userID int64) ([]model.Investment, error)
	GetByProjectID(ctx context.Context, projectID int64) ([]model.Investment, error)
}

// Investment service implements business logic for investment operations
type Investment struct {
	repo    InvestmentRepository
	project ProjectRepository
}

// NewInvestment creates a new investment service
func NewInvestment(repo InvestmentRepository, project ProjectRepository) *Investment {
	return &Investment{
		repo:    repo,
		project: project,
	}
}

// validateInvestment validates investment data
func (i *Investment) validateInvestment(investment model.Investment) error {
	// Validate investment user
	if investment.UserID <= 0 {
		logger.Error("Invalid investment user")
		return fmt.Errorf("%w", ErrInvalidInvestmentUser)
	}

	// Validate investment project
	if investment.ProjectID <= 0 {
		logger.Error("Invalid investment project")
		return fmt.Errorf("%w", ErrInvalidInvestmentProject)
	}

	// Validate investment amount
	if investment.Amount <= 0 {
		logger.Error("Invalid investment amount")
		return fmt.Errorf("%w", ErrInvalidInvestmentAmount)
	}

	return nil
}

// Create creates a new investment with validation and updates project's raised amount
func (i *Investment) Create(ctx context.Context, investment model.Investment) error {
	// Validate investment data
	if err := i.validateInvestment(investment); err != nil {
		return err
	}

	// Set invested at time if not provided
	if investment.InvestedAt == nil {
		now := time.Now()
		investment.InvestedAt = &now
	}

	// Create investment in repository
	if err := i.repo.Create(ctx, investment); err != nil {
		return err
	}

	// Update project's raised amount
	return i.UpdateRaisedAmount(ctx, int(investment.ProjectID), investment.Amount)
}

// Update updates an existing investment with validation
func (i *Investment) Update(ctx context.Context, investment model.Investment, userID int64) error {
	// Get existing investment to check ownership
	existingInvestment, err := i.repo.GetByID(ctx, investment.ID)
	if err != nil {
		return err
	}

	// Check if user is the owner of the investment
	if existingInvestment.UserID != userID {
		logger.Errorf("User %d attempted to update investment %d owned by %d", userID, investment.ID, existingInvestment.UserID)
		return fmt.Errorf("unauthorized: only the investment owner can update the investment")
	}

	// Calculate amount difference for project update
	amountDiff := investment.Amount - existingInvestment.Amount

	// Validate investment data
	if err := i.validateInvestment(investment); err != nil {
		return err
	}

	// Update investment in repository
	if err := i.repo.Update(ctx, investment); err != nil {
		return err
	}

	// Update project's raised amount if amount changed
	if amountDiff != 0 {
		return i.UpdateRaisedAmount(ctx, int(investment.ProjectID), amountDiff)
	}

	return nil
}

// Delete deletes an investment
func (i *Investment) Delete(ctx context.Context, id int64, userID int64) error {
	// Get existing investment to check ownership
	existingInvestment, err := i.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Check if user is the owner of the investment
	if existingInvestment.UserID != userID {
		logger.Errorf("User %d attempted to delete investment %d owned by %d", userID, id, existingInvestment.UserID)
		return fmt.Errorf("unauthorized: only the investment owner can delete the investment")
	}

	// Delete investment from repository
	if err := i.repo.Delete(ctx, id); err != nil {
		return err
	}

	// Update project's raised amount (subtract the investment amount)
	return i.UpdateRaisedAmount(ctx, int(existingInvestment.ProjectID), -existingInvestment.Amount)
}

// GetByID retrieves an investment by ID
func (i *Investment) GetByID(ctx context.Context, id int64, userID int64) (model.Investment, error) {
	// Get investment
	investment, err := i.repo.GetByID(ctx, id)
	if err != nil {
		return model.Investment{}, err
	}

	// Check if user is the owner of the investment
	if investment.UserID != userID {
		logger.Errorf("User %d attempted to access investment %d owned by %d", userID, id, investment.UserID)
		return model.Investment{}, fmt.Errorf("unauthorized: only the investment owner can access the investment")
	}

	return investment, nil
}

// GetByUserID lists investments by user ID
func (i *Investment) GetByUserID(ctx context.Context, userID int64, requestingUserID int64) ([]model.Investment, error) {
	// Check if the requesting user is the owner
	if userID != requestingUserID {
		logger.Errorf("User %d attempted to access investments of user %d", requestingUserID, userID)
		return nil, fmt.Errorf("unauthorized: users can only access their own investments")
	}

	return i.repo.GetByUserID(ctx, userID)
}

// GetByProjectID lists investments by project ID
func (i *Investment) GetByProjectID(ctx context.Context, projectID int64) ([]model.Investment, error) {
	return i.repo.GetByProjectID(ctx, projectID)
}

// UpdateRaisedAmount updates the amount raised for a project
func (i *Investment) UpdateRaisedAmount(ctx context.Context, projectID int, amount float64) error {
	// Get an existing project
	project, err := i.project.GetByID(ctx, projectID)
	if err != nil {
		return err
	}

	// Update amount raised
	project.AmountRaised += amount

	// Update project in repository
	return i.project.Update(ctx, project)
}
