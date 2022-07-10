package main

import (
	"context"
	"fmt"
	"github.com/d-ashesss/news-feed-bot/pkg/model"
	"gopkg.in/tucnak/telebot.v2"
	"log"
)

// botHandleStartCmd handles /start command.
//   Shows welcome message.
func (a *App) botHandleStartCmd(ctx context.Context, m *telebot.Message) {
	user := ctx.Value(BotCtxUser).(*model.Subscriber)
	if _, err := a.Bot.Send(
		m.Sender,
		fmt.Sprintf("Welcome, *%s* 🎉", user.ID),
		&telebot.SendOptions{ParseMode: telebot.ModeMarkdown},
	); err != nil {
		log.Printf("[bot] botHandleStartCmd() Failed to reply: %v", err)
	}
	a.botHandleMenuCmd(ctx, m)
}

// botHandleStartCmd handles /menu command.
//   Shows main menu.
func (a *App) botHandleMenuCmd(ctx context.Context, m *telebot.Message) {
	user := ctx.Value(BotCtxUser).(*model.Subscriber)
	if _, err := a.Bot.Send(
		m.Sender,
		fmt.Sprintf("Your menu, *%s*: 🗒", user.ID),
		&telebot.SendOptions{ParseMode: telebot.ModeMarkdown},
		NewBotMenuMain().Menu,
	); err != nil {
		log.Printf("[bot] botHandleMenuCmd() Failed to reply: %v", err)
	}
}

// botHandleCheckUpdatesCallback handles request to show unread updates.
func (a *App) botHandleCheckUpdatesCallback(ctx context.Context, cb *telebot.Callback) {
	user := ctx.Value(BotCtxUser).(*model.Subscriber)

	subs, err := a.SubscriptionModel.GetSubscriptionStatus(ctx, user)
	if err != nil {
		log.Printf("botHandleCheckUpdatesCallback: subscription status: %v", err)
		_ = a.Bot.Respond(cb, &telebot.CallbackResponse{Text: ""})
		return
	}
	_, _ = a.Bot.Edit(
		cb.Message,
		fmt.Sprintf("Your updates"),
		&telebot.SendOptions{ParseMode: telebot.ModeMarkdown},
		NewBotMenuCategoriesUpdates(subs).Menu,
	)
	_ = a.Bot.Respond(cb, &telebot.CallbackResponse{Text: ""})
}

// botHandleStartCmd handles /delete command.
//   Provides user with a choise to delete his data from the service.
//   - Confirm action will be handled by botHandleDeleteConfirmCallback
//   - Cancel action will be handled by botHandleDeleteCancelCallback
func (a *App) botHandleDeleteCmd(ctx context.Context, m *telebot.Message) {
	user := ctx.Value(BotCtxUser).(*model.Subscriber)

	if _, err := a.Bot.Send(
		m.Sender,
		fmt.Sprintf("You data is about to be deleted from our service, *%s* ♻", user.ID),
		&telebot.SendOptions{ParseMode: telebot.ModeMarkdown},
		NewBotMenuDelete().Menu,
	); err != nil {
		log.Printf("[bot] botHandleDeleteCmd() Failed to reply: %v", err)
	}
	if err := a.Bot.Delete(m); err != nil {
		log.Printf("[bot] botHandleDeleteCmd() Failed delete user message: %v", err)
	}
}

// botHandleDeleteConfirmCallback handles confirmation callback of Delete User menu.
func (a *App) botHandleDeleteConfirmCallback(ctx context.Context, cb *telebot.Callback) {
	user := ctx.Value(BotCtxUser).(*model.Subscriber)

	if err := a.SubscriberModel.Delete(ctx, user); err != nil {
		log.Printf("[bot] botHandleDeleteConfirmCallback() Failed to delete user: %v", err)
		_ = a.Bot.Respond(cb, &telebot.CallbackResponse{Text: "Something went wrong, please try again later.", ShowAlert: true})
		return
	}

	if _, err := a.Bot.Edit(
		cb.Message,
		fmt.Sprintf("Your data was successfully deleted, *%s* 👍", user.ID),
		&telebot.SendOptions{ParseMode: telebot.ModeMarkdown},
	); err != nil {
		log.Printf("[bot] botHandleDeleteConfirmCallback() Failed to edit message: %v", err)
	}
	if _, err := a.Bot.Send(
		cb.Sender,
		"You can always come back later, if you want. See you!",
	); err != nil {
		log.Printf("[bot] botHandleDeleteConfirmCallback() Failed to reply: %v", err)
	}
	_ = a.Bot.Respond(cb, &telebot.CallbackResponse{Text: ""})
}

// botHandleDeleteCancelCallback handles cancellation callback of Delete User menu.
func (a *App) botHandleDeleteCancelCallback(ctx context.Context, cb *telebot.Callback) {
	user := ctx.Value(BotCtxUser).(*model.Subscriber)

	if _, err := a.Bot.Edit(
		cb.Message,
		fmt.Sprintf("You will not be deleted, *%s* 👍", user.ID),
		&telebot.SendOptions{ParseMode: telebot.ModeMarkdown},
	); err != nil {
		log.Printf("[bot] botHandleDeleteCancelCallback() Failed to edit message: %v", err)
	}
	_ = a.Bot.Respond(cb, &telebot.CallbackResponse{Text: ""})
}

// botHandleTextMessage is an arbitrary method to handle any text message that was not handled by a specific handler.
func (a *App) botHandleTextMessage(_ context.Context, _ *telebot.Message) {
}
