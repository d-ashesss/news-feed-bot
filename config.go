package main

import (
	"context"
	"github.com/d-ashesss/news-feed-bot/secretmanager"
	"os"
)

type Config struct {
	TelegramToken string
	BaseURL       string
	WebPort       string
}

func loadConfig(ctx context.Context, projectID string, secretManager *secretmanager.SecretManager) Config {
	var telegramToken string
	var ok bool
	telegramToken, ok = os.LookupEnv("TELEGRAM_TOKEN")
	if !ok && secretManager != nil {
		telegramToken, _ = secretManager.GetSecret(ctx, "telegram-bot-token")
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

	return Config{
		TelegramToken: telegramToken,
		BaseURL:       baseURL,
		WebPort:       WebPort,
	}
}
