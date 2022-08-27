package main

import (
	"cloud.google.com/go/firestore"
	"context"
	firestoreDb "github.com/d-ashesss/news-feed-bot/pkg/db/firestore"
	"github.com/d-ashesss/news-feed-bot/pkg/feed/fetcher"
	"github.com/d-ashesss/news-feed-bot/pkg/model"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	ctx := context.Background()
	fsc, err := firestore.NewClient(ctx, projectID)
	defer func(fsc *firestore.Client) {
		_ = fsc.Close()
	}(fsc)
	if err != nil {
		log.Fatalf("failed to create firestore client: %v", err)
	}

	if len(os.Args) > 2 {
		log.Fatalf("Usage: fetch-updates [category-id]")
	}

	feedModel := firestoreDb.NewFeedModel(fsc)
	categoryModel := firestoreDb.NewCategoryModel(fsc)
	subscriberModel := firestoreDb.NewSubscriberModel(fsc, nil)
	updateModel := firestoreDb.NewUpdateModel(fsc)
	subscriptionModel := firestoreDb.NewSubscriptionModel(fsc, categoryModel, subscriberModel, updateModel)

	if len(os.Args) == 2 {
		catID := getCatId()
		cat, err := categoryModel.Get(ctx, catID)
		if err != nil {
			log.Fatalf("get category %q: %s", catID, err)
		}
		fetchCategory(ctx, feedModel, subscriptionModel, cat)
	} else {
		cats, err := categoryModel.GetAll(ctx)
		if err != nil {
			log.Fatalf("get categories: %v", err)
		}
		for _, cat := range cats {
			fetchCategory(ctx, feedModel, subscriptionModel, &cat)
		}
	}
}

func getCatId() string {
	return strings.TrimSpace(os.Args[1])
}

func fetchCategory(ctx context.Context, feedModel model.FeedModel, subscriptionModel model.SubscriptionModel, cat *model.Category) {
	feeds, err := feedModel.GetAll(ctx, cat)
	if err != nil {
		log.Printf("get feeds in %q: %s", cat.Name, err)
		return
	}

	f := fetcher.New(subscriptionModel)
	for _, feed := range feeds {
		if err := f.Fetch(ctx, &feed, cat); err != nil {
			log.Printf("fetch updates from %q: %v", feed.Title, err)
		}
		if err := feedModel.SetUpdated(ctx, &feed, time.Now().UTC()); err != nil {
			log.Printf("update feed %q: %v", feed.Title, err)
		}
	}
}
