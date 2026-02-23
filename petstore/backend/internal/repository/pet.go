package repository

import (
	"errors"
	"petstore/internal/database"
	"petstore/internal/models"

	"gorm.io/gorm"
)

// PetRepository handles pet data operations
type PetRepository struct {
	db *gorm.DB
}

// NewPetRepository creates a new pet repository
func NewPetRepository(db *gorm.DB) *PetRepository {
	return &PetRepository{db: db}
}

// Create creates a new pet
func (r *PetRepository) Create(pet *models.Pet) (*models.Pet, error) {
	entity := &database.PetEntity{
		Name:   pet.Name,
		Status: "available",
	}

	if pet.Status != nil {
		entity.Status = string(*pet.Status)
	}

	if pet.Category != nil && pet.Category.Id != nil {
		categoryID := uint(*pet.Category.Id)
		entity.CategoryID = &categoryID
	}

	if err := r.db.Create(entity).Error; err != nil {
		return nil, err
	}

	// Reload with associations
	if err := r.db.Preload("Category").First(entity, entity.ID).Error; err != nil {
		return nil, err
	}

	return r.toModel(entity), nil
}

// Update updates an existing pet
func (r *PetRepository) Update(pet *models.Pet) (*models.Pet, error) {
	if pet.Id == nil {
		return nil, errors.New("pet ID is required")
	}

	var entity database.PetEntity
	if err := r.db.First(&entity, *pet.Id).Error; err != nil {
		return nil, err
	}

	entity.Name = pet.Name

	if pet.Status != nil {
		entity.Status = string(*pet.Status)
	}

	if pet.Category != nil && pet.Category.Id != nil {
		categoryID := uint(*pet.Category.Id)
		entity.CategoryID = &categoryID
	}

	if err := r.db.Save(&entity).Error; err != nil {
		return nil, err
	}

	// Reload with associations
	if err := r.db.Preload("Category").First(&entity, entity.ID).Error; err != nil {
		return nil, err
	}

	return r.toModel(&entity), nil
}

// FindByID retrieves a pet by ID
func (r *PetRepository) FindByID(id int64) (*models.Pet, error) {
	var entity database.PetEntity
	if err := r.db.Preload("Category").First(&entity, id).Error; err != nil {
		return nil, err
	}

	return r.toModel(&entity), nil
}

// FindByStatus retrieves pets by status
func (r *PetRepository) FindByStatus(status string) ([]*models.Pet, error) {
	var entities []database.PetEntity
	if err := r.db.Preload("Category").Where("status = ?", status).Find(&entities).Error; err != nil {
		return nil, err
	}

	pets := make([]*models.Pet, len(entities))
	for i, entity := range entities {
		pets[i] = r.toModel(&entity)
	}

	return pets, nil
}

// Delete deletes a pet by ID
func (r *PetRepository) Delete(id int64) error {
	result := r.db.Delete(&database.PetEntity{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// toModel converts PetEntity to Pet model
func (r *PetRepository) toModel(entity *database.PetEntity) *models.Pet {
	id := int64(entity.ID)
	status := models.PetStatus(entity.Status)

	pet := &models.Pet{
		Id:     &id,
		Name:   entity.Name,
		Status: &status,
	}

	if entity.Category != nil {
		categoryID := int64(entity.Category.ID)
		categoryName := entity.Category.Name
		pet.Category = &models.Category{
			Id:   &categoryID,
			Name: &categoryName,
		}
	}

	return pet
}
