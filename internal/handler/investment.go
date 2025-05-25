package handler

import (
	"github.com/CryptoCrowd/internal/model"
	"github.com/CryptoCrowd/internal/service"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

// InvestmentHandler handles HTTP requests related to investments
type InvestmentHandler struct {
	investmentService *service.Investment
}

// NewInvestmentHandler creates a new investment handler
func NewInvestmentHandler(investmentService *service.Investment) *InvestmentHandler {
	return &InvestmentHandler{
		investmentService: investmentService,
	}
}

// Create handles the creation of a new investment
func (h *InvestmentHandler) Create(c *fiber.Ctx) error {
	var investment model.Investment
	if err := c.BodyParser(&investment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Get user ID from context (would be set by auth middleware)
	// For now, we'll get it from the request
	userIDStr := c.Get("X-User-ID")
	if userIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Set user ID from authenticated user
	investment.UserID = userID

	if err := h.investmentService.Create(c.Context(), investment); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(investment)
}

// Update handles the update of an existing investment
func (h *InvestmentHandler) Update(c *fiber.Ctx) error {
	var investment model.Investment
	if err := c.BodyParser(&investment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Get investment ID from URL
	investmentIDStr := c.Params("id")
	if investmentIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Investment ID is required",
		})
	}

	investmentID, err := strconv.ParseInt(investmentIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid investment ID",
		})
	}

	// Set investment ID from URL
	investment.ID = investmentID

	// Get user ID from context (would be set by auth middleware)
	// For now, we'll get it from the request
	userIDStr := c.Get("X-User-ID")
	if userIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	if err := h.investmentService.Update(c.Context(), investment, userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(investment)
}

// Delete handles the deletion of an investment
func (h *InvestmentHandler) Delete(c *fiber.Ctx) error {
	// Get investment ID from URL
	investmentIDStr := c.Params("id")
	if investmentIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Investment ID is required",
		})
	}

	investmentID, err := strconv.ParseInt(investmentIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid investment ID",
		})
	}

	// Get user ID from context (would be set by auth middleware)
	// For now, we'll get it from the request
	userIDStr := c.Get("X-User-ID")
	if userIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	if err := h.investmentService.Delete(c.Context(), investmentID, userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Investment deleted successfully",
	})
}

// GetByID handles the retrieval of an investment by ID
func (h *InvestmentHandler) GetByID(c *fiber.Ctx) error {
	// Get investment ID from URL
	investmentIDStr := c.Params("id")
	if investmentIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Investment ID is required",
		})
	}

	investmentID, err := strconv.ParseInt(investmentIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid investment ID",
		})
	}

	// Get user ID from context (would be set by auth middleware)
	// For now, we'll get it from the request
	userIDStr := c.Get("X-User-ID")
	if userIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	investment, err := h.investmentService.GetByID(c.Context(), investmentID, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(investment)
}

// GetByUserID handles the retrieval of investments by user ID
func (h *InvestmentHandler) GetByUserID(c *fiber.Ctx) error {
	// Get user ID from URL
	userIDStr := c.Params("user_id")
	if userIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Get requesting user ID from context (would be set by auth middleware)
	// For now, we'll get it from the request
	requestingUserIDStr := c.Get("X-User-ID")
	if requestingUserIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Requesting user ID is required",
		})
	}

	requestingUserID, err := strconv.ParseInt(requestingUserIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid requesting user ID",
		})
	}

	investments, err := h.investmentService.GetByUserID(c.Context(), userID, requestingUserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(investments)
}

// GetByProjectID handles the retrieval of investments by project ID
func (h *InvestmentHandler) GetByProjectID(c *fiber.Ctx) error {
	// Get project ID from URL
	projectIDStr := c.Params("project_id")
	if projectIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Project ID is required",
		})
	}

	projectID, err := strconv.ParseInt(projectIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid project ID",
		})
	}

	investments, err := h.investmentService.GetByProjectID(c.Context(), projectID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(investments)
}
