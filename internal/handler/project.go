package handler

import (
	"github.com/CryptoCrowd/internal/service"
	"github.com/gofiber/fiber/v2"
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
	return nil
}

// Update handles the update of an existing project
func (h *ProjectHandler) Update(c *fiber.Ctx) error {
	return nil
}

// Delete handles the deletion of a project
func (h *ProjectHandler) Delete(c *fiber.Ctx) error {
	return nil
}

// GetByID handles the retrieval of a project by ID
func (h *ProjectHandler) GetByID(c *fiber.Ctx) error {
	return nil
}

// List handles the listing of projects
func (h *ProjectHandler) List(c *fiber.Ctx) error {
	return nil
}

// ListByOwnerID handles the listing of projects by owner ID
func (h *ProjectHandler) ListByOwnerID(c *fiber.Ctx) error {
	return nil
}

// GetPhotosByProjectID handles the retrieval of photos for a project
func (h *ProjectHandler) GetPhotosByProjectID(c *fiber.Ctx) error {
	return nil
}
