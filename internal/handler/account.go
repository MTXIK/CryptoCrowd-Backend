package handler

import (
	"context"
	"github.com/CryptoCrowd/internal/model"
	"github.com/gofiber/fiber/v2"
)

type AccountServiceInterface interface {
	Create(ctx context.Context, acc model.Account, plainPassword string) error
	Update(ctx context.Context, acc model.Account) error
	UpdatePassword(ctx context.Context, email string, newPassword string) error
	Delete(ctx context.Context, email string) error
	GetByEmail(ctx context.Context, email string) (model.Account, error)
	List(ctx context.Context, searchTerm string) ([]model.Account, error)
}

// AccountHandler handles HTTP requests related to accounts
type AccountHandler struct {
	accountService AccountServiceInterface
}

// NewAccountHandler creates a new account handler
func NewAccountHandler(accountService AccountServiceInterface) *AccountHandler {
	return &AccountHandler{
		accountService: accountService,
	}
}

// Create handles the creation of a new account
func (h *AccountHandler) Create(c *fiber.Ctx) error {
	return nil
}

// Update handles the update of an existing account
func (h *AccountHandler) Update(c *fiber.Ctx) error {
	return nil
}

// UpdatePassword handles the update of an account's password
func (h *AccountHandler) UpdatePassword(c *fiber.Ctx) error {
	return nil
}

// Delete handles the deletion of an account
func (h *AccountHandler) Delete(c *fiber.Ctx) error {
	return nil
}

// GetByEmail handles the retrieval of an account by email
func (h *AccountHandler) GetByEmail(c *fiber.Ctx) error {
	return nil
}

// List handles the listing of accounts
func (h *AccountHandler) List(c *fiber.Ctx) error {
	return nil
}
