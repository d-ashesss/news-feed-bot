//go:build integration
// +build integration

package model

import (
	"cloud.google.com/go/firestore"
	"context"
	firestoreDb "github.com/d-ashesss/news-feed-bot/pkg/db/firestore"
	"github.com/d-ashesss/news-feed-bot/pkg/model"
	"testing"
	"time"
)

func TestFeedModel(t *testing.T) {
	ctx := context.Background()
	fsc, err := firestore.NewClient(ctx, firestore.DetectProjectID)
	defer func(fsc *firestore.Client) {
		_ = fsc.Close()
	}(fsc)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}
	resetData(t, ctx, fsc)

	feedModel := firestoreDb.NewFeedModel(fsc)
	categoryModel := firestoreDb.NewCategoryModel(fsc)

	cat1 := model.NewCategory("Cat1")
	if _, err := categoryModel.Create(ctx, cat1); err != nil {
		t.Fatalf("categoryModel.Create(%v): %v", cat1, err)
	}
	defer func() {
		if err := categoryModel.Delete(ctx, cat1); err != nil {
			t.Fatalf("categoryModel.Delete(%q): %v", cat1.Name, err)
		}
	}()

	cat1f1 := &model.Feed{Category: cat1, Title: "Cat1 Feed1"}

	t.Run("Create", func(t *testing.T) {
		t.Run("nil feed", func(t *testing.T) {
			var f *model.Feed
			if _, err := feedModel.Create(ctx, f); err != model.ErrInvalidFeed {
				t.Errorf("Create(%v): got %q; want ErrInvalidFeed", f, err)
			}
		})

		t.Run("nil category", func(t *testing.T) {
			f := &model.Feed{}
			if _, err := feedModel.Create(ctx, f); err != model.ErrInvalidCategory {
				t.Errorf("Create(%v): got %q; want ErrInvalidCategory", f, err)
			}
		})

		t.Run("invalid category", func(t *testing.T) {
			f := &model.Feed{Category: &model.Category{}}
			if _, err := feedModel.Create(ctx, f); err != model.ErrInvalidCategory {
				t.Errorf("Create(%v): got %q; want ErrInvalidCategory", f, err)
			}
		})

		t.Run("valid feed", func(t *testing.T) {
			ID, err := feedModel.Create(ctx, cat1f1)
			if err != nil {
				t.Fatalf("Create(%v): %v", cat1f1, err)
			}
			if ID == "" {
				t.Errorf("Create(%v): got empty ID", cat1f1)
			}
			if ID != cat1f1.ID {
				t.Errorf("Create(%v): got ID %q; want %q", cat1f1, ID, cat1f1.ID)
			}
		})
	})

	t.Run("Get", func(t *testing.T) {
		t.Run("empty id", func(t *testing.T) {
			if _, err := feedModel.Get(ctx, cat1, ""); err != model.ErrNotFound {
				t.Errorf("Get(%q, %q): got %q; want ErrNotFound", cat1.Name, "", err)
			}
		})

		t.Run("invalid id", func(t *testing.T) {
			if _, err := feedModel.Get(ctx, cat1, "nothing"); err != model.ErrNotFound {
				t.Errorf("Get(%q, %q): got %q; want ErrNotFound", cat1.Name, "nothing", err)
			}
		})

		t.Run("nil category", func(t *testing.T) {
			var cat *model.Category
			if _, err := feedModel.Get(ctx, cat, cat1f1.ID); err != model.ErrInvalidCategory {
				t.Errorf("Get(%v, %q): got %q; want ErrInvalidCategory", cat, cat1f1.Title, err)
			}
		})

		t.Run("empty category", func(t *testing.T) {
			cat := &model.Category{}
			if _, err := feedModel.Get(ctx, cat, cat1f1.ID); err != model.ErrInvalidCategory {
				t.Errorf("Get(%v, %q): got %q; want ErrInvalidCategory", cat, cat1f1.Title, err)
			}
		})

		t.Run("valid id", func(t *testing.T) {
			f, err := feedModel.Get(ctx, cat1, cat1f1.ID)
			if err != nil {
				t.Fatalf("Get(%q, %q): %v", cat1.Name, cat1f1.Title, err)
			}
			if f == nil {
				t.Fatalf("Get(%q, %q) = %v", cat1.Name, cat1f1.Title, f)
			}
			if f.ID != cat1f1.ID {
				t.Fatalf("Get(%q, %q) = %v", cat1.Name, cat1f1.Title, f)
			}
			if f.Category == nil || f.Category.ID != cat1.ID {
				t.Fatalf("Get(%q, %q) didn't load the category", cat1.Name, cat1f1.Title)
			}
		})
	})

	t.Run("GetAll", func(t *testing.T) {
		t.Run("nil category", func(t *testing.T) {
			var cat *model.Category
			if _, err := feedModel.GetAll(ctx, cat); err != model.ErrInvalidCategory {
				t.Errorf("GetAll(%v): got %q; want ErrInvalidCategory", cat, err)
			}
		})

		t.Run("empty category", func(t *testing.T) {
			cat := &model.Category{}
			if _, err := feedModel.GetAll(ctx, cat); err != model.ErrInvalidCategory {
				t.Errorf("GetAll(%v): got %q; want ErrInvalidCategory", cat, err)
			}
		})

		t.Run("valid category", func(t *testing.T) {
			feeds, err := feedModel.GetAll(ctx, cat1)
			if err != nil {
				t.Fatalf("GetAll(%v): %v", cat1.Name, err)
			}
			wantNum := 1
			if len(feeds) != wantNum {
				t.Errorf("GetAll(%v): got %d feeds; want %d", cat1.Name, len(feeds), wantNum)
			}
		})
	})

	t.Run("SetUpdated", func(t *testing.T) {
		t.Run("nil feed", func(t *testing.T) {
			var f *model.Feed
			u := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
			if err := feedModel.SetUpdated(ctx, f, u); err != model.ErrInvalidFeed {
				t.Errorf("SetUpdated(%v): got %q; want ErrInvalidFeed", f, err)
			}
		})

		t.Run("nil category", func(t *testing.T) {
			f := &model.Feed{ID: "test"}
			u := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
			if err := feedModel.SetUpdated(ctx, f, u); err != model.ErrInvalidCategory {
				t.Errorf("SetUpdated(%v): got %q; want ErrInvalidCategory", f, err)
			}
		})

		t.Run("invalid category", func(t *testing.T) {
			f := &model.Feed{ID: "test", Category: &model.Category{}}
			u := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
			if err := feedModel.SetUpdated(ctx, f, u); err != model.ErrInvalidCategory {
				t.Errorf("SetUpdated(%v): got %q; want ErrInvalidCategory", f, err)
			}
		})

		t.Run("valid feed", func(t *testing.T) {
			u := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
			if err := feedModel.SetUpdated(ctx, cat1f1, u); err != nil {
				t.Fatalf("SetUpdated(%q, %q): %v", cat1f1.Title, u.Format(time.RFC3339), err)
			}
			f, err := feedModel.Get(ctx, cat1, cat1f1.ID)
			if err != nil {
				t.Fatalf("Get(%q, %q): %v", cat1.Name, cat1f1.Title, err)
			}
			if !f.LastUpdate.Equal(u) {
				t.Errorf("SetUpdated(%q): got %q; want %q", cat1f1.Title, f.LastUpdate.Format(time.RFC3339), u.Format(time.RFC3339))
			}
		})
	})

	t.Run("Delete", func(t *testing.T) {
		t.Run("nil feed", func(t *testing.T) {
			var f *model.Feed
			if err := feedModel.Delete(ctx, f); err != model.ErrInvalidFeed {
				t.Errorf("Delete(%v): got %q; want ErrInvalidFeed", f, err)
			}
		})

		t.Run("nil category", func(t *testing.T) {
			f := &model.Feed{ID: "test"}
			if err := feedModel.Delete(ctx, f); err != model.ErrInvalidCategory {
				t.Errorf("Delete(%v): got %q; want ErrInvalidCategory", f, err)
			}
		})

		t.Run("invalid category", func(t *testing.T) {
			f := &model.Feed{ID: "test", Category: &model.Category{}}
			if err := feedModel.Delete(ctx, f); err != model.ErrInvalidCategory {
				t.Errorf("Delete(%v): got %q; want ErrInvalidCategory", f, err)
			}
		})

		t.Run("valid feed", func(t *testing.T) {
			if err := feedModel.Delete(ctx, cat1f1); err != nil {
				t.Fatalf("Delete(%q): got %q; want ErrInvalidCategory", cat1f1.Title, err)
			}
			if _, err := feedModel.Get(ctx, cat1, cat1f1.ID); err != model.ErrNotFound {
				t.Errorf("Get(%q, %q): got %q; want ErrNotFound for deleted feed", cat1.Name, cat1f1.Title, err)
			}
		})
	})
}
