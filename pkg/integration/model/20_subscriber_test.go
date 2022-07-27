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

func TestSubscriberModel(t *testing.T) {
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
	updateModel := firestoreDb.NewUpdateModel(fsc)
	subscriberModel := firestoreDb.NewSubscriberModel(fsc, updateModel)

	s1 := model.NewSubscriber("U1")

	t.Run("Create", func(t *testing.T) {
		t.Run("nil subscriber", func(t *testing.T) {
			if _, err := subscriberModel.Create(ctx, nil); err != model.ErrInvalidSubscriber {
				t.Errorf("Create(%v): want ErrInvalidSubscriber; got %q", nil, err)
			}
		})

		t.Run("invalid UserID", func(t *testing.T) {
			s := &model.Subscriber{}
			if _, err := subscriberModel.Create(ctx, s); err != model.ErrInvalidSubscriberID {
				t.Errorf("Create(%v): want ErrInvalidSubscriberID; got %q", s, err)
			}
		})

		t.Run("valid subscriber", func(t *testing.T) {
			ID, err := subscriberModel.Create(ctx, s1)
			if err != nil {
				t.Fatalf("Create(%v): %v", s1, err)
			}
			if ID == "" {
				t.Errorf("Create(%v): got empty ID", s1)
			}
			if ID != s1.ID {
				t.Errorf("Create(%v): got ID %q; want %q", s1, ID, s1.ID)
			}
		})
	})

	t.Run("Get", func(t *testing.T) {
		t.Run("empty ID", func(t *testing.T) {
			if _, err := subscriberModel.Get(ctx, ""); err != model.ErrNotFound {
				t.Errorf("Get(%q): want ErrNotFound; got %q", "", err)
			}
		})

		t.Run("invalid ID", func(t *testing.T) {
			if _, err := subscriberModel.Get(ctx, "nothing"); err != model.ErrNotFound {
				t.Errorf("Get(%q): want ErrNotFound; got %q", "", err)
			}
		})

		t.Run("valid ID", func(t *testing.T) {
			s, err := subscriberModel.Get(ctx, s1.UserID)
			if err != nil {
				t.Fatalf("Get(%q): %v", s1.UserID, err)
			}
			if s == nil {
				t.Fatalf("Get(%q) = %v", s1.UserID, s)
			}
			if s.UserID != s1.UserID {
				t.Errorf("Get(%q) = %q", s1.UserID, s.UserID)
			}
		})
	})

	t.Run("Delete", func(t *testing.T) {
		t.Run("nil subscriber", func(t *testing.T) {
			if err := subscriberModel.Delete(ctx, nil); err != model.ErrInvalidSubscriber {
				t.Errorf("Delete(%v): got %q; want ErrInvalidSubscriber", nil, err)
			}
		})

		t.Run("invalid subscriber", func(t *testing.T) {
			s := &model.Subscriber{}
			if err := subscriberModel.Delete(ctx, s); err != model.ErrInvalidSubscriber {
				t.Errorf("Delete(%q): got %q; want ErrInvalidSubscriber", s, err)
			}
		})

		t.Run("valid subscriber", func(t *testing.T) {
			if err := subscriberModel.Delete(ctx, s1); err != nil {
				t.Fatalf("Delete(%q): %v", s1.UserID, err)
			}
			if _, err := subscriberModel.Get(ctx, s1.UserID); err != model.ErrNotFound {
				t.Errorf("Get(%q): got %q; want ErrNotFound for deleted subscriber", s1.UserID, err)
			}
		})

		t.Run("deleted updates", func(t *testing.T) {
			cat := model.NewCategory("CatX")
			if _, err := categoryModel.Create(ctx, cat); err != nil {
				t.Fatalf("categoryModel.Create(%v): %v", cat, err)
			}
			defer func(t *testing.T) {
				t.Helper()
				if err := categoryModel.Delete(ctx, cat); err != nil {
					t.Fatalf("categoryModel.Delete(%q): %v", cat.Name, err)
				}
			}(t)
			s := &model.Subscriber{UserID: "Sx"}
			if _, err := subscriberModel.Create(ctx, s); err != nil {
				t.Fatalf("subscriberModel.Create(%v): %v", s, err)
			}
			up := &model.Update{Subscriber: s, Category: cat, Title: "CatXUpX"}
			if _, err := updateModel.Create(ctx, up); err != nil {
				t.Fatalf("updateModel.Create(%v): %v", up, err)
			}
			if err := subscriberModel.Delete(ctx, s); err != nil {
				t.Fatalf("subscriberModel.Delete(%q): %v", s.UserID, err)
			}
			count, err := updateModel.GetCountInCategory(ctx, s, cat)
			if err != nil {
				t.Fatalf("updateModel.GetFromCategory(%q, %q): %v", s.UserID, cat.Name, err)
				return
			}
			if count != 0 {
				t.Fatalf("updateModel.GetFromCategory(%q, %q): = %d; want 0", s.UserID, cat.Name, count)
			}
		})
	})
}
