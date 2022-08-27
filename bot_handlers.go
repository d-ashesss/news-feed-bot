package main

import (
	"context"
	"fmt"
	"github.com/d-ashesss/news-feed-bot/pkg/model"
	"gopkg.in/tucnak/telebot.v2"
	"log"
)

// botHandleStartCmd handles /start command.
//
//	Shows welcome message.
func (a *App) botHandleStartCmd(ctx context.Context, m *telebot.Message) {
	if _, err := a.Bot.Send(
		m.Sender,
		fmt.Sprintf("Welcome to this humble news bot!\n"+
			"Here you can receive news updates from the most famous world news agencies in the categories that you choose for yourself!\n"+
			"Please check out the menu to select the categories and start receiving the updates."),
		&telebot.SendOptions{ParseMode: telebot.ModeMarkdown},
	); err != nil {
		log.Printf("[bot] botHandleStartCmd() Failed to reply: %v", err)
	}
	a.botHandleMenuCmd(ctx, m)
}

// botHandleMenuCmd handles /menu command.
//
//	Shows main menu.
func (a *App) botHandleMenuCmd(_ context.Context, m *telebot.Message) {
	if _, err := a.Bot.Send(
		m.Sender,
		fmt.Sprintf("Please select the desired action:"),
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
		log.Printf("[bot] botHandleCheckUpdatesCallback(): subscription status: %v", err)
		return
	}
	selectedSubs := make([]model.Subscription, 0, len(subs))
	for _, sub := range subs {
		if sub.Subscribed {
			selectedSubs = append(selectedSubs, sub)
		}
	}
	if len(selectedSubs) == 0 {
		if _, err := a.Bot.Edit(
			cb.Message,
			fmt.Sprintf("You don't have any categories selected"),
			&telebot.SendOptions{ParseMode: telebot.ModeMarkdown},
			NewBotMenuNoCategoriesSelected().Menu,
		); err != nil {
			log.Printf("[bot] botHandleCheckUpdatesCallback(): Failed to edit message: %v", err)
		}
	} else {
		if _, err := a.Bot.Edit(
			cb.Message,
			fmt.Sprintf("Unread updates in categories you've selected:"),
			&telebot.SendOptions{ParseMode: telebot.ModeMarkdown},
			NewBotMenuCategoryUpdates(selectedSubs).Menu,
		); err != nil {
			log.Printf("[bot] botHandleCheckUpdatesCallback(): Failed to edit message: %v", err)
		}
	}
}

// botHandleCheckUpdatesCallback handles request to show the list of categories available for subscription.
func (a *App) botHandleSelectCategoriesCallback(ctx context.Context, cb *telebot.Callback) {
	user := ctx.Value(BotCtxUser).(*model.Subscriber)

	subs, err := a.SubscriptionModel.GetSubscriptionStatus(ctx, user)
	if err != nil {
		log.Printf("[bot] botHandleSelectCategoriesCallback(): subscription status: %v", err)
		return
	}
	if len(subs) == 0 {
		if _, err := a.Bot.Edit(
			cb.Message,
			fmt.Sprintf("Unfortunately I do not have any categories available at the moment, please come back later."),
			&telebot.SendOptions{ParseMode: telebot.ModeMarkdown},
		); err != nil {
			log.Printf("[bot] botHandleSelectCategoriesCallback(): Failed to edit message: %v", err)
		}
		return
	}
	if _, err := a.Bot.Edit(
		cb.Message,
		fmt.Sprintf("Select categories for which you would like to recieve updates:"),
		&telebot.SendOptions{ParseMode: telebot.ModeMarkdown},
		NewBotMenuSelectCategories(subs).Menu,
	); err != nil {
		log.Printf("[bot] botHandleSelectCategoriesCallback(): Failed to edit message: %v", err)
	}
}

// botHandleToggleCategoryCallback toggles selection of a category.
func (a *App) botHandleToggleCategoryCallback(ctx context.Context, cb *telebot.Callback) {
	user := ctx.Value(BotCtxUser).(*model.Subscriber)

	cat, err := a.CategoryModel.Get(ctx, cb.Data)
	if err != nil {
		log.Printf("[bot] botHandleToggleCategoryCallback(): get category: %v", err)
		return
	}
	if user.HasCategory(*cat) {
		if err := a.SubscriptionModel.Unsubscribe(ctx, user, *cat); err != nil {
			log.Printf("[bot] botHandleToggleCategoryCallback(): unsubscribe: %v", err)
			return
		}
	} else {
		if err := a.SubscriptionModel.Subscribe(ctx, user, *cat); err != nil {
			log.Printf("[bot] botHandleToggleCategoryCallback(): subscribe: %v", err)
			return
		}
	}

	subs, err := a.SubscriptionModel.GetSubscriptionStatus(ctx, user)
	if err != nil {
		log.Printf("[bot] botHandleToggleCategoryCallback(): subscription status: %v", err)
		return
	}
	if _, err := a.Bot.Edit(
		cb.Message,
		fmt.Sprintf("Select categories for which you would like to recieve updates:"),
		&telebot.SendOptions{ParseMode: telebot.ModeMarkdown},
		NewBotMenuSelectCategories(subs).Menu,
	); err != nil {
		log.Printf("[bot] botHandleToggleCategoryCallback(): Failed to edit message: %v", err)
	}
}

