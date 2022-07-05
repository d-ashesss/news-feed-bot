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
	updateModel := firestoreDb.NewUpdateModel(fsc)
	subscriptionModel := firestoreDb.NewSubscriptionModel(fsc, categoryModel, updateModel)

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

	s1 := model.NewSubscriber("U1")
	if _, err = subscriptionModel.CreateSubscriber(ctx, s1); err != nil {
		t.Fatalf("failed to create subscriber %v: %v", s1, err)
	}
	s2 := model.NewSubscriber("U2")
	if _, err = subscriptionModel.CreateSubscriber(ctx, s2); err != nil {
		t.Fatalf("failed to create subscriber %v: %v", s2, err)
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
				t.Errorf("Get(%q): %v", cat4.Name, err)
				return
			}
			if cat.ID != cat4.ID || cat.Name != cat4.Name {
				t.Errorf("Get(%q): got %v, want %v", cat4.Name, cat, cat4)
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
				t.Errorf("Delete(%q): %v", cat4.Name, err)
			}
			_, err := categoryModel.Get(ctx, cat4.ID)
			if _, ok := err.(firestorm.NotFoundError); !ok {
				t.Errorf("Get(%q): got %v, want NotFoundError for deleted category", cat4.Name, err)
				return
			}
		})
	})

	t.Run("UpdateModel", func(t *testing.T) {
		upTitle := "Cat1Up1"
		t.Run("Create", func(t *testing.T) {
			t.Run("invalid subscriber", func(t *testing.T) {
				up := model.Update{Subscriber: nil, Category: cat1, Title: upTitle}
				if _, err := updateModel.Create(ctx, &up); err != model.ErrInvalidSubscriber {
					t.Errorf("Create(%v): want ErrInvalidSubscriber, got %q", up, err)
					return
				}
			})

			t.Run("invalid category", func(t *testing.T) {
				up := model.Update{Subscriber: s1, Category: nil, Title: upTitle}
				if _, err := updateModel.Create(ctx, &up); err != model.ErrInvalidCategory {
					t.Errorf("Create(%v): want ErrInvalidCategory, got %q", up, err)
					return
				}
			})

			t.Run("valid update", func(t *testing.T) {
				up := model.Update{Subscriber: s1, Category: cat1, Title: upTitle}
				_, err := updateModel.Create(ctx, &up)
				if err != nil {
					t.Errorf("Create(%v): %v", up, err)
					return
				}
				if len(up.ID) == 0 {
					t.Errorf("Create(): want ID field to be set")
				}
			})
		})

		t.Run("GetFromCategory", func(t *testing.T) {
			t.Run("invalid subscriber", func(t *testing.T) {
				if _, err := updateModel.GetFromCategory(ctx, nil, cat2); err != model.ErrInvalidSubscriber {
					t.Errorf("GetFromCategory(%v, %q): want ErrInvalidSubscriber, got %q", nil, cat2.Name, err)
					return
				}
			})

			t.Run("invalid category", func(t *testing.T) {
				if _, err := updateModel.GetFromCategory(ctx, s1, nil); err != model.ErrInvalidCategory {
					t.Errorf("GetFromCategory(%q, %v): want ErrInvalidCategory, got %q", s1.UserID, nil, err)
					return
				}
			})

			t.Run("category has no updates", func(t *testing.T) {
				if _, err := updateModel.GetFromCategory(ctx, s1, cat2); err == nil {
					t.Errorf("GetFromCategory(%q, %q): want error", s1.UserID, cat2.Name)
					return
				}
			})

			t.Run("subscriber has no updates", func(t *testing.T) {
				if _, err := updateModel.GetFromCategory(ctx, s2, cat1); err == nil {
					t.Errorf("GetFromCategory(%q, %q): want error", s2.UserID, cat1.Name)
					return
				}
			})

			t.Run("has update", func(t *testing.T) {
				up, err := updateModel.GetFromCategory(ctx, s1, cat1)
				if err != nil {
					t.Errorf("GetFromCategory(%q, %q): %v", s1.UserID, cat1.Name, err)
					return
				}
				if up.Title != upTitle {
					t.Errorf("GetFromCategory(%q, %q): got %v, want %v", s1.UserID, cat1.Name, up.Title, upTitle)
				}
			})
		})

		t.Run("GetCountInCategory", func(t *testing.T) {
			t.Run("invalid subscriber", func(t *testing.T) {
				if _, err := updateModel.GetCountInCategory(ctx, nil, cat2); err != model.ErrInvalidSubscriber {
					t.Errorf("GetCountInCategory(%v, %q): want ErrInvalidSubscriber, got %q", nil, cat2.Name, err)
					return
				}
			})

			t.Run("invalid category", func(t *testing.T) {
				if _, err := updateModel.GetCountInCategory(ctx, s1, nil); err != model.ErrInvalidCategory {
					t.Errorf("GetCountInCategory(%q, %v): want ErrInvalidCategory, got %q", s1.UserID, nil, err)
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
			up, err := updateModel.GetFromCategory(ctx, s1, cat1)
			if err != nil {
				t.Errorf("GetFromCategory(%q, %q): %v", s1.UserID, cat1.Name, err)
				return
			}
			if err := updateModel.Delete(ctx, up); err != nil {
				t.Errorf("Delete(%v): %v", up.Title, err)
			}
			if _, err = updateModel.GetFromCategory(ctx, s1, cat1); err == nil {
				t.Errorf("GetFromCategory(%q, %q): want error for deleted update", s1.UserID, cat2.Name)
				return
			}
		})
	})

	t.Run("SubscriptionModel", func(t *testing.T) {
		t.Run("GetSubscriber", func(t *testing.T) {
			s, err := subscriptionModel.GetSubscriber(ctx, s1.ID)
			if err != nil {
				t.Errorf("GetSubscriber(%q): %v", s1.UserID, err)
				return
			}
			if s.UserID != s1.UserID {
				t.Errorf("GetSubscriber(%q) = %q", s1.UserID, s.UserID)
			}
		})

		t.Run("GetSubscriptionStatus", func(t *testing.T) {
			subs, err := subscriptionModel.GetSubscriptionStatus(ctx, s1)
			if err != nil {
				t.Errorf("GetSubscriptionStatus(): %v", err)
				return
			}
			wantNum := 3
			if len(subs) != wantNum {
				t.Errorf("GetSubscriptionStatus(): got %d categories, want %d", len(subs), wantNum)
			}
			for _, sub := range subs {
				if sub.Subscribed {
					t.Errorf("GetSubscriptionStatus(): %v: want unsubscribed", sub.Category.Name)
				}
				if sub.Unread > 0 {
					t.Errorf("GetSubscriptionStatus(): %v: got %d unreads, want 0", sub.Category.Name, sub.Unread)
				}
			}
		})

		t.Run("Subscribe", func(t *testing.T) {
			if err := subscriptionModel.Subscribe(ctx, s1, *cat1); err != nil {
				t.Errorf("Subscribe(%q, %q): %v", s1.UserID, cat1.Name, err)
			}
			if err := subscriptionModel.Subscribe(ctx, s1, *cat2); err != nil {
				t.Errorf("Subscribe(%q, %q): %v", s1.UserID, cat2.Name, err)
			}
			if err := subscriptionModel.Subscribe(ctx, s2, *cat1); err != nil {
				t.Errorf("Subscribe(%q, %q): %v", s2.UserID, cat1.Name, err)
			}
			if err := subscriptionModel.Subscribe(ctx, s2, *cat3); err != nil {
				t.Errorf("Subscribe(%q, %q): %v", s2.UserID, cat3.Name, err)
			}

			t.Run(s1.UserID, func(t *testing.T) {
				if len(s1.Categories) != 2 {
					t.Errorf("%q has %d categories, want 2", s1.UserID, len(s1.Categories))
				}
				subs, err := subscriptionModel.GetSubscriptionStatus(ctx, s1)
				if err != nil {
					t.Errorf("GetSubscriptionStatus(%q): %v", s1.UserID, err)
					return
				}
				for _, sub := range subs {
					if sub.Category.ID == cat1.ID && !sub.Subscribed {
						t.Errorf("GetSubscriptionStatus(%q): %v: want subscribed", s1.UserID, sub.Category.Name)
					}
					if sub.Category.ID == cat2.ID && !sub.Subscribed {
						t.Errorf("GetSubscriptionStatus(%q): %v: want subscribed", s1.UserID, sub.Category.Name)
					}
					if sub.Category.ID == cat3.ID && sub.Subscribed {
						t.Errorf("GetSubscriptionStatus(%q): %v: want unsubscribed", s1.UserID, sub.Category.Name)
					}
				}
			})

			t.Run(s2.UserID, func(t *testing.T) {
				if len(s2.Categories) != 2 {
					t.Errorf("%q has %d categories, want 2", s2.UserID, len(s2.Categories))
				}
				subs, err := subscriptionModel.GetSubscriptionStatus(ctx, s2)
				if err != nil {
					t.Errorf("GetSubscriptionStatus(%q): %v", s2.UserID, err)
					return
				}
				for _, sub := range subs {
					if sub.Category.ID == cat1.ID && !sub.Subscribed {
						t.Errorf("GetSubscriptionStatus(%q): %v: want subscribed", s2.UserID, sub.Category.Name)
					}
					if sub.Category.ID == cat2.ID && sub.Subscribed {
						t.Errorf("GetSubscriptionStatus(%q): %v: want unsubscribed", s2.UserID, sub.Category.Name)
					}
					if sub.Category.ID == cat3.ID && !sub.Subscribed {
						t.Errorf("GetSubscriptionStatus(%q): %v: want subscribed", s2.UserID, sub.Category.Name)
					}
				}
			})
		})

		t.Run("Unsubscribe", func(t *testing.T) {
			if err := subscriptionModel.Unsubscribe(ctx, s2, *cat3); err != nil {
				t.Errorf("Unsubscribe(%q, %q): %v", s2.UserID, cat3.Name, err)
			}

			t.Run(s2.UserID, func(t *testing.T) {
				if len(s2.Categories) != 1 {
					t.Errorf("%q has %d categories, want 1", s2.UserID, len(s2.Categories))
				}
				subs, err := subscriptionModel.GetSubscriptionStatus(ctx, s2)
				if err != nil {
					t.Errorf("GetSubscriptionStatus(%q): %v", s2.UserID, err)
					return
				}
				for _, sub := range subs {
					if sub.Category.ID == cat1.ID && !sub.Subscribed {
						t.Errorf("GetSubscriptionStatus(%q): %v: want subscribed", s2.UserID, sub.Category.Name)
					}
					if sub.Category.ID == cat2.ID && sub.Subscribed {
						t.Errorf("GetSubscriptionStatus(%q): %v: want unsubscribed", s2.UserID, sub.Category.Name)
					}
					if sub.Category.ID == cat3.ID && sub.Subscribed {
						t.Errorf("GetSubscriptionStatus(%q): %v: want unsubscribed", s2.UserID, sub.Category.Name)
					}
				}
			})
		})

		t.Run("AddUpdate", func(t *testing.T) {
			t.Run("no category", func(t *testing.T) {
				up := model.Update{Title: "No Cat"}
				if err := subscriptionModel.AddUpdate(ctx, up); err != model.ErrInvalidCategory {
					t.Errorf("AddUpdate(%q): want ErrInvalidCategory error, got %q", up.Title, err)
				}
			})
			up1 := model.Update{Category: cat1, Title: "Cat1 Up2"}
			if err := subscriptionModel.AddUpdate(ctx, up1); err != nil {
				t.Errorf("AddUpdate(%q): %v", up1.Title, err)
			}
			up2 := model.Update{Category: cat1, Title: "Cat1 Up3"}
			if err := subscriptionModel.AddUpdate(ctx, up2); err != nil {
				t.Errorf("AddUpdate(%q): %v", up2.Title, err)
			}
			up3 := model.Update{Category: cat2, Title: "Cat2 Up1"}
			if err := subscriptionModel.AddUpdate(ctx, up3); err != nil {
				t.Errorf("AddUpdate(%q): %v", up3.Title, err)
			}

			t.Run(s1.UserID, func(t *testing.T) {
				subs, err := subscriptionModel.GetSubscriptionStatus(ctx, s1)
				if err != nil {
					t.Errorf("GetSubscriptionStatus(%q): %v", s1.UserID, err)
					return
				}
				for _, sub := range subs {
					if sub.Category.ID == cat1.ID && sub.Unread != 2 {
						t.Errorf("GetSubscriptionStatus(%q): %v: got %d unreads, want %d", s1.UserID, sub.Category.Name, sub.Unread, 2)
					}
					if sub.Category.ID == cat2.ID && sub.Unread != 1 {
						t.Errorf("GetSubscriptionStatus(%q): %v: got %d unreads, want %d", s1.UserID, sub.Category.Name, sub.Unread, 1)
					}
					if sub.Category.ID == cat3.ID && sub.Unread != 0 {
						t.Errorf("GetSubscriptionStatus(%q): %v: got %d unreads, want %d", s1.UserID, sub.Category.Name, sub.Unread, 0)
					}
				}
			})

			t.Run(s2.UserID, func(t *testing.T) {
				subs, err := subscriptionModel.GetSubscriptionStatus(ctx, s2)
				if err != nil {
					t.Errorf("GetSubscriptionStatus(%q): %v", s2.UserID, err)
					return
				}
				for _, sub := range subs {
					if sub.Category.ID == cat1.ID && sub.Unread != 2 {
						t.Errorf("GetSubscriptionStatus(%q): %v: got %d unreads, want %d", s2.UserID, sub.Category.Name, sub.Unread, 2)
					}
					if sub.Category.ID == cat2.ID && sub.Unread != 0 {
						t.Errorf("GetSubscriptionStatus(%q): %v: got %d unreads, want %d", s2.UserID, sub.Category.Name, sub.Unread, 0)
					}
					if sub.Category.ID == cat3.ID && sub.Unread != 0 {
						t.Errorf("GetSubscriptionStatus(%q): %v: got %d unreads, want %d", s2.UserID, sub.Category.Name, sub.Unread, 0)
					}
				}
			})
		})
	})
}

func resetData(ctx context.Context, fso *firestorm.FSClient) {
	cats := make([]model.Category, 0)
	_ = fso.NewRequest().QueryEntities(ctx, fso.NewRequest().ToCollection(model.Category{}).Query, &cats)()
	if err := fso.NewRequest().DeleteEntities(ctx, cats)(); err != nil {
		log.Fatalf("resetData: failed to delete categories: %v", err)
	}

	sbs := make([]model.Subscriber, 0)
	_ = fso.NewRequest().QueryEntities(ctx, fso.NewRequest().ToCollection(model.Subscriber{}).Query, &sbs)()
	if err := fso.NewRequest().DeleteEntities(ctx, sbs)(); err != nil {
		log.Fatalf("resetData: failed to delete subscriptions: %v", err)
	}
}
