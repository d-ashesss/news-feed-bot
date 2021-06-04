package bot

import (
	"errors"
	"gopkg.in/tucnak/telebot.v2"
	"os"
	"strings"
)

var token string

func init() {
	token = os.Getenv("TG_TOKEN")
}

func New() (*Bot, error) {
	bot, err := telebot.NewBot(telebot.Settings{
		Token: token,
	})
	if err != nil {
		return nil, err
	}
	return &Bot{
		Bot: bot,
	}, nil
}

type Bot struct {
	*telebot.Bot
}

// WebhookURL retrieves the webhook URL if it was previously set.
func (b *Bot) WebhookURL() (string, error) {
	h, err := b.GetWebhook()
	if err != nil {
		return "", err
	}
	return b.webhookURL(h)
}

func (b *Bot) webhookURL(h *telebot.Webhook) (string, error) {
	if h == nil || len(h.Listen) == 0 {
		return "", errors.New("bot: received webhook is invalid")
	}
	return h.Listen, nil
}

// SetWebhookURL sets new webhook URL.
func (b *Bot) SetWebhookURL(u string) error {
	h := &telebot.Webhook{Endpoint: &telebot.WebhookEndpoint{PublicURL: u}}
	if err := b.SetWebhook(h); err != nil {
		return err
	}
	return nil
}

func GetUserName(user *telebot.User) string {
	if len(user.Username) > 0 {
		return "@" + user.Username
	}
	return strings.Trim(user.FirstName+" "+user.LastName, " ")
}
