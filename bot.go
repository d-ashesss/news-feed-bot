package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/d-ashesss/news-feed-bot/bot"
	"github.com/d-ashesss/news-feed-bot/pkg/model"
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
	a.Bot.Handle("/start", a.botHandleMessage(botCtx, a.botHandleStart))
	a.Bot.Handle("/menu", a.botHandleMessage(botCtx, a.botHandleMenu))
	a.Bot.Handle("/delete", a.botHandleMessage(botCtx, a.botHandleDelete))
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

// botHandleMessage initializes common middleware stack to handle TG message
func (a *App) botHandleMessage(ctx context.Context, h func(ctx context.Context, m *telebot.Message)) func(m *telebot.Message) {
	return botMessageHandlerWithContext(
		ctx,
		h,
		a.botMiddlewareMessageGetUser,
		a.botMiddlewareMessageLogMessage,
	)
}

// botMessageHandlerWithContext adapts telebot handler to the form `func(ctx context.Context, m *telebot.Message)` allowing to pass context.
// middleware in stack will be executed in LIFO order
func botMessageHandlerWithContext(
	ctx context.Context,
	handler func(ctx context.Context, m *telebot.Message),
	stack ...func(func(ctx context.Context, m *telebot.Message)) func(ctx context.Context, m *telebot.Message),
) func(m *telebot.Message) {
	next := handler
	for _, fn := range stack {
		next = fn(next)
	}
	return func(m *telebot.Message) {
		next(ctx, m)
	}
}

// botMiddlewareMessageLogMessage logs incoming message
func (a *App) botMiddlewareMessageLogMessage(next func(ctx context.Context, m *telebot.Message)) func(ctx context.Context, m *telebot.Message) {
	return func(ctx context.Context, m *telebot.Message) {
		log.Printf("[bot] Incoiming message from %s: %q", bot.GetUserName(m.Sender), m.Text)

		next(ctx, m)
	}
}

// botMiddlewareMessageGetUser loads existing model.User or creating a new one
func (a *App) botMiddlewareMessageGetUser(next func(ctx context.Context, m *telebot.Message)) func(ctx context.Context, m *telebot.Message) {
	return func(ctx context.Context, m *telebot.Message) {
		user, err := a.UserModel.GetByTelegramID(ctx, m.Sender.ID)
		if err != nil {
			log.Printf("[bot] No user for TG ID %d: %v", m.Sender.ID, err)
			user = model.NewUser(m.Sender.ID)
			if err := a.UserModel.Create(ctx, user); err != nil {
				log.Printf("[bot] Failed to create user: %q", err)
				return
			} else {
				log.Printf("[bot] Created user %q", user.ID)
			}
		} else {
			log.Printf("[bot] Found user %q with TG ID %d", user.ID, user.TelegramID)
		}

		ctx = context.WithValue(ctx, BotCtxUser, user)
		next(ctx, m)
	}
}

func (a *App) botHandleStart(ctx context.Context, m *telebot.Message) {
	user := ctx.Value(BotCtxUser).(*model.User)
	if _, err := a.Bot.Send(
		m.Sender,
		fmt.Sprintf("Welcome, *%s* 🎉", user.ID),
		&telebot.SendOptions{ParseMode: telebot.ModeMarkdown},
	); err != nil {
		log.Printf("[bot] botHandleStart() Failed to reply: %v", err)
	}
}

func (a *App) botHandleMenu(ctx context.Context, m *telebot.Message) {
	user := ctx.Value(BotCtxUser).(*model.User)
	if _, err := a.Bot.Send(
		m.Sender,
		fmt.Sprintf("Your menu, *%s*: 🗒", user.ID),
		&telebot.SendOptions{ParseMode: telebot.ModeMarkdown},
	); err != nil {
		log.Printf("[bot] botHandleMenu() Failed to reply: %v", err)
	}
}

func (a *App) botHandleDelete(ctx context.Context, m *telebot.Message) {
	user := ctx.Value(BotCtxUser).(*model.User)
	if _, err := a.Bot.Send(
		m.Sender,
		fmt.Sprintf("You are going to be deleted, *%s* ♻️", user.ID),
		&telebot.SendOptions{ParseMode: telebot.ModeMarkdown},
	); err != nil {
		log.Printf("[bot] botHandleDelete() Failed to reply: %v", err)
	}
	if err := a.UserModel.Delete(ctx, user); err != nil {
		log.Printf("[bot] botHandleDelete() Failed to delete user: %v", err)
		return
	}
	if _, err := a.Bot.Send(
		m.Sender,
		fmt.Sprintf("Done 👍️"),
	); err != nil {
		log.Printf("[bot] botHandleDelete() Failed to reply: %v", err)
	}
}

func (a *App) botHandleTextMessage(_ context.Context, _ *telebot.Message) {
}
