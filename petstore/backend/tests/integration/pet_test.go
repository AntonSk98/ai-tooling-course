package integration

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"petstore/internal/models"
	"testing"
)

func TestAddPet(t *testing.T) {
	app := setupTestApp(t)

	t.Run("valid pet without category", func(t *testing.T) {
		status := models.PetStatusAvailable
		req := models.CreatePetRequest{
			Name:   "Balu",
			Status: &status,
		}

		body, _ := json.Marshal(req)
		httpReq := httptest.NewRequest("POST", "/pet", bytes.NewReader(body))
		httpReq.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(httpReq)
		if err != nil {
			t.Fatalf("Failed to perform request: %v", err)
		}

		if resp.StatusCode != 200 {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var result models.Pet
		json.NewDecoder(resp.Body).Decode(&result)

		if result.Name != req.Name {
			t.Errorf("Expected name %s, got %s", req.Name, result.Name)
		}
		if result.Id == nil {
			t.Error("Expected ID to be set")
		}
	})

	t.Run("valid pet with existing category", func(t *testing.T) {
		// First create a category
		categoryReq := models.CreateCategoryRequest{Name: "Dogs"}
		catBody, _ := json.Marshal(categoryReq)
		catHttpReq := httptest.NewRequest("POST", "/category", bytes.NewReader(catBody))
		catHttpReq.Header.Set("Content-Type", "application/json")
		catResp, _ := app.Test(catHttpReq)

		var createdCategory models.Category
		json.NewDecoder(catResp.Body).Decode(&createdCategory)

		// Now create a pet with that category reference
		status := models.PetStatusAvailable
		req := models.CreatePetRequest{
			Name:   "Rex",
			Status: &status,
			Category: &models.CategoryReference{
				Id: *createdCategory.Id,
			},
		}

		body, _ := json.Marshal(req)
		httpReq := httptest.NewRequest("POST", "/pet", bytes.NewReader(body))
		httpReq.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(httpReq)
		if err != nil {
			t.Fatalf("Failed to perform request: %v", err)
		}

		if resp.StatusCode != 200 {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var result models.Pet
		json.NewDecoder(resp.Body).Decode(&result)

		if result.Name != req.Name {
			t.Errorf("Expected name %s, got %s", req.Name, result.Name)
		}
		if result.Category == nil || result.Category.Id == nil {
			t.Error("Expected category to be set")
		}
	})

	t.Run("pet with non-existent category", func(t *testing.T) {
		status := models.PetStatusAvailable
		req := models.CreatePetRequest{
			Name:   "Max",
			Status: &status,
			Category: &models.CategoryReference{
				Id: 9999,
			},
		}

		body, _ := json.Marshal(req)
		httpReq := httptest.NewRequest("POST", "/pet", bytes.NewReader(body))
		httpReq.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(httpReq)

		if resp.StatusCode != 422 {
			t.Errorf("Expected status 422, got %d", resp.StatusCode)
		}
	})

	t.Run("missing required field", func(t *testing.T) {
		req := models.CreatePetRequest{
			Name: "",
		}

		body, _ := json.Marshal(req)
		httpReq := httptest.NewRequest("POST", "/pet", bytes.NewReader(body))
		httpReq.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(httpReq)

		if resp.StatusCode != 422 {
			t.Errorf("Expected status 422, got %d", resp.StatusCode)
		}
	})

	t.Run("invalid status", func(t *testing.T) {
		name := "Balu"
		status := models.PetStatus("invalid")
		pet := models.Pet{
			Name:   name,
			Status: &status,
		}

		body, _ := json.Marshal(pet)
		req := httptest.NewRequest("POST", "/pet", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		if resp.StatusCode != 422 {
			t.Errorf("Expected status 422, got %d", resp.StatusCode)
		}
	})
}

func TestUpdatePet(t *testing.T) {
	t.Run("valid update", func(t *testing.T) {
		app := setupTestApp(t)

		// Create a pet first
		name := "Balu"
		status := models.PetStatusAvailable
		pet := models.Pet{
			Name:   name,
			Status: &status,
		}
		body, _ := json.Marshal(pet)
		req := httptest.NewRequest("POST", "/pet", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		var createdPet models.Pet
		json.NewDecoder(resp.Body).Decode(&createdPet)

		// Update the pet
		updatedName := "Balu Updated"
		updatedStatus := models.PetStatusSold
		updatePet := models.Pet{
			Id:     createdPet.Id,
			Name:   updatedName,
			Status: &updatedStatus,
		}

		body, _ = json.Marshal(updatePet)
		req = httptest.NewRequest("PUT", "/pet", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ = app.Test(req)

		if resp.StatusCode != 200 {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var result models.Pet
		json.NewDecoder(resp.Body).Decode(&result)

		if result.Name != updatedName {
			t.Errorf("Expected name %s, got %s", updatedName, result.Name)
		}
	})

	t.Run("pet not found", func(t *testing.T) {
		app := setupTestApp(t)
		nonExistentID := int64(9999)
		updatePet := models.Pet{
			Id:   &nonExistentID,
			Name: "Test",
		}

		body, _ := json.Marshal(updatePet)
		req := httptest.NewRequest("PUT", "/pet", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		if resp.StatusCode != 404 {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})
}

func TestGetPetByID(t *testing.T) {
	t.Run("valid pet id", func(t *testing.T) {
		app := setupTestApp(t)

		// Create a pet first
		name := "Balu"
		pet := models.Pet{Name: name}
		body, _ := json.Marshal(pet)
		req := httptest.NewRequest("POST", "/pet", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		var createdPet models.Pet
		json.NewDecoder(resp.Body).Decode(&createdPet)

		// Get the pet by ID
		req = httptest.NewRequest("GET", "/pet/1", nil)
		resp, _ = app.Test(req)

		if resp.StatusCode != 200 {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var result models.Pet
		json.NewDecoder(resp.Body).Decode(&result)

		if result.Name != name {
			t.Errorf("Expected name %s, got %s", name, result.Name)
		}
	})

	t.Run("pet not found", func(t *testing.T) {
		app := setupTestApp(t)

		req := httptest.NewRequest("GET", "/pet/9999", nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != 404 {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})

	t.Run("invalid pet id", func(t *testing.T) {
		app := setupTestApp(t)
		req := httptest.NewRequest("GET", "/pet/invalid", nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != 400 {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})
}

func TestFindPetsByStatus(t *testing.T) {
	app := setupTestApp(t)

	// Create pets with different statuses
	statuses := []models.PetStatus{models.PetStatusAvailable, models.PetStatusPending, models.PetStatusSold}
	for _, status := range statuses {
		pet := models.Pet{Name: "Pet " + string(status), Status: &status}
		body, _ := json.Marshal(pet)
		req := httptest.NewRequest("POST", "/pet", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		app.Test(req)
	}

	t.Run("find by available status", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/pet/findByStatus?status=available", nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != 200 {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var pets []models.Pet
		json.NewDecoder(resp.Body).Decode(&pets)

		if len(pets) == 0 {
			t.Error("Expected at least one pet")
		}
	})

	t.Run("invalid status", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/pet/findByStatus?status=invalid", nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != 400 {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})
}

func TestUpdatePetWithForm(t *testing.T) {
	t.Run("valid form update", func(t *testing.T) {
		app := setupTestApp(t)

		// Create a pet first
		pet := models.Pet{Name: "Balu"}
		body, _ := json.Marshal(pet)
		req := httptest.NewRequest("POST", "/pet", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		app.Test(req)

		// Update pet with form
		req = httptest.NewRequest("POST", "/pet/1?name=UpdatedName&status=sold", nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != 200 {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var result models.Pet
		json.NewDecoder(resp.Body).Decode(&result)

		if result.Name != "UpdatedName" {
			t.Errorf("Expected name UpdatedName, got %s", result.Name)
		}
	})

	t.Run("pet not found", func(t *testing.T) {
		app := setupTestApp(t)
		req := httptest.NewRequest("POST", "/pet/9999?name=Test", nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != 404 {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})
}

func TestDeletePet(t *testing.T) {
	t.Run("valid delete", func(t *testing.T) {
		app := setupTestApp(t)

		// Create a pet first
		pet := models.Pet{Name: "Balu"}
		body, _ := json.Marshal(pet)
		req := httptest.NewRequest("POST", "/pet", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		app.Test(req)

		// Delete the pet
		req = httptest.NewRequest("DELETE", "/pet/1", nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != 200 {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		// Verify pet is deleted
		req = httptest.NewRequest("GET", "/pet/1", nil)
		resp, _ = app.Test(req)

		if resp.StatusCode != 404 {
			t.Errorf("Expected status 404 after deletion, got %d", resp.StatusCode)
		}
	})

	t.Run("pet not found", func(t *testing.T) {
		app := setupTestApp(t)
		req := httptest.NewRequest("DELETE", "/pet/9999", nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != 404 {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})
}
