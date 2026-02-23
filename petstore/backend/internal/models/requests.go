package models

// CategoryReference is used when referencing an existing category by ID only
type CategoryReference struct {
	Id int64 `json:"id" example:"1"`
}

// CreatePetRequest defines the request body for creating a pet (without ID)
type CreatePetRequest struct {
	Name     string             `json:"name" example:"Balu"`
	Status   *PetStatus         `json:"status,omitempty" example:"available"`
	Category *CategoryReference `json:"category,omitempty"`
}

// CreateCategoryRequest defines the request body for creating a category (without ID)
type CreateCategoryRequest struct {
	Name string `json:"name" example:"Dogs"`
}

// ToCategory converts CreateCategoryRequest to Category model
func (r *CreateCategoryRequest) ToCategory() *Category {
	return &Category{
		Name: &r.Name,
	}
}

// ToPet converts CreatePetRequest to Pet model
func (r *CreatePetRequest) ToPet() *Pet {
	var category *Category
	if r.Category != nil {
		category = &Category{
			Id: &r.Category.Id,
		}
	}

	return &Pet{
		Name:     r.Name,
		Status:   r.Status,
		Category: category,
	}
}
