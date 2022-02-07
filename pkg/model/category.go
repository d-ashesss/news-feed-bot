package model

import (
	"context"
	"fmt"
)

// Category represents a category entity.
type Category struct {
	ID string `firestore:"-"`
}

// NewCategory initializes new Category.
func NewCategory() *Category {
	return &Category{}
}

// CategoryStore is an interface wrapper for a DB engine.
type CategoryStore interface {
	Create(ctx context.Context, u interface{}) (string, error)
	Get(ctx context.Context, id string, u interface{}) error
	Delete(ctx context.Context, id string) error
}

// CategoryModel data model for Category.
type CategoryModel struct {
	db CategoryStore
}

// NewCategoryModel initializes CategoryModel.
func NewCategoryModel(db CategoryStore) *CategoryModel {
	return &CategoryModel{db: db}
}

func (m *CategoryModel) Create(ctx context.Context, c *Category) (err error) {
	if len(c.ID) > 0 {
		return fmt.Errorf("create categore: provided categore is not new")
	}
	c.ID, err = m.db.Create(ctx, c)
	return
}

func (m *CategoryModel) Get(ctx context.Context, id string) (*Category, error) {
	var c Category
	if err := m.db.Get(ctx, id, &c); err != nil {
		return nil, err
	}
	c.ID = id
	return &c, nil
}

func (m *CategoryModel) Delete(ctx context.Context, c *Category) error {
	return m.db.Delete(ctx, c.ID)
}
