package database

import (
	"time"

	"gorm.io/gorm"
)

// PetEntity represents a pet in the database
type PetEntity struct {
	ID         uint           `gorm:"primaryKey;autoIncrement"`
	Name       string         `gorm:"not null;size:100"`
	Status     string         `gorm:"size:20;default:'available'"`
	CategoryID *uint          `gorm:"index"`
	Category   *CategoryEntity `gorm:"foreignKey:CategoryID"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

// TableName specifies the table name for PetEntity
func (PetEntity) TableName() string {
	return "pets"
}

// CategoryEntity represents a category in the database
type CategoryEntity struct {
	ID        uint           `gorm:"primaryKey;autoIncrement"`
	Name      string         `gorm:"not null;size:100;unique"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName specifies the table name for CategoryEntity
func (CategoryEntity) TableName() string {
	return "categories"
}
