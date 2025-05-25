package handler

import (
	"github.com/CryptoCrowd/internal/model"
	"github.com/CryptoCrowd/internal/service"
	"github.com/gofiber/fiber/v2"
)

// AccountHandler handles HTTP requests related to accounts
type AccountHandler struct {
	accountService *service.Account
}

// NewAccountHandler creates a new account handler
func NewAccountHandler(accountService *service.Account) *AccountHandler {
	return &AccountHandler{
		accountService: accountService,
	}
}

// Create handles the creation of a new account
func (h *AccountHandler) Create(c *fiber.Ctx) error {
	var account model.Account
	if err := c.BodyParser(&account); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Get password from request
	password := c.FormValue("password")
	if password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Password is required",
		})
	}

	if err := h.accountService.Create(c.Context(), account, password); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(account)
}

// Update handles the update of an existing account
func (h *AccountHandler) Update(c *fiber.Ctx) error {
	var account model.Account
	if err := c.BodyParser(&account); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.accountService.Update(c.Context(), account); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(account)
}

// UpdatePassword handles the update of an account's password
func (h *AccountHandler) UpdatePassword(c *fiber.Ctx) error {
	email := c.FormValue("email")
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email is required",
		})
	}

	newPassword := c.FormValue("new_password")
	if newPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "New password is required",
		})
	}

	if err := h.accountService.UpdatePassword(c.Context(), email, newPassword); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Password updated successfully",
	})
}

// Delete handles the deletion of an account
func (h *AccountHandler) Delete(c *fiber.Ctx) error {
	email := c.Params("email")
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email is required",
		})
	}

	if err := h.accountService.Delete(c.Context(), email); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Account deleted successfully",
	})
}

// GetByEmail handles the retrieval of an account by email
func (h *AccountHandler) GetByEmail(c *fiber.Ctx) error {
	email := c.Params("email")
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email is required",
		})
	}

	account, err := h.accountService.GetByEmail(c.Context(), email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(account)
}

// List handles the listing of accounts
func (h *AccountHandler) List(c *fiber.Ctx) error {
	searchTerm := c.Query("search", "")

	accounts, err := h.accountService.List(c.Context(), searchTerm)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(accounts)
}
