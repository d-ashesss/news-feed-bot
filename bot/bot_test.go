package bot

import (
	"gopkg.in/tucnak/telebot.v2"
	"testing"
)

func TestBot_webhookURL(t *testing.T) {
	tests := []struct {
		name    string
		webhook *telebot.Webhook
		want    string
	}{
		{
			name:    "Nil",
			webhook: nil,
			want:    "",
		},
		{
			name:    "EmptyURL",
			webhook: &telebot.Webhook{Listen: ""},
			want:    "",
		},
		{
			name:    "ValidURL",
			webhook: &telebot.Webhook{Listen: "http://localhost/path"},
			want:    "http://localhost/path",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := Bot{}
			got, err := b.webhookURL(tt.webhook)

			if len(tt.want) == 0 && err == nil {
				t.Errorf("want error")
			}
			if len(tt.want) > 0 && err != nil {
				t.Errorf("got error %v", err)
			}
			if len(tt.want) > 0 && got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGetUserName(t *testing.T) {
	tests := []struct {
		name string
		user *telebot.User
		want string
	}{
		{
			name: "Username",
			user: &telebot.User{Username: "u-name"},
			want: "@u-name",
		},
		{
			name: "Firstname",
			user: &telebot.User{FirstName: "f-name"},
			want: "f-name",
		},
		{
			name: "Lastname",
			user: &telebot.User{LastName: "l-name"},
			want: "l-name",
		},
		{
			name: "Fullname",
			user: &telebot.User{FirstName: "f-name", LastName: "l-name"},
			want: "f-name l-name",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetUserName(tt.user)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}
