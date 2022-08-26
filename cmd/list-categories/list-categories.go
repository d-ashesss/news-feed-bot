package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	firestoreDb "github.com/d-ashesss/news-feed-bot/pkg/db/firestore"
	"log"
	"os"
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

	feedModel := firestoreDb.NewFeedModel(fsc)
	categoryModel := firestoreDb.NewCategoryModel(fsc)

	cats, err := categoryModel.GetAll(ctx)
	if err != nil {
		log.Fatalf("failed to get categories: %s", err)
	}

	for _, cat := range cats {
		fmt.Printf("%s: %s\n", cat.Name, cat.ID)
		feeds, err := feedModel.GetAll(ctx, &cat)
		if err != nil {
			log.Printf("get %s feeds: %v", cat.Name, err)
		}
		for _, feed := range feeds {
			fmt.Printf("\t%s %q [%s]: %s\n", feed.ID, feed.Title, feed.LastUpdate.Format(time.RFC822), feed.URL)
		}
	}
}
