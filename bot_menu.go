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

const (
	BotBtnBackToMainMenuLabel = "⬅️ Back to main menu"
	BotBtnBackToMainMenuID    = "btnBackToMainMenu"
)

type BotMenuNoUpdatesInCategory struct {
	Menu *telebot.ReplyMarkup

	BtnBack telebot.Btn
}

func NewBotMenuNoUpdatesInCategory() *BotMenuNoUpdatesInCategory {
	m := &BotMenuNoUpdatesInCategory{
		Menu: &telebot.ReplyMarkup{},
	}
	m.BtnBack = m.Menu.Data(BotMenuCategoryNextUpdateBtnBackLabel, BotMenuMainBtnCheckUpdatesID)
	m.Menu.Inline(
		m.Menu.Row(m.BtnBack),
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
	backBtn := m.Menu.Data(BotBtnBackToMainMenuLabel, BotBtnBackToMainMenuID)
	m.Menu.Inline(
		m.Menu.Row(m.BtnSelectCategories),
		m.Menu.Row(backBtn),
	)
	return m
}

type BotMenuCategoryUpdates struct {
	Menu *telebot.ReplyMarkup
}

const BotMenuCategoryUpdatesBtnCategoryUpdatesID = "btnMenuCategoryUpdates"

func NewBotMenuCategoryUpdates(subs []model.Subscription) *BotMenuCategoryUpdates {
	m := &BotMenuCategoryUpdates{
		Menu: &telebot.ReplyMarkup{},
	}
	rows := make([]telebot.Row, 0, len(subs)+1)
	for _, sub := range subs {
		label := fmt.Sprintf("%s (%d)", sub.Category.Name, sub.Unread)
		btn := m.Menu.Data(label, BotMenuCategoryUpdatesBtnCategoryUpdatesID, sub.Category.ID)
		rows = append(rows, m.Menu.Row(btn))
	}
	backBtn := m.Menu.Data(BotBtnBackToMainMenuLabel, BotBtnBackToMainMenuID)
	rows = append(rows, m.Menu.Row(backBtn))
	m.Menu.Inline(rows...)
	return m
}

type BotMenuCategoryNextUpdate struct {
	Menu *telebot.ReplyMarkup

	BtnBack telebot.Btn
	BtnNext telebot.Btn
}

const (
	BotMenuCategoryNextUpdateBtnBackLabel = "⬅️ Back to categories️"
	BotMenuCategoryNextUpdateBtnNextLabel = "Next ➡️"
)

func NewBotMenuCategoryNextUpdate(cat *model.Category) *BotMenuCategoryNextUpdate {
	m := &BotMenuCategoryNextUpdate{
		Menu: &telebot.ReplyMarkup{},
	}
	m.BtnBack = m.Menu.Data(BotMenuCategoryNextUpdateBtnBackLabel, BotMenuMainBtnCheckUpdatesID, cat.ID)
	m.BtnNext = m.Menu.Data(BotMenuCategoryNextUpdateBtnNextLabel, BotMenuCategoryUpdatesBtnCategoryUpdatesID, cat.ID)
	m.Menu.Inline(m.Menu.Row(m.BtnBack, m.BtnNext))
	return m
}

const BotMenuSelectCategoriesBtnToggleCategoryID = "btnMenuToggleCategory"

type BotMenuSelectCategories struct {
	Menu *telebot.ReplyMarkup
}

func NewBotMenuSelectCategories(subs []model.Subscription) *BotMenuSelectCategories {
	m := &BotMenuSelectCategories{
		Menu: &telebot.ReplyMarkup{},
	}
	rows := make([]telebot.Row, 0, len(subs)+1)
	for _, sub := range subs {
		label := sub.Category.Name
		if sub.Subscribed {
			label = "✅ " + label
		}
		btn := m.Menu.Data(label, BotMenuSelectCategoriesBtnToggleCategoryID, sub.Category.ID)
		rows = append(rows, m.Menu.Row(btn))
	}
	backBtn := m.Menu.Data(BotBtnBackToMainMenuLabel, BotBtnBackToMainMenuID)
	rows = append(rows, m.Menu.Row(backBtn))
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
