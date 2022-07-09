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

func TestSubscriptionModel(t *testing.T) {
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
	subscriptionModel := firestoreDb.NewSubscriptionModel(fsc, categoryModel, subscriberModel, updateModel)

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

	const s1UserId = "U1"
	s1 := model.NewSubscriber(s1UserId)
	if _, err = subscriberModel.Create(ctx, s1); err != nil {
		t.Fatalf("failed to create subscriber %v: %v", s1, err)
	}
	const s2UserId = "U2"
	s2 := model.NewSubscriber(s2UserId)
	if _, err = subscriberModel.Create(ctx, s2); err != nil {
		t.Fatalf("failed to create subscriber %v: %v", s2, err)
	}

	const cat1up1Title = "Cat1 Up1"
	cat1up1 := model.Update{Category: cat1, Title: cat1up1Title, Date: time.Now().Add(-3 * time.Minute)}
	const cat1up2Title = "Cat1 Up2"
	cat1up2 := model.Update{Category: cat1, Title: cat1up2Title, Date: time.Now().Add(-2 * time.Minute)}
	const cat2up1Title = "Cat2 Up1"
	cat2up1 := model.Update{Category: cat2, Title: cat2up1Title, Date: time.Now().Add(-1 * time.Minute)}

	t.Run("SubscriptionModel", func(t *testing.T) {
		t.Run("GetSubscriptionStatus", func(t *testing.T) {
			t.Run("nil subscriber", func(t *testing.T) {
				var s *model.Subscriber
				if _, err := subscriptionModel.GetSubscriptionStatus(ctx, s); err != model.ErrInvalidSubscriber {
					t.Fatalf("GetSubscriptionStatus(%q): got %q; want ErrInvalidSubscriber", s1.UserID, err)
				}
			})

			t.Run("empty subscriber", func(t *testing.T) {
				s := &model.Subscriber{}
				subs, err := subscriptionModel.GetSubscriptionStatus(ctx, s)
				if err != nil {
					t.Fatalf("GetSubscriptionStatus(%v): %v", s, err)
				}
				assertGetSubscriptionStatusSubscribed(t, subs, map[string]bool{
					cat1.ID: false,
					cat2.ID: false,
					cat3.ID: false,
				})
				assertGetSubscriptionStatusUnread(t, subs, map[string]int{
					cat1.ID: 0,
					cat2.ID: 0,
					cat3.ID: 0,
				})
			})

			t.Run("valid subscriber", func(t *testing.T) {
				subs, err := subscriptionModel.GetSubscriptionStatus(ctx, s1)
				if err != nil {
					t.Fatalf("GetSubscriptionStatus(%q): %v", s1.UserID, err)
				}
				wantNum := 3
				if len(subs) != wantNum {
					t.Errorf("GetSubscriptionStatus(%q): got %d categories; want %d", s1.UserID, len(subs), wantNum)
				}
				assertGetSubscriptionStatusSubscribed(t, subs, map[string]bool{
					cat1.ID: false,
					cat2.ID: false,
					cat3.ID: false,
				})
				assertGetSubscriptionStatusUnread(t, subs, map[string]int{
					cat1.ID: 0,
					cat2.ID: 0,
					cat3.ID: 0,
				})
			})
		})

		t.Run("Subscribe", func(t *testing.T) {
			t.Run("nil subscriber", func(t *testing.T) {
				var s *model.Subscriber
				if err := subscriptionModel.Subscribe(ctx, s, *cat1); err != model.ErrInvalidSubscriber {
					t.Fatalf("Subscribe(%v, %q): got %q; want ErrInvalidSubscriber", s, cat1.Name, err)
				}
			})

			t.Run("empty category", func(t *testing.T) {
				cat := model.Category{}
				if err := subscriptionModel.Subscribe(ctx, s1, cat); err != model.ErrInvalidCategory {
					t.Fatalf("Subscribe(%q, %v): got %q; want ErrInvalidCategory", s1.UserID, cat, err)
				}
			})

			if err := subscriptionModel.Subscribe(ctx, s1, *cat1); err != nil {
				t.Fatalf("Subscribe(%q, %q): %v", s1.UserID, cat1.Name, err)
			}
			if err := subscriptionModel.Subscribe(ctx, s1, *cat2); err != nil {
				t.Fatalf("Subscribe(%q, %q): %v", s1.UserID, cat2.Name, err)
			}
			if err := subscriptionModel.Subscribe(ctx, s2, *cat1); err != nil {
				t.Fatalf("Subscribe(%q, %q): %v", s2.UserID, cat1.Name, err)
			}
			if err := subscriptionModel.Subscribe(ctx, s2, *cat3); err != nil {
				t.Fatalf("Subscribe(%q, %q): %v", s2.UserID, cat3.Name, err)
			}

			t.Run(s1UserId, func(t *testing.T) {
				wantNum := 2
				if len(s1.Categories) != wantNum {
					t.Errorf("%q has %d categories; want %d", s1.UserID, len(s1.Categories), wantNum)
				}
				subs, err := subscriptionModel.GetSubscriptionStatus(ctx, s1)
				if err != nil {
					t.Fatalf("GetSubscriptionStatus(%q): %v", s1.UserID, err)
				}
				assertGetSubscriptionStatusSubscribed(t, subs, map[string]bool{
					cat1.ID: true,
					cat2.ID: true,
					cat3.ID: false,
				})
			})

			t.Run(s2UserId, func(t *testing.T) {
				wantNum := 2
				if len(s2.Categories) != wantNum {
					t.Errorf("%q has %d categories; want %d", s2.UserID, len(s2.Categories), wantNum)
				}
				subs, err := subscriptionModel.GetSubscriptionStatus(ctx, s2)
				if err != nil {
					t.Fatalf("GetSubscriptionStatus(%q): %v", s2.UserID, err)
				}
				assertGetSubscriptionStatusSubscribed(t, subs, map[string]bool{
					cat1.ID: true,
					cat2.ID: false,
					cat3.ID: true,
				})
			})
		})

		t.Run("Unsubscribe", func(t *testing.T) {
			t.Run("nil subscriber", func(t *testing.T) {
				var s *model.Subscriber
				if err := subscriptionModel.Unsubscribe(ctx, s, *cat1); err != model.ErrInvalidSubscriber {
					t.Fatalf("Unsubscribe(%v, %q): got %q; want ErrInvalidSubscriber", s, cat1.Name, err)
				}
			})

			t.Run("empty category", func(t *testing.T) {
				cat := &model.Category{}
				if err := subscriptionModel.Unsubscribe(ctx, s1, *cat); err != model.ErrInvalidCategory {
					t.Fatalf("Unsubscribe(%q, %v): got %q; want ErrInvalidCategory", s1.UserID, cat, err)
				}
			})

			if err := subscriptionModel.Unsubscribe(ctx, s2, *cat3); err != nil {
				t.Fatalf("Unsubscribe(%q, %q): %v", s2.UserID, cat3.Name, err)
			}

			t.Run(s2UserId, func(t *testing.T) {
				wantNum := 1
				if len(s2.Categories) != wantNum {
					t.Errorf("%q has %d categories; want %d", s2.UserID, len(s2.Categories), wantNum)
				}
				subs, err := subscriptionModel.GetSubscriptionStatus(ctx, s2)
				if err != nil {
					t.Fatalf("GetSubscriptionStatus(%q): %v", s2.UserID, err)
				}
				assertGetSubscriptionStatusSubscribed(t, subs, map[string]bool{
					cat1.ID: true,
					cat2.ID: false,
					cat3.ID: false,
				})
			})
		})

		t.Run("AddUpdate", func(t *testing.T) {
			t.Run("no category", func(t *testing.T) {
				up := model.Update{Title: "No Cat"}
				if err := subscriptionModel.AddUpdate(ctx, up); err != model.ErrInvalidCategory {
					t.Errorf("AddUpdate(%q): got %q; want ErrInvalidCategory", up.Title, err)
				}
			})

			if err := subscriptionModel.AddUpdate(ctx, cat1up1); err != nil {
				t.Fatalf("AddUpdate(%q): %v", cat1up1.Title, err)
			}
			if err := subscriptionModel.AddUpdate(ctx, cat1up2); err != nil {
				t.Fatalf("AddUpdate(%q): %v", cat1up2.Title, err)
			}
			if err := subscriptionModel.AddUpdate(ctx, cat2up1); err != nil {
				t.Fatalf("AddUpdate(%q): %v", cat2up1.Title, err)
			}

			t.Run(s1UserId, func(t *testing.T) {
				subs, err := subscriptionModel.GetSubscriptionStatus(ctx, s1)
				if err != nil {
					t.Fatalf("GetSubscriptionStatus(%q): %v", s1.UserID, err)
				}
				assertGetSubscriptionStatusUnread(t, subs, map[string]int{
					cat1.ID: 2,
					cat2.ID: 1,
					cat3.ID: 0,
				})
			})

			t.Run(s2UserId, func(t *testing.T) {
				subs, err := subscriptionModel.GetSubscriptionStatus(ctx, s2)
				if err != nil {
					t.Fatalf("GetSubscriptionStatus(%q): %v", s2.UserID, err)
				}
				assertGetSubscriptionStatusUnread(t, subs, map[string]int{
					cat1.ID: 2,
					cat2.ID: 0,
					cat3.ID: 0,
				})
			})
		})

		t.Run("ShiftUpdate", func(t *testing.T) {
			t.Run("nil subscriber", func(t *testing.T) {
				var s *model.Subscriber
				_, err := subscriptionModel.ShiftUpdate(ctx, s, *cat1)
				if err != model.ErrInvalidSubscriber {
					t.Fatalf("ShiftUpdate(%v, %q): got %q; want ErrInvalidSubscriber", s, cat1.Name, err)
				}
			})

			t.Run("empty category", func(t *testing.T) {
				cat := &model.Category{}
				_, err := subscriptionModel.ShiftUpdate(ctx, s1, *cat)
				if err != model.ErrInvalidCategory {
					t.Fatalf("ShiftUpdate(%q, %v): got %q; want ErrInvalidCategory", s1.UserID, cat, err)
				}
			})

			t.Run(cat1up1Title, func(t *testing.T) {
				up, err := subscriptionModel.ShiftUpdate(ctx, s1, *cat1)
				if err != nil {
					t.Fatalf("ShiftUpdate(%q, %q): %v", s1.UserID, cat1.Name, err)
				}
				if up.Title != cat1up1.Title {
					t.Errorf("ShiftUpdate(%q, %q): = %q; want %q", s1.UserID, cat1.Name, up.Title, cat1up1.Title)
				}
			})
			t.Run(cat1up2Title, func(t *testing.T) {
				up, err := subscriptionModel.ShiftUpdate(ctx, s1, *cat1)
				if err != nil {
					t.Fatalf("ShiftUpdate(%q, %q): %v", s1.UserID, cat1.Name, err)
				}
				if up.Title != cat1up2.Title {
					t.Errorf("ShiftUpdate(%q, %q): = %q; want %q", s1.UserID, cat1.Name, up.Title, cat1up2.Title)
				}
			})
			t.Run("no update left", func(t *testing.T) {
				_, err := subscriptionModel.ShiftUpdate(ctx, s1, *cat1)
				if err != model.ErrNoUpdates {
					t.Fatalf("ShiftUpdate(%q, %q): got %q; want ErrNoUpdates", s1.UserID, cat1.Name, err)
				}
			})
			t.Run(cat2up1Title, func(t *testing.T) {
				up, err := subscriptionModel.ShiftUpdate(ctx, s1, *cat2)
				if err != nil {
					t.Fatalf("ShiftUpdate(%q, %q): %v", s1.UserID, cat2.Name, err)
				}
				if up.Title != cat2up1.Title {
					t.Errorf("ShiftUpdate(%q, %q): = %q; want %q", s1.UserID, cat2.Name, up.Title, cat2up1.Title)
				}
			})
		})
	})
}

func assertGetSubscriptionStatusSubscribed(t *testing.T, subs []model.Subscription, want map[string]bool) {
	t.Helper()
	for _, sub := range subs {
		wantStatus, ok := want[sub.Category.ID]
		if !ok {
			t.Errorf("GetSubscriptionStatus(%q): unexpected category %v", t.Name(), sub.Category)
			continue
		}
		if sub.Subscribed != wantStatus {
			t.Errorf("GetSubscriptionStatus(%q): %v: got Subscribed = %v; want %v", t.Name(), sub.Category.Name, sub.Subscribed, wantStatus)
		}
	}
}

func assertGetSubscriptionStatusUnread(t *testing.T, subs []model.Subscription, want map[string]int) {
	t.Helper()
	for _, sub := range subs {
		wantUnread, ok := want[sub.Category.ID]
		if !ok {
			t.Errorf("GetSubscriptionStatus(%q): unexpected category %v", t.Name(), sub.Category)
			continue
		}
		if sub.Unread != wantUnread {
			t.Errorf("GetSubscriptionStatus(%q): %v: got Unread = %v; want %v", t.Name(), sub.Category.Name, sub.Unread, wantUnread)
		}
	}
}
