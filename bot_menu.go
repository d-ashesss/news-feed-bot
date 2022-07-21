package main

import (
	"fmt"
	"github.com/d-ashesss/news-feed-bot/pkg/model"
	"gopkg.in/tucnak/telebot.v2"
)

const (
	BotMenuMainBtnCheckUpdatesLabel = "Check for updates"
	BotMenuMainBtnCheckUpdatesID    = "btnMenuMainCheckUpdates"

	BotMenuMainBtnSelectCategoriesLabel = "Select categories"
	BotMenuMainBtnSelectCategoriesID    = "btnMenuMainSelectCategories"
)

type BotMenuMain struct {
	Menu *telebot.ReplyMarkup

	BtnCheckUpdates     telebot.Btn
	BtnSelectCategories telebot.Btn
}

func NewBotMenuMain() *BotMenuMain {
	m := &BotMenuMain{
		Menu: &telebot.ReplyMarkup{},
	}
	m.BtnCheckUpdates = m.Menu.Data(BotMenuMainBtnCheckUpdatesLabel, BotMenuMainBtnCheckUpdatesID)
	m.BtnSelectCategories = m.Menu.Data(BotMenuMainBtnSelectCategoriesLabel, BotMenuMainBtnSelectCategoriesID)
	m.Menu.Inline(
		m.Menu.Row(m.BtnCheckUpdates),
		m.Menu.Row(m.BtnSelectCategories),
	)
	return m
}

type BotMenuNoUpdatesInCategory struct {
	Menu *telebot.ReplyMarkup

	BtnCheckUpdates telebot.Btn
}

func NewBotMenuNoUpdatesInCategory() *BotMenuNoUpdatesInCategory {
	m := &BotMenuNoUpdatesInCategory{
		Menu: &telebot.ReplyMarkup{},
	}
	m.BtnCheckUpdates = m.Menu.Data(BotMenuMainBtnCheckUpdatesLabel, BotMenuMainBtnCheckUpdatesID)
	m.Menu.Inline(
		m.Menu.Row(m.BtnCheckUpdates),
	)
	return m
}

type BotMenuNoCategoriesSelected struct {
	Menu *telebot.ReplyMarkup

	BtnSelectCategories telebot.Btn
}

func NewBotMenuNoCategoriesSelected() *BotMenuNoCategoriesSelected {
	m := &BotMenuNoCategoriesSelected{
		Menu: &telebot.ReplyMarkup{},
	}
	m.BtnSelectCategories = m.Menu.Data(BotMenuMainBtnSelectCategoriesLabel, BotMenuMainBtnSelectCategoriesID)
	m.Menu.Inline(
		m.Menu.Row(m.BtnSelectCategories),
	)
	return m
}

type BotMenuCategoryUpdates struct {
	Menu *telebot.ReplyMarkup
}

const BotMenuBtnCategoryUpdatesID = "btnMenuCategoryUpdates"

func NewBotMenuCategoryUpdates(subs []model.Subscription) *BotMenuCategoryUpdates {
	m := &BotMenuCategoryUpdates{
		Menu: &telebot.ReplyMarkup{},
	}
	rows := make([]telebot.Row, 0, len(subs))
	for _, sub := range subs {
		label := fmt.Sprintf("%s (%d)", sub.Category.Name, sub.Unread)
		btn := m.Menu.Data(label, BotMenuBtnCategoryUpdatesID, sub.Category.ID)
		rows = append(rows, m.Menu.Row(btn))
	}
	m.Menu.Inline(rows...)
	return m
}

type BotMenuCategoryNextUpdate struct {
	Menu *telebot.ReplyMarkup

	BtnNext telebot.Btn
}

const (
	BotMenuBtnCategoryNextUpdateLabel = "Next"
	BotMenuBtnCategoryNextUpdateID    = "btnMenuCategoryNextUpdate"
)

func NewBotMenuCategoryNextUpdate(cat *model.Category) *BotMenuCategoryNextUpdate {
	m := &BotMenuCategoryNextUpdate{
		Menu: &telebot.ReplyMarkup{},
	}
	m.BtnNext = m.Menu.Data(BotMenuBtnCategoryNextUpdateLabel, BotMenuBtnCategoryNextUpdateID, cat.ID)
	m.Menu.Inline(m.Menu.Row(m.BtnNext))
	return m
}

const BotMenuBtnToggleCategoryID = "btnMenuToggleCategory"

type BotMenuSelectCategories struct {
	Menu *telebot.ReplyMarkup
}

func NewBotMenuSelectCategories(subs []model.Subscription) *BotMenuSelectCategories {
	m := &BotMenuSelectCategories{
		Menu: &telebot.ReplyMarkup{},
	}
	rows := make([]telebot.Row, 0, len(subs))
	for _, sub := range subs {
		label := sub.Category.Name
		if sub.Subscribed {
			label = "✅ " + label
		}
		btn := m.Menu.Data(label, BotMenuBtnToggleCategoryID, sub.Category.ID)
		rows = append(rows, m.Menu.Row(btn))
	}
	m.Menu.Inline(rows...)
	return m
}

const (
	BotMenuDeleteBtnConfirmLabel = "✔️ Confirm"
	BotMenuDeleteBtnConfirmID    = "btnMenuDeleteConfirm"
	BotMenuDeleteBtnCancelLabel  = "❌ Cancel"
	BotMenuDeleteBtnCancelID     = "btnMenuDeleteCancel"
)

// BotMenuDelete represents Delete User menu.
type BotMenuDelete struct {
	Menu *telebot.ReplyMarkup

	BtnConfirm telebot.Btn
	BtnCancel  telebot.Btn
}

// NewBotMenuDelete initializes new BotMenuDelete.
func NewBotMenuDelete() *BotMenuDelete {
	m := &BotMenuDelete{
		Menu: &telebot.ReplyMarkup{},
	}
	m.BtnConfirm = m.Menu.Data(BotMenuDeleteBtnConfirmLabel, BotMenuDeleteBtnConfirmID)
	m.BtnCancel = m.Menu.Data(BotMenuDeleteBtnCancelLabel, BotMenuDeleteBtnCancelID)
	m.Menu.Inline(m.Menu.Row(m.BtnConfirm, m.BtnCancel))
	return m
}
