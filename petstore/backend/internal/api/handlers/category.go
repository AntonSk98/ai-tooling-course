package handlers

import (
	"errors"
	"petstore/internal/domain"
	"petstore/internal/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// CategoryHandler handles category-related HTTP requests
type CategoryHandler struct {
	service *domain.CategoryService
}

// NewCategoryHandler creates a new category handler
func NewCategoryHandler(service *domain.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

// AddCategory handles POST /category
// @Summary Add a new category to the store
// @Description Add a new category with the input payload. ID will be auto-generated, omit it from the request.
// @Tags categories
// @Accept json
// @Produce json
// @Param category body models.CreateCategoryRequest true "Category object that needs to be added (without id)"
// @Success 200 {object} models.Category
// @Failure 400 {object} map[string]string
// @Failure 422 {object} map[string]string
// @Router /category [post]
func (h *CategoryHandler) AddCategory(c *fiber.Ctx) error {
	var req models.CreateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	category := req.ToCategory()
	createdCategory, err := h.service.CreateCategory(category)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(createdCategory)
}

// UpdateCategory handles PUT /category
// @Summary Update an existing category
// @Description Update an existing category by Id
// @Tags categories
// @Accept json
// @Produce json
// @Param category body models.Category true "Category object that needs to be updated"
// @Success 200 {object} models.Category
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 422 {object} map[string]string
// @Router /category [put]
func (h *CategoryHandler) UpdateCategory(c *fiber.Ctx) error {
	var category models.Category
	if err := c.BodyParser(&category); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	updatedCategory, err := h.service.UpdateCategory(&category)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Category not found",
			})
		}
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(updatedCategory)
}

// GetAllCategories handles GET /category/listAll
// @Summary Get all categories
// @Description Returns a list of all categories
// @Tags categories
// @Produce json
// @Success 200 {array} models.Category
// @Failure 500 {object} map[string]string
// @Router /category/listAll [get]
func (h *CategoryHandler) GetAllCategories(c *fiber.Ctx) error {
	categories, err := h.service.GetAllCategories()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unexpected error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(categories)
}

// DeleteCategory handles DELETE /category/:categoryId
// @Summary Deletes a category
// @Description Delete a category by ID
// @Tags categories
// @Param categoryId path int true "Category ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /category/{categoryId} [delete]
func (h *CategoryHandler) DeleteCategory(c *fiber.Ctx) error {
	categoryID, err := c.ParamsInt("categoryId")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid category ID",
		})
	}

	if err := h.service.DeleteCategory(int64(categoryID)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Category not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete category",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Category deleted successfully",
	})
}
