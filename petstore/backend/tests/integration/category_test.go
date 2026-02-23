package integration

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"petstore/internal/models"
	"testing"
)

func TestAddCategory(t *testing.T) {
	app := setupTestApp(t)

	t.Run("valid category", func(t *testing.T) {
		name := "Dogs"
		category := models.Category{
			Name: &name,
		}

		body, _ := json.Marshal(category)
		req := httptest.NewRequest("POST", "/category", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to perform request: %v", err)
		}

		if resp.StatusCode != 200 {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var result models.Category
		json.NewDecoder(resp.Body).Decode(&result)

		if *result.Name != name {
			t.Errorf("Expected name %s, got %s", name, *result.Name)
		}
		if result.Id == nil {
			t.Error("Expected ID to be set")
		}
	})

	t.Run("missing required field", func(t *testing.T) {
		category := models.Category{}

		body, _ := json.Marshal(category)
		req := httptest.NewRequest("POST", "/category", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		if resp.StatusCode != 422 {
			t.Errorf("Expected status 422, got %d", resp.StatusCode)
		}
	})
}

func TestUpdateCategory(t *testing.T) {
	t.Run("valid update", func(t *testing.T) {
		app := setupTestApp(t)

		// Create a category first
		name := "Birds"
		category := models.Category{Name: &name}
		body, _ := json.Marshal(category)
		req := httptest.NewRequest("POST", "/category", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		var createdCategory models.Category
		json.NewDecoder(resp.Body).Decode(&createdCategory)

		// Update the category
		updatedName := "Cats"
		updateCategory := models.Category{
			Id:   createdCategory.Id,
			Name: &updatedName,
		}

		body, _ = json.Marshal(updateCategory)
		req = httptest.NewRequest("PUT", "/category", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ = app.Test(req)

		if resp.StatusCode != 200 {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var result models.Category
		json.NewDecoder(resp.Body).Decode(&result)

		if *result.Name != updatedName {
			t.Errorf("Expected name %s, got %s", updatedName, *result.Name)
		}
	})

	t.Run("category not found", func(t *testing.T) {
		app := setupTestApp(t)

		nonExistentID := int64(9999)
		name := "Test"
		updateCategory := models.Category{
			Id:   &nonExistentID,
			Name: &name,
		}

		body, _ := json.Marshal(updateCategory)
		req := httptest.NewRequest("PUT", "/category", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		if resp.StatusCode != 404 {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})

	t.Run("missing required field", func(t *testing.T) {
		app := setupTestApp(t)

		// Create a category first
		name := "Fish"
		category := models.Category{Name: &name}
		body, _ := json.Marshal(category)
		req := httptest.NewRequest("POST", "/category", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		var createdCategory models.Category
		json.NewDecoder(resp.Body).Decode(&createdCategory)

		// Try to update without name
		updateCategory := models.Category{
			Id: createdCategory.Id,
		}

		body2, _ := json.Marshal(updateCategory)
		req2 := httptest.NewRequest("PUT", "/category", bytes.NewReader(body2))
		req2.Header.Set("Content-Type", "application/json")

		resp2, _ := app.Test(req2)

		if resp2.StatusCode != 422 {
			t.Errorf("Expected status 422, got %d", resp2.StatusCode)
		}
	})
}

func TestGetAllCategories(t *testing.T) {
	t.Run("get all categories", func(t *testing.T) {
		app := setupTestApp(t)

		// Create multiple categories
		categories := []string{"Dogs", "Cats", "Birds"}
		for _, catName := range categories {
			name := catName
			category := models.Category{Name: &name}
			body, _ := json.Marshal(category)
			req := httptest.NewRequest("POST", "/category", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			app.Test(req)
		}

		// Get all categories
		req := httptest.NewRequest("GET", "/category/listAll", nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != 200 {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var result []models.Category
		json.NewDecoder(resp.Body).Decode(&result)

		if len(result) != 3 {
			t.Errorf("Expected 3 categories, got %d", len(result))
		}
	})
}
