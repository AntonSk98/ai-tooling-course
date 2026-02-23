package domain

import (
	"errors"
	"petstore/internal/models"
	"petstore/internal/repository"
)

// PetService handles business logic for pets
type PetService struct {
	repo         *repository.PetRepository
	categoryRepo *repository.CategoryRepository
}

// NewPetService creates a new pet service
func NewPetService(repo *repository.PetRepository, categoryRepo *repository.CategoryRepository) *PetService {
	return &PetService{
		repo:         repo,
		categoryRepo: categoryRepo,
	}
}

// CreatePet creates a new pet
func (s *PetService) CreatePet(pet *models.Pet) (*models.Pet, error) {
	if pet.Name == "" {
		return nil, errors.New("pet name is required")
	}

	if pet.Status != nil {
		if err := s.validateStatus(string(*pet.Status)); err != nil {
			return nil, err
		}
	}

	// Validate category exists if provided
	if pet.Category != nil && pet.Category.Id != nil {
		_, err := s.categoryRepo.FindByID(*pet.Category.Id)
		if err != nil {
			return nil, errors.New("category does not exist, please create it first")
		}
	}

	// Ignore client-provided ID - database will auto-generate
	pet.Id = nil

	return s.repo.Create(pet)
}

// UpdatePet updates an existing pet
func (s *PetService) UpdatePet(pet *models.Pet) (*models.Pet, error) {
	if pet.Id == nil {
		return nil, errors.New("pet ID is required")
	}

	if pet.Name == "" {
		return nil, errors.New("pet name is required")
	}

	if pet.Status != nil {
		if err := s.validateStatus(string(*pet.Status)); err != nil {
			return nil, err
		}
	}

	return s.repo.Update(pet)
}

// GetPetByID retrieves a pet by ID
func (s *PetService) GetPetByID(id int64) (*models.Pet, error) {
	return s.repo.FindByID(id)
}

// FindPetsByStatus retrieves pets by status
func (s *PetService) FindPetsByStatus(status string) ([]*models.Pet, error) {
	if err := s.validateStatus(status); err != nil {
		return nil, err
	}

	return s.repo.FindByStatus(status)
}

// UpdatePetWithForm updates a pet using form data
func (s *PetService) UpdatePetWithForm(id int64, name *string, status *string) (*models.Pet, error) {
	pet, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if name != nil && *name != "" {
		pet.Name = *name
	}

	if status != nil && *status != "" {
		if err := s.validateStatus(*status); err != nil {
			return nil, err
		}
		petStatus := models.PetStatus(*status)
		pet.Status = &petStatus
	}

	return s.repo.Update(pet)
}

// DeletePet deletes a pet by ID
func (s *PetService) DeletePet(id int64) error {
	return s.repo.Delete(id)
}

// validateStatus checks if the status is valid
func (s *PetService) validateStatus(status string) error {
	validStatuses := []string{"available", "pending", "sold"}
	for _, valid := range validStatuses {
		if status == valid {
			return nil
		}
	}
	return errors.New("invalid status value")
}
