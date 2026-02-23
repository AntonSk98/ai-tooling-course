package api

import (
	"petstore/internal/api/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	swagger "github.com/gofiber/swagger"
)

// SetupRoutes configures all application routes
func SetupRoutes(app *fiber.App, petHandler *handlers.PetHandler, categoryHandler *handlers.CategoryHandler) {
	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Swagger documentation
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	// Pet routes
	app.Post("/pet", petHandler.AddPet)
	app.Put("/pet", petHandler.UpdatePet)
	app.Get("/pet/findByStatus", petHandler.FindPetsByStatus)
	app.Get("/pet/:petId", petHandler.GetPetByID)
	app.Post("/pet/:petId", petHandler.UpdatePetWithForm)
	app.Delete("/pet/:petId", petHandler.DeletePet)

	// Category routes
	app.Post("/category", categoryHandler.AddCategory)
	app.Put("/category", categoryHandler.UpdateCategory)
	app.Get("/category/listAll", categoryHandler.GetAllCategories)
	app.Delete("/category/:categoryId", categoryHandler.DeleteCategory)
}
