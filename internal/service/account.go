package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/CryptoCrowd/internal/logger"
	"github.com/CryptoCrowd/internal/model"
)

var (
	ErrInvalidUsername = errors.New("invalid username")
	ErrInvalidEmail    = errors.New(("invalid email"))
	ErrInvalidRole     = errors.New("invalid role")
	ErrEmptyPass       = errors.New("password cannot be empty")
)

type AccountRepository interface {
	Create(ctx context.Context, acc model.Account, plainPassword string) error
	Update(ctx context.Context, acc model.Account) error
	UpdatePassword(ctx context.Context, email string, newPassword string) error
	Delete(ctx context.Context, email string) error
	GetByEmail(ctx context.Context, email string) (model.Account, error)
	List(ctx context.Context, searchTerm string) ([]model.Account, error)
}
type Account struct {
	repo        AccountRepository
	emailRegexp *regexp.Regexp
}

func NewAccount(repo AccountRepository) *Account {
	reg, _ := regexp.Compile(`/^[a-zA-Z0–9._%+-]+@[a-zA-Z0–9.-]+\.[a-zA-Z]{2,}$/`)

	return &Account{
		repo:        repo,
		emailRegexp: reg,
	}
}

func (a *Account) Create(ctx context.Context, acc model.Account, plainPassword string) error {

	if acc.Username == "" {
		logger.Error("Invalid username")
		return fmt.Errorf("%w", ErrInvalidUsername)
	}

	if !a.emailRegexp.MatchString(acc.Email) {
		logger.Errorf("Invalid email: %s", acc.Email)
		return fmt.Errorf("%w: %s", ErrInvalidEmail, acc.Email)
	}

	if acc.Role == "" {
		logger.Error("Invalid role")
		return fmt.Errorf("%w", ErrInvalidRole)
	}

	return a.repo.Create(ctx, acc, plainPassword)
}

func (a *Account) Update(ctx context.Context, acc model.Account) error {
	if acc.Username == "" {
		logger.Error("Invalid username")
		return fmt.Errorf("%w", ErrInvalidUsername)
	}

	if !a.emailRegexp.MatchString(acc.Email) {
		logger.Errorf("Invalid email: %s", acc.Email)
		return fmt.Errorf("%w: %s", ErrInvalidEmail, acc.Email)
	}

	if acc.Role == "" {
		logger.Error("Invalid role")
		return fmt.Errorf("%w", ErrInvalidRole)
	}

	return a.repo.Update(ctx, acc)
}

func (a *Account) UpdatePassword(ctx context.Context, email string, newPassword string) error {
	if !a.emailRegexp.MatchString(email) {
		logger.Errorf("Invalid email: %s", email)
		return fmt.Errorf("%w: %s", ErrInvalidEmail, email)
	}

	if newPassword == "" {
		logger.Error("Empty new password")
		return fmt.Errorf("%w", ErrEmptyPass)
	}

	return a.repo.UpdatePassword(ctx, email, newPassword)
}

func (a *Account) Delete(ctx context.Context, email string) error {
	if !a.emailRegexp.MatchString(email) {
		logger.Errorf("Invalid email: %s", email)
		return fmt.Errorf("%w: %s", ErrInvalidEmail, email)
	}

	return a.repo.Delete(ctx, email)
}

func (a *Account) GetByEmail(ctx context.Context, email string) (model.Account, error) {
	if !a.emailRegexp.MatchString(email) {
		logger.Errorf("Invalid email: %s", email)
		return model.Account{}, fmt.Errorf("%w: %s", ErrInvalidEmail, email)
	}

	return a.repo.GetByEmail(ctx, email)
}

func (a *Account) List(ctx context.Context, searchTerm string) ([]model.Account, error) {
	return a.repo.List(ctx, searchTerm)
}
