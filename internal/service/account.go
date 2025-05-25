package service

import (
	"context"
	"github.com/CryptoCrowd/internal/model"
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
	repo AccountRepository
}

func (a *Account) Create(ctx context.Context, acc model.Account, plainPassword string) error {
	if acc.Username == "" {

	}
	return a.repo.Create(ctx, acc, plainPassword)
}
