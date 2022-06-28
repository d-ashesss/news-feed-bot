//go:build integration
// +build integration

package integration

import (
	"cloud.google.com/go/firestore"
	"context"
	firestoreDb "github.com/d-ashesss/news-feed-bot/pkg/db/firestore"
	"github.com/d-ashesss/news-feed-bot/pkg/model"
	"github.com/jschoedt/go-firestorm"
	"log"
	"testing"
)

func TestModel(t *testing.T) {
	ctx := context.Background()
	fsc, err := firestore.NewClient(ctx, firestore.DetectProjectID)
	defer func(fsc *firestore.Client) {
		_ = fsc.Close()
	}(fsc)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}
	fso := firestorm.New(fsc, "ID", "Parent")

	resetData(ctx, fso)
	defer resetData(ctx, fso)

	categoryModel := firestoreDb.NewCategoryModel(fsc)

	cat1 := model.NewCategory("Cat1")
	if _, err := categoryModel.Create(ctx, cat1); err != nil {
		t.Fatalf("failed to create category %v: %v", cat1, err)
	}
	cat2 := model.NewCategory("Cat2")
	if _, err := categoryModel.Create(ctx, cat2); err != nil {
		t.Fatalf("failed to create category %v: %v", cat2, err)
	}
	cat3 := model.NewCategory("Cat3")
	if _, err := categoryModel.Create(ctx, cat3); err != nil {
		t.Fatalf("failed to create category %v: %v", cat3, err)
	}

	t.Run("CategoryModel", func(t *testing.T) {
		cat4 := model.NewCategory("Cat4")
		if _, err := categoryModel.Create(ctx, cat4); err != nil {
			t.Errorf("Create(%v): %v", cat4, err)
			return
		}

		t.Run("Get", func(t *testing.T) {
			cat, err := categoryModel.Get(ctx, cat4.ID)
			if err != nil {
				t.Errorf("Get(%q): %v", cat4.ID, err)
				return
			}
			if cat.ID != cat4.ID || cat.Name != cat4.Name {
				t.Errorf("Get(%q): got %v, want %v", cat4.ID, cat, cat4)
			}
		})

		t.Run("GetAll", func(t *testing.T) {
			cats, err := categoryModel.GetAll(ctx)
			if err != nil {
				t.Errorf("GetAll(): %v", err)
				return
			}
			wantNum := 4
			if len(cats) != wantNum {
				t.Errorf("GetAll(): got %d categories, want %d", len(cats), wantNum)
			}
		})

		t.Run("Delete", func(t *testing.T) {
			if err := categoryModel.Delete(ctx, cat4); err != nil {
				t.Errorf("Delete(%q): %v", cat4.ID, err)
			}
			_, err := categoryModel.Get(ctx, cat4.ID)
			if _, ok := err.(firestorm.NotFoundError); !ok {
				t.Errorf("Delete(): got %v, want NotFoundError for deleted category", err)
				return
			}
		})
	})
}

func resetData(ctx context.Context, fso *firestorm.FSClient) {
	cats := make([]model.Category, 0)
	_ = fso.NewRequest().QueryEntities(ctx, fso.NewRequest().ToCollection(model.Category{}).Query, &cats)()
	if err := fso.NewRequest().DeleteEntities(ctx, cats)(); err != nil {
		log.Fatalf("resetData: failed to delete categories: %v", err)
	}
}
