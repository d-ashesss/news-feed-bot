//go:build integration
// +build integration

package model

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/d-ashesss/news-feed-bot/pkg/model"
	"github.com/jschoedt/go-firestorm"
	"testing"
)

func resetData(t *testing.T, ctx context.Context, c *firestore.Client) {
	t.Helper()

	fsc := firestorm.New(c, "ID", "")

	cats := make([]model.Category, 0)
	_ = fsc.NewRequest().QueryEntities(ctx, fsc.NewRequest().ToCollection(model.Category{}).Query, &cats)()
	if err := fsc.NewRequest().DeleteEntities(ctx, cats)(); err != nil {
		t.Fatalf("resetData: failed to delete categories: %v", err)
	}

	sbs := make([]model.Subscriber, 0)
	_ = fsc.NewRequest().QueryEntities(ctx, fsc.NewRequest().ToCollection(model.Subscriber{}).Query, &sbs)()
	if err := fsc.NewRequest().DeleteEntities(ctx, sbs)(); err != nil {
		t.Fatalf("resetData: failed to delete subscriptions: %v", err)
	}
}
