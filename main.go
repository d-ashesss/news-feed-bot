package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/d-ashesss/news-feed-bot/bot"
	"github.com/d-ashesss/news-feed-bot/http"
	firestoreDb "github.com/d-ashesss/news-feed-bot/pkg/db/firestore"
	memoryDb "github.com/d-ashesss/news-feed-bot/pkg/db/memory"
	"github.com/d-ashesss/news-feed-bot/pkg/model"
	"github.com/d-ashesss/news-feed-bot/secretmanager"
	"log"
	"os"
)

var projectID string

func init() {
	projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")

	log.SetOutput(os.Stdout)
	log.SetFlags(0)
}

func main() {
	ctx := context.Background()
	secretManager, err := secretmanager.New(ctx, projectID)
	if err != nil {
		log.Printf("[main] Failed to init secret manager: %v", err)
		log.Printf("[main] Only local configuration will be used")
	} else {
		defer secretManager.Close()
	}

	config := loadConfig(ctx, projectID, secretManager)

	var userStore model.UserStore
	fstore, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Printf("[main] Failed to init Firestore: %v", err)
		userStore = memoryDb.NewUserStore()
	} else {
		defer func() { _ = fstore.Close() }()
		userStore = firestoreDb.NewUserStore(fstore)
	}

	httpServer := http.NewServer(config.WebPort)

	userModel := model.NewUserModel(userStore)

	app := NewApp(config, httpServer, userModel)

	b, err := bot.New(config.TelegramToken)
	if err != nil {
		log.Printf("[main] Failed to init TG bot: %v", err)
	} else {
		if err = app.SetBot(b); err != nil {
			log.Printf("[main] Failed to integrate TG bot with the app: %v", err)
		}
	}

	app.Run()
}
