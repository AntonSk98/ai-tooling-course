package handlers

import (
	"errors"
	"petstore/internal/domain"
	"petstore/internal/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// PetHandler handles pet-related HTTP requests
type PetHandler struct {
	service *domain.PetService
}

// NewPetHandler creates a new pet handler
func NewPetHandler(service *domain.PetService) *PetHandler {
	return &PetHandler{service: service}
}

// AddPet handles POST /pet
// @Summary Add a new pet to the store
// @Description Add a new pet with the input payload. ID will be auto-generated, omit it from the request.
// @Tags pets
// @Accept json
// @Produce json
// @Param pet body models.CreatePetRequest true "Pet object that needs to be added (without id)"
// @Success 200 {object} models.Pet
// @Failure 400 {object} map[string]string
// @Failure 422 {object} map[string]string
// @Router /pet [post]
func (h *PetHandler) AddPet(c *fiber.Ctx) error {
	var req models.CreatePetRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	pet := req.ToPet()
	createdPet, err := h.service.CreatePet(pet)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(createdPet)
}

// UpdatePet handles PUT /pet
// @Summary Update an existing pet
// @Description Update an existing pet by Id
// @Tags pets
// @Accept json
// @Produce json
// @Param pet body models.Pet true "Pet object that needs to be updated"
// @Success 200 {object} models.Pet
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 422 {object} map[string]string
// @Router /pet [put]
func (h *PetHandler) UpdatePet(c *fiber.Ctx) error {
	var pet models.Pet
	if err := c.BodyParser(&pet); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	updatedPet, err := h.service.UpdatePet(&pet)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Pet not found",
			})
		}
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(updatedPet)
}

// GetPetByID handles GET /pet/:petId
// @Summary Find pet by ID
// @Description Returns a single pet
// @Tags pets
// @Produce json
// @Param petId path int true "ID of pet to return"
// @Success 200 {object} models.Pet
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /pet/{petId} [get]
func (h *PetHandler) GetPetByID(c *fiber.Ctx) error {
	petID, err := strconv.ParseInt(c.Params("petId"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID supplied",
		})
	}

	pet, err := h.service.GetPetByID(petID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Pet not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unexpected error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(pet)
}

// FindPetsByStatus handles GET /pet/findByStatus
// @Summary Finds pets by status
// @Description Multiple status values can be provided
// @Tags pets
// @Produce json
// @Param status query string false "Status values" Enums(available, pending, sold)
// @Success 200 {array} models.Pet
// @Failure 400 {object} map[string]string
// @Router /pet/findByStatus [get]
func (h *PetHandler) FindPetsByStatus(c *fiber.Ctx) error {
	status := c.Query("status", "available")

	pets, err := h.service.FindPetsByStatus(status)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid status value",
		})
	}

	return c.Status(fiber.StatusOK).JSON(pets)
}

// UpdatePetWithForm handles POST /pet/:petId
// @Summary Updates a pet in the store with form data
// @Description Updates a pet with form data
// @Tags pets
// @Produce json
// @Param petId path int true "ID of pet that needs to be updated"
// @Param name query string false "Name of pet"
// @Param status query string false "Status of pet"
// @Success 200 {object} models.Pet
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /pet/{petId} [post]
func (h *PetHandler) UpdatePetWithForm(c *fiber.Ctx) error {
	petID, err := strconv.ParseInt(c.Params("petId"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID supplied",
		})
	}

	name := c.Query("name")
	status := c.Query("status")

	var namePtr, statusPtr *string
	if name != "" {
		namePtr = &name
	}
	if status != "" {
		statusPtr = &status
	}

	pet, err := h.service.UpdatePetWithForm(petID, namePtr, statusPtr)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Pet not found",
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(pet)
}

// DeletePet handles DELETE /pet/:petId
// @Summary Deletes a pet
// @Description Deletes a pet by ID
// @Tags pets
// @Param petId path int true "Pet id to delete"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /pet/{petId} [delete]
func (h *PetHandler) DeletePet(c *fiber.Ctx) error {
	petID, err := strconv.ParseInt(c.Params("petId"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid pet value",
		})
	}

	err = h.service.DeletePet(petID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Pet not found",
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid pet value",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Pet deleted",
	})
}
