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

func TestUpdateModel(t *testing.T) {
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
	subscriberModel := firestoreDb.NewSubscriberModel(fsc)
	updateModel := firestoreDb.NewUpdateModel(fsc)

	cat1 := model.NewCategory("Cat1")
	if _, err := categoryModel.Create(ctx, cat1); err != nil {
		t.Fatalf("categoryModel.Create(%v): %v", cat1, err)
	}
	cat2 := model.NewCategory("Cat2")
	if _, err := categoryModel.Create(ctx, cat2); err != nil {
		t.Fatalf("categoryModel.Create(%v): %v", cat2, err)
	}
	defer func() {
		if err := categoryModel.Delete(ctx, cat1); err != nil {
			t.Fatalf("categoryModel.Delete(%q): %v", cat1.Name, err)
		}
		if err := categoryModel.Delete(ctx, cat2); err != nil {
			t.Fatalf("categoryModel.Delete(%q): %v", cat2.Name, err)
		}
	}()

	s1 := &model.Subscriber{UserID: "S1"}
	if _, err := subscriberModel.Create(ctx, s1); err != nil {
		t.Fatalf("subscriberModel.Create(%v): %v", s1, err)
	}
	s2 := &model.Subscriber{UserID: "S2"}
	if _, err := subscriberModel.Create(ctx, s2); err != nil {
		t.Fatalf("subscriberModel.Create(%v): %v", s2, err)
	}
	defer func() {
		if err := subscriberModel.Delete(ctx, s1); err != nil {
			t.Fatalf("subscriberModel.Delete(%q): %v", s1.UserID, err)
		}
		if err := subscriberModel.Delete(ctx, s2); err != nil {
			t.Fatalf("subscriberModel.Delete(%q): %v", s2.UserID, err)
		}
	}()

	cat1up1 := &model.Update{Subscriber: s1, Category: cat1, Title: "Cat1Up1"}

	t.Run("Create", func(t *testing.T) {
		t.Run("nil subscriber", func(t *testing.T) {
			var s *model.Subscriber
			up := model.Update{Subscriber: s, Category: cat1}
			if _, err := updateModel.Create(ctx, &up); err != model.ErrInvalidSubscriber {
				t.Errorf("Create(%v): got %q; want ErrInvalidSubscriber", up, err)
			}
		})

		t.Run("invalid subscriber", func(t *testing.T) {
			up := model.Update{Subscriber: &model.Subscriber{}, Category: cat1}
			if _, err := updateModel.Create(ctx, &up); err != model.ErrInvalidSubscriber {
				t.Errorf("Create(%v): got %q; want ErrInvalidSubscriber", up, err)
			}
		})

		t.Run("nil category", func(t *testing.T) {
			var cat *model.Category
			up := model.Update{Subscriber: s1, Category: cat}
			if _, err := updateModel.Create(ctx, &up); err != model.ErrInvalidCategory {
				t.Errorf("Create(%v): got %q; want ErrInvalidCategory", up, err)
			}
		})

		t.Run("invalid category", func(t *testing.T) {
			up := model.Update{Subscriber: s1, Category: &model.Category{}}
			if _, err := updateModel.Create(ctx, &up); err != model.ErrInvalidCategory {
				t.Errorf("Create(%v): got %q; want ErrInvalidCategory", up, err)
			}
		})

		t.Run("valid update", func(t *testing.T) {
			ID, err := updateModel.Create(ctx, cat1up1)
			if err != nil {
				t.Fatalf("Create(%v): %v", cat1up1, err)
			}
			if ID == "" {
				t.Errorf("Create(%v): got empty ID", cat1up1)
			}
			if ID != cat1up1.ID {
				t.Errorf("Create(%v): got ID %q; want %q", cat1up1, ID, cat1up1.ID)
			}
		})
	})

	t.Run("GetFromCategory", func(t *testing.T) {
		t.Run("nil subscriber", func(t *testing.T) {
			if _, err := updateModel.GetFromCategory(ctx, nil, cat2); err != model.ErrInvalidSubscriber {
				t.Errorf("GetFromCategory(%v, %q): got %q; want ErrInvalidSubscriber", nil, cat2.Name, err)
				return
			}
		})

		t.Run("invalid subscriber", func(t *testing.T) {
			if _, err := updateModel.GetFromCategory(ctx, &model.Subscriber{}, cat2); err != model.ErrInvalidSubscriber {
				t.Errorf("GetFromCategory({}, %q): got %q; want ErrInvalidSubscriber", cat2.Name, err)
				return
			}
		})

		t.Run("nil category", func(t *testing.T) {
			if _, err := updateModel.GetFromCategory(ctx, s1, nil); err != model.ErrInvalidCategory {
				t.Errorf("GetFromCategory(%q, %v): got %q; want ErrInvalidCategory", s1.UserID, nil, err)
				return
			}
		})

		t.Run("invalid category", func(t *testing.T) {
			if _, err := updateModel.GetFromCategory(ctx, s1, &model.Category{}); err != model.ErrInvalidCategory {
				t.Errorf("GetFromCategory(%q, {}): got %q; want ErrInvalidCategory", s1.UserID, err)
				return
			}
		})

		t.Run("category has no updates", func(t *testing.T) {
			if _, err := updateModel.GetFromCategory(ctx, s1, cat2); err != model.ErrNoUpdates {
				t.Errorf("GetFromCategory(%q, %q): got %q; want ErrNoUpdates", s1.UserID, cat2.Name, err)
				return
			}
		})

		t.Run("subscriber has no updates", func(t *testing.T) {
			if _, err := updateModel.GetFromCategory(ctx, s2, cat1); err != model.ErrNoUpdates {
				t.Errorf("GetFromCategory(%q, %q): got %q; want ErrNoUpdates", s2.UserID, cat1.Name, err)
				return
			}
		})

		t.Run("subscriber has update", func(t *testing.T) {
			up, err := updateModel.GetFromCategory(ctx, s1, cat1)
			if err != nil {
				t.Errorf("GetFromCategory(%q, %q): %v", s1.UserID, cat1.Name, err)
				return
			}
			if up.Title != cat1up1.Title {
				t.Errorf("GetFromCategory(%q, %q) = %q; want %q", s1.UserID, cat1.Name, up.Title, cat1up1.Title)
			}
		})
	})

	t.Run("GetCountInCategory", func(t *testing.T) {
		t.Run("nil subscriber", func(t *testing.T) {
			if _, err := updateModel.GetCountInCategory(ctx, nil, cat2); err != model.ErrInvalidSubscriber {
				t.Errorf("GetCountInCategory(%v, %q): got %q; want ErrInvalidSubscriber", nil, cat2.Name, err)
				return
			}
		})

		t.Run("invalid subscriber", func(t *testing.T) {
			if _, err := updateModel.GetCountInCategory(ctx, &model.Subscriber{}, cat2); err != model.ErrInvalidSubscriber {
				t.Errorf("GetCountInCategory({}, %q): got %q; want ErrInvalidSubscriber", cat2.Name, err)
				return
			}
		})

		t.Run("nil category", func(t *testing.T) {
			if _, err := updateModel.GetCountInCategory(ctx, s1, nil); err != model.ErrInvalidCategory {
				t.Errorf("GetCountInCategory(%q, %v): got %q; want ErrInvalidCategory", s1.UserID, nil, err)
				return
			}
		})

		t.Run("invalid category", func(t *testing.T) {
			if _, err := updateModel.GetCountInCategory(ctx, s1, &model.Category{}); err != model.ErrInvalidCategory {
				t.Errorf("GetCountInCategory(%q, {}): got %q; want ErrInvalidCategory", s1.UserID, err)
				return
			}
		})

		t.Run("empty category", func(t *testing.T) {
			wantCount := 0
			count, err := updateModel.GetCountInCategory(ctx, s1, cat2)
			if err != nil {
				t.Errorf("GetCountInCategory(%q, %q): %v", s1.UserID, cat2.Name, err)
				return
			}
			if count != wantCount {
				t.Errorf("GetCountInCategory(%q, %q): got %d updates, want %d", s1.UserID, cat2.Name, count, wantCount)
			}
		})
		t.Run("not empty category", func(t *testing.T) {
			wantCount := 1
			count, err := updateModel.GetCountInCategory(ctx, s1, cat1)
			if err != nil {
				t.Errorf("GetCountInCategory(%q, %q): %v", s1.UserID, cat1.Name, err)
				return
			}
			if count != wantCount {
				t.Errorf("GetCountInCategory(%q, %q): got %d updates, want %d", s1.UserID, cat1.Name, count, wantCount)
			}
		})
	})

	t.Run("Delete", func(t *testing.T) {
		if err := updateModel.Delete(ctx, cat1up1); err != nil {
			t.Fatalf("Delete(%q): %v", cat1up1.Title, err)
		}
		if _, err = updateModel.GetFromCategory(ctx, s1, cat1); err != model.ErrNoUpdates {
			t.Errorf("GetFromCategory(%q, %q): got %q; want ErrNoUpdates for deleted update", s1.UserID, cat2.Name, err)
		}
	})

	t.Run("DeleteForSubscriber", func(t *testing.T) {
		s := &model.Subscriber{UserID: "Sx"}
		if _, err := subscriberModel.Create(ctx, s); err != nil {
			t.Fatalf("subscriberModel.Create(%v): %v", s, err)
		}
		defer func(t *testing.T) {
			t.Helper()
			if err := subscriberModel.Delete(ctx, s); err != nil {
				t.Fatalf("subscriberModel.Delete(%q): %v", s.UserID, err)
			}
		}(t)
		up := &model.Update{Subscriber: s, Category: cat1, Title: "Cat1UpX"}
		if _, err := updateModel.Create(ctx, up); err != nil {
			t.Fatalf("Create(%v): %v", up, err)
		}

		if err := updateModel.DeleteForSubscriber(ctx, s); err != nil {
			t.Fatalf("DeleteForSubscriber(%q): %v", up.Title, err)
		}

		count, err := updateModel.GetCountInCategory(ctx, s, cat1)
		if err != nil {
			t.Fatalf("GetFromCategory(%q, %q): %v", s.UserID, cat1.Name, err)
			return
		}
		if count != 0 {
			t.Fatalf("GetFromCategory(%q, %q): = %d; want 0", s.UserID, cat1.Name, count)
		}
	})
}
