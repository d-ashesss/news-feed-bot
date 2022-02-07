package model

import (
	"context"
	"errors"
	"github.com/d-ashesss/news-feed-bot/mocks"
	"github.com/stretchr/testify/mock"
	"testing"
)

var (
	testCategoryError = errors.New("test error")
)

const (
	testCategoryID = "test-ID"
)

func TestCategoryModel_Create(t *testing.T) {
	t.Run("user exists", func(t *testing.T) {
		db := &mocks.CategoryStore{}

		m := NewCategoryModel(db)
		c := NewCategory()
		c.ID = "new-ID"

		if err := m.Create(context.Background(), c); err == nil {
			t.Errorf("Create() expected error")
		}
	})

	t.Run("create error", func(t *testing.T) {
		db := &mocks.CategoryStore{}
		db.On("Create", mock.Anything, mock.Anything).Return("", testCategoryError)

		m := NewCategoryModel(db)
		c := NewCategory()

		if err := m.Create(context.Background(), c); err != testCategoryError {
			t.Errorf("Create() got error %q, want %q", err, testCategoryError)
		}
	})

	t.Run("user created", func(t *testing.T) {
		db := &mocks.CategoryStore{}
		db.On("Create", mock.Anything, mock.Anything).Return(func(ctx context.Context, o interface{}) string {
			return testCategoryID
		}, nil)

		m := NewCategoryModel(db)
		c := NewCategory()

		if err := m.Create(context.Background(), c); err != nil {
			t.Errorf("Create() unexpected error = %v", err)
		}
		if c.ID != testCategoryID {
			t.Errorf("Create() got user ID %q, want %q", c.ID, testCategoryID)
		}
	})
}

func TestCategoryModel_Get(t *testing.T) {
	t.Run("user not found", func(t *testing.T) {
		db := &mocks.CategoryStore{}
		db.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(testCategoryError)
		m := NewCategoryModel(db)

		c, err := m.Get(context.Background(), testCategoryID)
		if err != testCategoryError {
			t.Errorf("Get() got error %q, want %q", err, testCategoryError)
		}
		if c != nil {
			t.Errorf("Get() returned unexpected user")
		}
	})

	t.Run("user found", func(t *testing.T) {
		db := &mocks.CategoryStore{}
		db.On("Get", mock.Anything, mock.MatchedBy(func(id string) bool {
			return id == testCategoryID
		}), mock.Anything).Return(nil)
		m := NewCategoryModel(db)

		c, err := m.Get(context.Background(), testCategoryID)
		if err != nil {
			t.Errorf("Get() unexpected error = %v", err)
		}
		if c.ID != testCategoryID {
			t.Errorf("Get() got user with ID %q, want %q", c.ID, testCategoryID)
		}
	})
}

func TestCategoryModel_Delete(t *testing.T) {
	t.Run("category deleted", func(t *testing.T) {
		db := &mocks.CategoryStore{}
		db.On("Delete", mock.Anything, mock.MatchedBy(func(id string) bool {
			return id == testCategoryID
		})).Return(nil)
		m := NewCategoryModel(db)
		c := NewCategory()
		c.ID = testCategoryID
		if err := m.Delete(context.Background(), c); err != nil {
			t.Errorf("Delete() unexpected error = %v", err)
		}
	})
}
