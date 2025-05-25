package handler

import (
	"github.com/CryptoCrowd/internal/model"
	"github.com/CryptoCrowd/internal/service"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

// ProjectHandler handles HTTP requests related to projects
type ProjectHandler struct {
	projectService *service.Project
}

// NewProjectHandler creates a new project handler
func NewProjectHandler(projectService *service.Project) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
	}
}

// Create handles the creation of a new project
func (h *ProjectHandler) Create(c *fiber.Ctx) error {
	var project model.Project
	if err := c.BodyParser(&project); err != nil {
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

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Set owner ID from authenticated user
	project.OwnerID = userID

	if err := h.projectService.Create(c.Context(), project); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(project)
}

// Update handles the update of an existing project
func (h *ProjectHandler) Update(c *fiber.Ctx) error {
	var project model.Project
	if err := c.BodyParser(&project); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Get project ID from URL
	projectIDStr := c.Params("id")
	if projectIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Project ID is required",
		})
	}

	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid project ID",
		})
	}

	// Set project ID from URL
	project.ID = projectID

	// Get user ID from context (would be set by auth middleware)
	// For now, we'll get it from the request
	userIDStr := c.Get("X-User-ID")
	if userIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	if err := h.projectService.Update(c.Context(), project, userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(project)
}

// Delete handles the deletion of a project
func (h *ProjectHandler) Delete(c *fiber.Ctx) error {
	// Get project ID from URL
	projectIDStr := c.Params("id")
	if projectIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Project ID is required",
		})
	}

	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid project ID",
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

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	if err := h.projectService.Delete(c.Context(), projectID, userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Project deleted successfully",
	})
}

// GetByID handles the retrieval of a project by ID
func (h *ProjectHandler) GetByID(c *fiber.Ctx) error {
	// Get project ID from URL
	projectIDStr := c.Params("id")
	if projectIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Project ID is required",
		})
	}

	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid project ID",
		})
	}

	project, err := h.projectService.GetByID(c.Context(), projectID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(project)
}

// List handles the listing of projects
func (h *ProjectHandler) List(c *fiber.Ctx) error {
	searchTerm := c.Query("search", "")

	projects, err := h.projectService.List(c.Context(), searchTerm)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(projects)
}

// ListByOwnerID handles the listing of projects by owner ID
func (h *ProjectHandler) ListByOwnerID(c *fiber.Ctx) error {
	// Get owner ID from URL
	ownerIDStr := c.Params("owner_id")
	if ownerIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Owner ID is required",
		})
	}

	ownerID, err := strconv.ParseInt(ownerIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid owner ID",
		})
	}

	searchTerm := c.Query("search", "")

	projects, err := h.projectService.ListByOwnerID(c.Context(), ownerID, searchTerm)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(projects)
}

// GetPhotosByProjectID handles the retrieval of photos for a project
func (h *ProjectHandler) GetPhotosByProjectID(c *fiber.Ctx) error {
	// Get project ID from URL
	projectIDStr := c.Params("id")
	if projectIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Project ID is required",
		})
	}

	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid project ID",
		})
	}

	photos, err := h.projectService.GetPhotosByProjectID(c.Context(), projectID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(photos)
}
