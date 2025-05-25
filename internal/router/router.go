package router

import (
	"github.com/CryptoCrowd/internal/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// SetupRouter configures the Fiber router with all routes
func SetupRouter(
	accountHandler *handler.AccountHandler,
	projectHandler *handler.ProjectHandler,
	investmentHandler *handler.InvestmentHandler,
) *fiber.App {
	app := fiber.New(fiber.Config{
		// Enable strict routing
		StrictRouting: true,
		// Enable case-sensitive routing
		CaseSensitive: true,
		// Set app name
		AppName: "CryptoCrowd API",
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, X-User-ID",
	}))

	// API routes
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Account routes
	accounts := v1.Group("/accounts")
	accounts.Post("/", accountHandler.Create)
	accounts.Put("/", accountHandler.Update)
	accounts.Put("/password", accountHandler.UpdatePassword)
	accounts.Delete("/:email", accountHandler.Delete)
	accounts.Get("/:email", accountHandler.GetByEmail)
	accounts.Get("/", accountHandler.List)

	// Project routes
	projects := v1.Group("/projects")
	projects.Post("/", projectHandler.Create)
	projects.Put("/:id", projectHandler.Update)
	projects.Delete("/:id", projectHandler.Delete)
	projects.Get("/:id", projectHandler.GetByID)
	projects.Get("/", projectHandler.List)
	projects.Get("/owner/:owner_id", projectHandler.ListByOwnerID)
	projects.Get("/:id/photos", projectHandler.GetPhotosByProjectID)

	// Investment routes
	investments := v1.Group("/investments")
	investments.Post("/", investmentHandler.Create)
	investments.Put("/:id", investmentHandler.Update)
	investments.Delete("/:id", investmentHandler.Delete)
	investments.Get("/:id", investmentHandler.GetByID)
	investments.Get("/user/:user_id", investmentHandler.GetByUserID)
	investments.Get("/project/:project_id", investmentHandler.GetByProjectID)

	return app
}
