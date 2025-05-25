package service

import (
	"context"
	"errors"
	"regexp"

	"github.com/CryptoCrowd/internal/model"
)

var (
	ErrInvalidUsername = errors.New("invalid username")
	ErrInvalidEmail    = errors.New("invalid email")
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
	reg, _ := regexp.Compile(`^[a-zA-Z0–9._%+-]+@[a-zA-Z0–9.-]+\.[a-zA-Z]{2,}$`)

	return &Account{
		repo:        repo,
		emailRegexp: reg,
	}
}

func (a *Account) Create(ctx context.Context, acc model.Account, plainPassword string) error {
	return nil
}

func (a *Account) Update(ctx context.Context, acc model.Account) error {
	return nil
}

func (a *Account) UpdatePassword(ctx context.Context, email string, newPassword string) error {
	return nil
}

func (a *Account) Delete(ctx context.Context, email string) error {
	return nil
}

func (a *Account) GetByEmail(ctx context.Context, email string) (model.Account, error) {
	return model.Account{}, nil
}

func (a *Account) List(ctx context.Context, searchTerm string) ([]model.Account, error) {
	return nil, nil
}
