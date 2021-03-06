package main

import "gopkg.in/tucnak/telebot.v2"

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
