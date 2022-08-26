package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	firestoreDb "github.com/d-ashesss/news-feed-bot/pkg/db/firestore"
	"github.com/d-ashesss/news-feed-bot/pkg/feed/fetcher"
	"github.com/d-ashesss/news-feed-bot/pkg/model"
	"log"
	"os"
	"strings"
	"time"
)

func init() {
	log.SetFlags(0)
}

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

	if len(os.Args) != 3 {
		log.Fatalf("Usage: add-feed <category-id> <feed-url>")
	}

	feedModel := firestoreDb.NewFeedModel(fsc)
	categoryModel := firestoreDb.NewCategoryModel(fsc)
	subscriberModel := firestoreDb.NewSubscriberModel(fsc, nil)
	updateModel := firestoreDb.NewUpdateModel(fsc)
	subscriptionModel := firestoreDb.NewSubscriptionModel(fsc, categoryModel, subscriberModel, updateModel)

	catID := getCatId(os.Args)
	if len(catID) == 0 {
		log.Fatalf("invalid category")
	}
	cat, err := categoryModel.Get(ctx, catID)
	if err != nil {
		log.Fatalf("get category %q: %s", catID, err)
	}

	url := getFeedUrl(os.Args)
	if len(url) == 0 {
		log.Fatalf("invalid feed URL")
	}
	f := fetcher.New(subscriptionModel)
	title, err := f.GetTitle(ctx, url)
	if err != nil {
		log.Fatalf("get feed title: %v", err)
	}

	feed := &model.Feed{
		Category:   cat,
		Title:      title,
		URL:        url,
		LastUpdate: time.Now().UTC(),
	}
	if _, err := feedModel.Create(ctx, feed); err != nil {
		log.Fatalf("create feed: %v", err)
	}
	fmt.Printf("%s %q [%s]: %s\n", feed.ID, feed.Title, feed.LastUpdate.Format(time.RFC822), feed.URL)
}

func getCatId(args []string) string {
	id := args[1]
	return strings.TrimSpace(id)
}

func getFeedUrl(args []string) string {
	url := args[2]
	return strings.TrimSpace(url)
}
