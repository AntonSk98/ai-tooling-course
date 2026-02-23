package main

import (
	"log"
	"petstore/internal/api"
	"petstore/internal/api/handlers"
	"petstore/internal/database"
	"petstore/internal/domain"
	"petstore/internal/repository"
	"petstore/pkg/config"

	_ "petstore/docs"

	"github.com/gofiber/fiber/v2"
)

// @title Petstore API
// @version 1.0
// @description REST API for a petstore application built with Go Fiber
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@petstore.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /
// @schemes http

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	if err := database.Connect(cfg.DBPath); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err := database.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	// Run migrations
	if err := database.Migrate(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize repositories
	petRepo := repository.NewPetRepository(database.DB)
	categoryRepo := repository.NewCategoryRepository(database.DB)

	// Initialize services
	petService := domain.NewPetService(petRepo, categoryRepo)
	categoryService := domain.NewCategoryService(categoryRepo)

	// Initialize handlers
	petHandler := handlers.NewPetHandler(petService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Petstore API v1.0",
	})

	// Serve static files (UI)
	app.Static("/", "../ui")

	// Setup API routes
	api.SetupRoutes(app, petHandler, categoryHandler)

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	log.Printf("UI available at: http://localhost:%s", cfg.Port)
	log.Printf("API documentation at: http://localhost:%s/swagger/index.html", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
