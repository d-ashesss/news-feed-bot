package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	firestoreDb "github.com/d-ashesss/news-feed-bot/pkg/db/firestore"
	"github.com/d-ashesss/news-feed-bot/pkg/model"
	"log"
	"os"
	"strings"
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

	if len(os.Args) != 2 {
		log.Fatalf("Usage: create-category <category-name>")
	}

	categoryModel := firestoreDb.NewCategoryModel(fsc)
	catName := getCatName(os.Args)
	if len(catName) == 0 {
		log.Fatalf("invalid category name")

	}
	cat := model.NewCategory(catName)
	if _, err := categoryModel.Create(ctx, cat); err != nil {
		log.Fatalf("failed to create category: %s", err)
	}
	fmt.Println(cat.ID)
}

func getCatName(args []string) string {
	name := args[1]
	return strings.TrimSpace(name)
}
