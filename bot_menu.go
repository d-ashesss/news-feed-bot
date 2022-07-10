package main

import (
	"github.com/d-ashesss/news-feed-bot/pkg/model"
	"gopkg.in/tucnak/telebot.v2"
)

const (
	BotMenuMainBtnCheckUpdatesLabel = "Check for updates"
	BotMenuMainBtnCheckUpdatesID    = "btnMenuMainCheckUpdates"
)

type BotMenuMain struct {
	Menu *telebot.ReplyMarkup

	BtnCheckUpdates telebot.Btn
}

func NewBotMenuMain() *BotMenuMain {
	m := &BotMenuMain{
		Menu: &telebot.ReplyMarkup{},
	}
	m.BtnCheckUpdates = m.Menu.Data(BotMenuMainBtnCheckUpdatesLabel, BotMenuMainBtnCheckUpdatesID)
	m.Menu.Inline(
		m.Menu.Row(m.BtnCheckUpdates),
	)
	return m
}

type BotMenuCategoriesUpdates struct {
	Menu *telebot.ReplyMarkup
}

func NewBotMenuCategoriesUpdates(_ []model.Subscription) *BotMenuCategoriesUpdates {
	m := &BotMenuCategoriesUpdates{
		Menu: &telebot.ReplyMarkup{},
	}
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
