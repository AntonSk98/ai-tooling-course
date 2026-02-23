package domain

import (
	"errors"
	"petstore/internal/models"
	"petstore/internal/repository"
)

// CategoryService handles business logic for categories
type CategoryService struct {
	repo *repository.CategoryRepository
}

// NewCategoryService creates a new category service
func NewCategoryService(repo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

// CreateCategory creates a new category
func (s *CategoryService) CreateCategory(category *models.Category) (*models.Category, error) {
	if category.Name == nil || *category.Name == "" {
		return nil, errors.New("category name is required")
	}

	// Ignore client-provided ID - database will auto-generate
	category.Id = nil

	return s.repo.Create(category)
}

// UpdateCategory updates an existing category
func (s *CategoryService) UpdateCategory(category *models.Category) (*models.Category, error) {
	if category.Id == nil {
		return nil, errors.New("category ID is required")
	}

	if category.Name == nil || *category.Name == "" {
		return nil, errors.New("category name is required")
	}

	return s.repo.Update(category)
}

// GetAllCategories retrieves all categories
func (s *CategoryService) GetAllCategories() ([]*models.Category, error) {
	return s.repo.FindAll()
}

// DeleteCategory deletes a category by ID
func (s *CategoryService) DeleteCategory(id int64) error {
	return s.repo.Delete(id)
}
