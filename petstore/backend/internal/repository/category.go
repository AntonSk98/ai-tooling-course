package repository

import (
	"errors"
	"petstore/internal/database"
	"petstore/internal/models"

	"gorm.io/gorm"
)

// CategoryRepository handles category data operations
type CategoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository creates a new category repository
func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// Create creates a new category
func (r *CategoryRepository) Create(category *models.Category) (*models.Category, error) {
	entity := &database.CategoryEntity{
		Name: *category.Name,
	}

	if err := r.db.Create(entity).Error; err != nil {
		return nil, err
	}

	return r.toModel(entity), nil
}

// Update updates an existing category
func (r *CategoryRepository) Update(category *models.Category) (*models.Category, error) {
	if category.Id == nil {
		return nil, errors.New("category ID is required")
	}

	var entity database.CategoryEntity
	if err := r.db.First(&entity, *category.Id).Error; err != nil {
		return nil, err
	}

	if category.Name != nil {
		entity.Name = *category.Name
	}

	if err := r.db.Save(&entity).Error; err != nil {
		return nil, err
	}

	return r.toModel(&entity), nil
}

// FindAll retrieves all categories
func (r *CategoryRepository) FindAll() ([]*models.Category, error) {
	var entities []database.CategoryEntity
	if err := r.db.Find(&entities).Error; err != nil {
		return nil, err
	}

	categories := make([]*models.Category, len(entities))
	for i, entity := range entities {
		categories[i] = r.toModel(&entity)
	}

	return categories, nil
}

// FindByID retrieves a category by ID
func (r *CategoryRepository) FindByID(id int64) (*models.Category, error) {
	var entity database.CategoryEntity
	if err := r.db.First(&entity, id).Error; err != nil {
		return nil, err
	}

	return r.toModel(&entity), nil
}

// Delete deletes a category by ID (soft delete)
func (r *CategoryRepository) Delete(id int64) error {
	result := r.db.Delete(&database.CategoryEntity{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// toModel converts CategoryEntity to Category model
func (r *CategoryRepository) toModel(entity *database.CategoryEntity) *models.Category {
	id := int64(entity.ID)
	name := entity.Name
	return &models.Category{
		Id:   &id,
		Name: &name,
	}
}
