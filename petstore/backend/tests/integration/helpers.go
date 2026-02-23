package integration

import (
	"net/http/httptest"
	"petstore/internal/api"
	"petstore/internal/api/handlers"
	"petstore/internal/database"
	"petstore/internal/domain"
	"petstore/internal/repository"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func setupTestApp(t *testing.T) *fiber.App {
	// Create a unique in-memory database for each test (no cache sharing)
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Run migrations
	if err := db.AutoMigrate(&database.CategoryEntity{}, &database.PetEntity{}); err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	// Initialize repositories
	petRepo := repository.NewPetRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)

	// Initialize services
	petService := domain.NewPetService(petRepo, categoryRepo)
	categoryService := domain.NewCategoryService(categoryRepo)

	// Initialize handlers
	petHandler := handlers.NewPetHandler(petService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	// Create Fiber app
	app := fiber.New()
	api.SetupRoutes(app, petHandler, categoryHandler)

	return app
}

func TestHealthEndpoint(t *testing.T) {
	app := setupTestApp(t)

	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}
