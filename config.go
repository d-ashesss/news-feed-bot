package main

import (
	"context"
	"github.com/d-ashesss/news-feed-bot/secretmanager"
	"log"
	"os"
)

type Config struct {
	TelegramToken   string
	BaseURL         string
	WebPort         string
	BotWebhookMode  bool
	BotResetWebhook bool
}

func loadConfig(ctx context.Context, projectID string, secretManager *secretmanager.SecretManager) Config {
	var telegramToken string
	var ok bool
	telegramToken, ok = os.LookupEnv("TELEGRAM_TOKEN")
	if !ok && secretManager != nil {
		var err error
		if telegramToken, err = secretManager.GetSecret(ctx, "telegram-bot-token"); err != nil {
			log.Printf("[config] secretManager.GetSecret: %v", err)
		}
	}

	baseURL := os.Getenv("APP_BASE_URL")
	if len(baseURL) == 0 {
		if len(projectID) > 0 {
			baseURL = "https://" + projectID + ".appspot.com"
		} else {
			baseURL = "http://localhost"
		}
	}

	WebPort := os.Getenv("PORT")

	_, BotWebhookMode := os.LookupEnv("BOT_WEBHOOK_MODE")
	_, BotResetWebhook := os.LookupEnv("BOT_RESET_WEBHOOK")

	return Config{
		TelegramToken:   telegramToken,
		BaseURL:         baseURL,
		WebPort:         WebPort,
		BotWebhookMode:  BotWebhookMode,
		BotResetWebhook: BotResetWebhook,
	}
}
