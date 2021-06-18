package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/d-ashesss/news-feed-bot/bot"
	"github.com/google/uuid"
	"gopkg.in/tucnak/telebot.v2"
	"log"
	"net/http"
	"net/url"
	"os"
)

var (
	projectID, baseURL string
)

func init() {
	projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
	baseURL = os.Getenv("APP_BASE_URL")
	if len(baseURL) == 0 {
		if len(projectID) > 0 {
			baseURL = "https://" + projectID + ".appspot.com"
		} else {
			baseURL = "http://localhost"
		}
	}
}

func (a *App) SetBot(bot *bot.Bot) error {
	if bot == nil {
		return errors.New("invalid bot instance")
	}
	p, err := getBotWebhookPath(bot)
	if err != nil {
		p, err = createBotWebhookPath(bot)
		if err != nil {
			return errors.New(fmt.Sprintf("unable to create webhook: %v", err))
		}
	}
	a.Bot = bot
	a.HttpServer.Post(p, a.handleBotUpdate)
	a.Bot.Handle(telebot.OnText, a.handleBotMessage)
	return nil
}

func getBotWebhookPath(bot *bot.Bot) (string, error) {
	u, err := bot.WebhookURL()
	if err != nil {
		return "", err
	}
	parts, err := url.Parse(u)
	if err != nil {
		return "", err
	}
	return parts.Path, nil
}

func createBotWebhookPath(bot *bot.Bot) (string, error) {
	p := "/update/" + uuid.NewString()
	if err := bot.SetWebhookURL(baseURL + p); err != nil {
		return "", errors.New(fmt.Sprintf("unable to set webhook: %v", err))
	}
	log.Printf("Created new TG webhook %q", baseURL+p)
	return p, nil
}

func (a *App) handleBotUpdate(res http.ResponseWriter, r *http.Request) {
	if a.Bot == nil {
		res.WriteHeader(500)
		return
	}
	var update telebot.Update
	err := json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		log.Printf("[bot] Cannot decode update: %v", err)
		res.WriteHeader(500)
		return
	}
	a.Bot.ProcessUpdate(update)
}

func (a *App) handleBotMessage(m *telebot.Message) {
	log.Printf("[bot] Incoiming message from %s: %q", bot.GetUserName(m.Sender), m.Text)
	if m.Text == "/start" {
		if _, err := a.Bot.Send(m.Sender, "Welcome 🎉"); err != nil {
			log.Printf("[bot] Failed to reply: %v", err)
		}
	}
}
