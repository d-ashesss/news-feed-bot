package model

import (
	"context"
)

// Category represents a category entity.
type Category struct {
	ID   string
	Name string
}

// NewCategory initializes new Category.
func NewCategory(name string) *Category {
	return &Category{Name: name}
}

// CategoryModel is a data model for Category.
type CategoryModel interface {
	Create(ctx context.Context, c *Category) (string, error)
	Get(ctx context.Context, id string) (*Category, error)
	GetAll(ctx context.Context) ([]Category, error)
	Delete(ctx context.Context, c *Category) error
}