// botHandleCategoryUpdatesCallback shows the oldest update from selected category.
func (a *App) botHandleCategoryUpdatesCallback(ctx context.Context, cb *telebot.Callback) {
	user := ctx.Value(BotCtxUser).(*model.Subscriber)

	cat, err := a.CategoryModel.Get(ctx, cb.Data)
	if err != nil {
		log.Printf("[bot] botHandleCategoryUpdatesCallback(): get category: %v", err)
		return
	}

	up, err := a.SubscriptionModel.ShiftUpdate(ctx, user, *cat)
	if err == model.ErrNoUpdates {
		if _, err := a.Bot.Edit(
			cb.Message,
			fmt.Sprintf("You don't have any updates available in category *%v*", cat.Name),
			&telebot.SendOptions{ParseMode: telebot.ModeMarkdown},
			NewBotMenuNoUpdatesInCategory().Menu,
		); err != nil {
			log.Printf("[bot] botHandleCategoryUpdatesCallback(): Failed to edit message: %v", err)
		}
		return
	}
	if err != nil {
		log.Printf("[bot] botHandleCategoryUpdatesCallback(): shift update: %v", err)
		return
	}

	if _, err := a.Bot.Edit(
		cb.Message,
		up.FormatMessage(),
		&telebot.SendOptions{ParseMode: telebot.ModeMarkdown},
		NewBotMenuCategoryNextUpdate(cat, up.URL).Menu,
	); err != nil {
		log.Printf("[bot] botHandleCategoryUpdatesCallback(): Failed to show update: %v", err)
	}
}

// botHandleNextUpdateCallback removes the menu from previous update
// and shows the oldest update from selected category.
func (a *App) botHandleNextUpdateCallback(ctx context.Context, cb *telebot.Callback) {
	if _, err := a.Bot.Edit(cb.Message, NewUpdatedBotMenuUpdateURL(cb.Message.ReplyMarkup).Menu); err != nil {
		log.Printf("[bot] botHandleNextUpdateCallback(): Failed to remove menu from update: %v", err)
	}

	user := ctx.Value(BotCtxUser).(*model.Subscriber)

	cat, err := a.CategoryModel.Get(ctx, cb.Data)
	if err != nil {
		log.Printf("[bot] botHandleCategoryUpdatesCallback(): get category: %v", err)
		return
	}

	up, err := a.SubscriptionModel.ShiftUpdate(ctx, user, *cat)
	if err == model.ErrNoUpdates {
		if _, err := a.Bot.Send(
			cb.Sender,
			fmt.Sprintf("You don't have any updates available in category *%v*", cat.Name),
			&telebot.SendOptions{ParseMode: telebot.ModeMarkdown},
			NewBotMenuNoUpdatesInCategory().Menu,
		); err != nil {
			log.Printf("[bot] botHandleCategoryUpdatesCallback(): Failed to edit message: %v", err)
		}
		return
	}
	if err != nil {
		log.Printf("[bot] botHandleCategoryUpdatesCallback(): shift update: %v", err)
		return
	}

	if _, err := a.Bot.Send(
		cb.Sender,
		up.FormatMessage(),
		&telebot.SendOptions{ParseMode: telebot.ModeMarkdown},
		NewBotMenuCategoryNextUpdate(cat, up.URL).Menu,
	); err != nil {
		log.Printf("[bot] botHandleCategoryUpdatesCallback(): Failed to show update: %v", err)
	}
}

// botHandleDeleteCmd handles /delete command.
//
//	Provides user with a choise to delete his data from the service.
//	- Confirm action will be handled by botHandleDeleteConfirmCallback
//	- Cancel action will be handled by botHandleDeleteCancelCallback
func (a *App) botHandleDeleteCmd(_ context.Context, m *telebot.Message) {
	if _, err := a.Bot.Send(
		m.Sender,
		fmt.Sprintf("You data is about to be deleted from our service"),
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
		fmt.Sprintf("Your data was successfully deleted 👍"),
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
}

// botHandleDeleteCancelCallback handles cancellation callback of Delete User menu.
func (a *App) botHandleDeleteCancelCallback(_ context.Context, cb *telebot.Callback) {
	if _, err := a.Bot.Edit(
		cb.Message,
		fmt.Sprintf("Your data will not be deleted 👍"),
		&telebot.SendOptions{ParseMode: telebot.ModeMarkdown},
	); err != nil {
		log.Printf("[bot] botHandleDeleteCancelCallback() Failed to edit message: %v", err)
	}
}

// botHandleTextMessage is an arbitrary method to handle any text message that was not handled by a specific handler.
func (a *App) botHandleTextMessage(_ context.Context, _ *telebot.Message) {
}
