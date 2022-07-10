package main

import (
	"context"
	"github.com/d-ashesss/news-feed-bot/bot"
	"github.com/d-ashesss/news-feed-bot/pkg/model"
	"gopkg.in/tucnak/telebot.v2"
	"log"
)

// botHandleMessage initializes common middleware stack to handle TG message.
func (a *App) botHandleMessage(ctx context.Context, h func(ctx context.Context, m *telebot.Message)) func(m *telebot.Message) {
	return botMessageHandlerWithContext(
		ctx,
		h,
		a.botMiddlewareMessageGetUser,
		a.botMiddlewareMessageLogMessage,
	)
}

// botMessageHandlerWithContext adapts telebot handler to the form `func(ctx context.Context, m *telebot.Message)`
// allowing to pass context.
//   Middleware in stack will be executed in LIFO order.
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

// botMiddlewareMessageLogMessage logs incoming message.
func (a *App) botMiddlewareMessageLogMessage(next func(ctx context.Context, m *telebot.Message)) func(ctx context.Context, m *telebot.Message) {
	return func(ctx context.Context, m *telebot.Message) {
		log.Printf("[bot] Incoiming message from %s: %q", bot.GetUserName(m.Sender), m.Text)

		next(ctx, m)
	}
}

// botMiddlewareMessageGetUser loads existing model.User or creating a new one.
func (a *App) botMiddlewareMessageGetUser(next func(ctx context.Context, m *telebot.Message)) func(ctx context.Context, m *telebot.Message) {
	return func(ctx context.Context, m *telebot.Message) {
		ctx, err := loadUser(ctx, a.UserModel, m.Sender)
		if err != nil {
			log.Printf("[bot] Failed to load user: %q", err)
			return
		}
		next(ctx, m)
	}
}

// botHandleCallback initializes common middleware stack to handle TG message.
func (a *App) botHandleCallback(ctx context.Context, h func(ctx context.Context, cb *telebot.Callback)) func(cb *telebot.Callback) {
	return botCallbackHandlerWithContext(
		ctx,
		h,
		a.botMiddlewareCallbackGetUser,
	)
}

// botCallbackHandlerWithContext adapts telebot handler to the form `func(ctx context.Context, cb *telebot.Callback)`
// allowing to pass context.
//   Middleware in stack will be executed in LIFO order.
func botCallbackHandlerWithContext(
	ctx context.Context,
	handler func(ctx context.Context, cb *telebot.Callback),
	stack ...func(func(ctx context.Context, cb *telebot.Callback)) func(ctx context.Context, cb *telebot.Callback),
) func(cb *telebot.Callback) {
	next := handler
	for _, fn := range stack {
		next = fn(next)
	}
	return func(cb *telebot.Callback) {
		next(ctx, cb)
	}
}

// botMiddlewareCallbackGetUser loads existing model.User or creating a new one.
func (a *App) botMiddlewareCallbackGetUser(next func(ctx context.Context, cb *telebot.Callback)) func(ctx context.Context, cb *telebot.Callback) {
	return func(ctx context.Context, cb *telebot.Callback) {
		ctx, err := loadUser(ctx, a.UserModel, cb.Sender)
		if err != nil {
			log.Printf("[bot] Failed to load user: %q", err)
			return
		}
		next(ctx, cb)
	}
}

func loadUser(ctx context.Context, m *model.UserModel, s *telebot.User) (context.Context, error) {
	user, err := m.GetByTelegramID(ctx, s.ID)
	if err != nil {
		log.Printf("[bot] No user for TG ID %d", s.ID)
		user = model.NewUser(s.ID)
		if err := m.Create(ctx, user); err != nil {
			return ctx, err
		} else {
			log.Printf("[bot] Created user %q", user.ID)
		}
	} else {
		log.Printf("[bot] Found user %q with TG ID %d", user.ID, user.TelegramID)
	}

	return context.WithValue(ctx, BotCtxUser, user), nil
}
