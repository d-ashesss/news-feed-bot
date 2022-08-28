package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/d-ashesss/news-feed-bot/bot"
	"github.com/google/uuid"
	"gopkg.in/tucnak/telebot.v2"
	"log"
	"net/http"
	"net/url"
)

const (
	BotCtxUser = "user"
)

func (a *App) SetBot(bot *bot.Bot) error {
	if bot == nil {
		return errors.New("invalid bot instance")
	}
	if a.Config.BotResetWebhook {
		if err := bot.RemoveWebhook(); err != nil {
			return fmt.Errorf("unable to remove webhook: %v", err)
		}
	}
	if a.Config.BotWebhookMode {
		p, err := getBotWebhookPath(bot)
		if err != nil {
			p, err = createBotWebhookPath(a.Config.BaseURL, bot)
			if err != nil {
				return fmt.Errorf("unable to create webhook: %v", err)
			}
		}
		a.HttpServer.Post(p, a.botHandleWebhookUpdate)
	} else {
		if err := bot.RemoveWebhook(); err != nil {
			return fmt.Errorf("unable to remove webhook: %v", err)
		}
	}
	a.Bot = bot
	botCtx := context.TODO()
	a.Bot.Handle(telebot.OnText, a.botHandleMessage(botCtx, a.botHandleTextMessage))
	a.Bot.Handle("/start", a.botHandleMessage(botCtx, a.botHandleStartCmd))
	a.Bot.Handle("/menu", a.botHandleMessage(botCtx, a.botHandleMenuCmd))
	a.Bot.Handle("/delete", a.botHandleMessage(botCtx, a.botHandleDeleteCmd))

	a.Bot.Handle(&telebot.Btn{Unique: BotBtnBackToMainMenuID}, a.botHandleCallback(botCtx, a.botHandleBackToMainMenuCallback))

	menuMain := NewBotMenuMain()
	a.Bot.Handle(&menuMain.BtnCheckUpdates, a.botHandleCallback(botCtx, a.botHandleCheckUpdatesCallback))
	a.Bot.Handle(&menuMain.BtnSelectCategories, a.botHandleCallback(botCtx, a.botHandleSelectCategoriesCallback))

	a.Bot.Handle(&telebot.Btn{Unique: BotMenuSelectCategoriesBtnToggleCategoryID}, a.botHandleCallback(botCtx, a.botHandleToggleCategoryCallback))
	a.Bot.Handle(&telebot.Btn{Unique: BotMenuCategoryUpdatesBtnCategoryUpdatesID}, a.botHandleCallback(botCtx, a.botHandleCategoryUpdatesCallback))

	menuDelete := NewBotMenuDelete()
	a.Bot.Handle(&menuDelete.BtnConfirm, a.botHandleCallback(botCtx, a.botHandleDeleteConfirmCallback))
	a.Bot.Handle(&menuDelete.BtnCancel, a.botHandleCallback(botCtx, a.botHandleDeleteCancelCallback))
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

func createBotWebhookPath(baseURL string, bot *bot.Bot) (string, error) {
	p := "/update/" + uuid.NewString()
	if err := bot.SetWebhookURL(baseURL + p); err != nil {
		return "", fmt.Errorf("unable to set webhook: %v", err)
	}
	log.Printf("[bot] Created new TG webhook %q", baseURL+p)
	return p, nil
}

func (a *App) botHandleWebhookUpdate(res http.ResponseWriter, r *http.Request) {
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
