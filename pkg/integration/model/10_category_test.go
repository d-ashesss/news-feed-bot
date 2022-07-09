//go:build integration
// +build integration

package model

import (
	"cloud.google.com/go/firestore"
	"context"
	firestoreDb "github.com/d-ashesss/news-feed-bot/pkg/db/firestore"
	"github.com/d-ashesss/news-feed-bot/pkg/model"
	"testing"
)

func TestCategoryModel(t *testing.T) {
	ctx := context.Background()
	fsc, err := firestore.NewClient(ctx, firestore.DetectProjectID)
	defer func(fsc *firestore.Client) {
		_ = fsc.Close()
	}(fsc)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}
	resetData(t, ctx, fsc)

	categoryModel := firestoreDb.NewCategoryModel(fsc)

	cat1 := model.NewCategory("Cat1")
	cat2 := model.NewCategory("Cat2")

	t.Run("Create", func(t *testing.T) {
		t.Run("nil category", func(t *testing.T) {
			var nilCat *model.Category
			if _, err := categoryModel.Create(ctx, nilCat); err != model.ErrInvalidCategory {
				t.Errorf("Create(%v): got %q; want ErrInvalidCategory", nil, err)
			}
		})

		t.Run("empty category", func(t *testing.T) {
			cat := &model.Category{}
			if _, err := categoryModel.Create(ctx, cat); err != model.ErrInvalidCategoryName {
				t.Errorf("Create(%v): got %q; want ErrInvalidCategoryName", cat, err)
			}
		})

		t.Run("valid category", func(t *testing.T) {
			ID, err := categoryModel.Create(ctx, cat1)
			if err != nil {
				t.Fatalf("Create(%v): %v", cat1, err)
			}
			if ID == "" {
				t.Errorf("Create(%v): got empty ID", cat1)
			}
			if ID != cat1.ID {
				t.Errorf("Create(%v): got ID %q; want %q", cat1, ID, cat1.ID)
			}
		})
	})

	t.Run("Get", func(t *testing.T) {
		t.Run("empty ID", func(t *testing.T) {
			if _, err := categoryModel.Get(ctx, ""); err != model.ErrNotFound {
				t.Errorf("Get(%q): got %q; want ErrNotFound", "", err)
			}
		})

		t.Run("invalid ID", func(t *testing.T) {
			if _, err := categoryModel.Get(ctx, "nothing"); err != model.ErrNotFound {
				t.Errorf("Get(%q): got %q; want ErrNotFound", "", err)
			}
		})

		t.Run("valid ID", func(t *testing.T) {
			cat, err := categoryModel.Get(ctx, cat1.ID)
			if err != nil {
				t.Fatalf("Get(%q): %v", cat1.Name, err)
			}
			if cat == nil {
				t.Fatalf("Get(%q) = %v", cat1.Name, cat)
			}
			if cat.ID != cat1.ID || cat.Name != cat1.Name {
				t.Errorf("Get(%q) = %v", cat1.Name, cat)
			}
		})
	})

	if _, err := categoryModel.Create(ctx, cat2); err != nil {
		t.Fatalf("Create(%v): %v", cat2, err)
	}
	defer func() {
		if err := categoryModel.Delete(ctx, cat2); err != nil {
			t.Fatalf("Delete(%v): %v", cat2, err)
		}
	}()

	t.Run("GetAll", func(t *testing.T) {
		cats, err := categoryModel.GetAll(ctx)
		if err != nil {
			t.Fatalf("GetAll(): %v", err)
		}
		wantNum := 2
		if len(cats) != wantNum {
			t.Errorf("GetAll(): got %d categories; want %d", len(cats), wantNum)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		t.Run("nil category", func(t *testing.T) {
			var nilCat *model.Category
			if err := categoryModel.Delete(ctx, nilCat); err != model.ErrInvalidCategory {
				t.Errorf("Delete(%v): got %q; want ErrInvalidCategory", nilCat, err)
			}
		})

		t.Run("empty category", func(t *testing.T) {
			cat := &model.Category{}
			if err := categoryModel.Delete(ctx, cat); err != model.ErrInvalidCategory {
				t.Errorf("Delete(%v): got %q; want ErrInvalidCategory", cat, err)
			}
		})

		t.Run("valid category", func(t *testing.T) {
			if err := categoryModel.Delete(ctx, cat1); err != nil {
				t.Fatalf("Delete(%v): %v", cat1, err)
			}
			if _, err := categoryModel.Get(ctx, cat1.ID); err != model.ErrNotFound {
				t.Errorf("Get(%q): got %q; want ErrNotFound for deleted category", cat1.Name, err)
			}
		})
	})
}
