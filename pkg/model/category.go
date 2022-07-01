package model

import (
	"context"
)

// Category represents a category entity.
type Category struct {
	ID   string // ID is an internal ID.
	Name string // Name is a name of the category.
}

// NewCategory initializes new Category.
func NewCategory(name string) *Category {
	return &Category{Name: name}
}

// CategoryModel is a data model for Category.
type CategoryModel interface {
	// Create saves a Category entity into the DB.
	Create(ctx context.Context, c *Category) (string, error)
	// Get retrieves a Category entity from the DB.
	Get(ctx context.Context, id string) (*Category, error)
	// GetAll retrieves all Category entities from the DB.
	GetAll(ctx context.Context) ([]Category, error)
	// Delete deletes a Category entity from the DB.
	Delete(ctx context.Context, c *Category) error
}
