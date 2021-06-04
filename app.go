package main

import (
	"NewsFeedBot/bot"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/google/uuid"
	"gopkg.in/tucnak/telebot.v2"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
)

var (
	projectID, baseURL string
)

func init() {
	projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
	baseURL = os.Getenv("APP_BASE_URL")
	if len(baseURL) == 0 {
		if len(projectID) > 0 {
			baseURL = "https://" + projectID + ".appspot.com"
		} else {
			baseURL = "http://localhost"
		}
	}
}

type App struct {
	HttpServer *Server
	Bot        *bot.Bot
}

func (a *App) SetBot(bot *bot.Bot) error {
	if bot == nil {
		return errors.New("invalid bot instance")
	}
	p, err := getBotWebhookPath(bot)
	if err != nil {
		p, err = createBotWebhookPath(bot)
		if err != nil {
			return errors.New(fmt.Sprintf("unable to create webhook: %v", err))
		}
	}
	a.Bot = bot
	a.HttpServer.Post(p, a.handleBotUpdate)
	a.Bot.Handle(telebot.OnText, a.handleBotMessage)
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

func createBotWebhookPath(bot *bot.Bot) (string, error) {
	p := "/update/" + uuid.NewString()
	if err := bot.SetWebhookURL(baseURL + p); err != nil {
		return "", errors.New(fmt.Sprintf("unable to set webhook: %v", err))
	}
	log.Printf("Created new TG webhook %q", baseURL+p)
	return p, nil
}

func (a *App) Run() {
	if a.Bot == nil {
		log.Printf("Running in botless mode")
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		log.Printf("Starting HTTP server")
		if err := a.HttpServer.Run(); err != nil {
			log.Printf("Http server error: %s", err)
		}
		signals <- syscall.SIGQUIT
	}()

	sig := <-signals
	log.Printf("Received signal %s", sig)
	log.Printf("Stopping HTTP server")
	a.HttpServer.Shutdown()
	log.Printf("Gracefully exiting")
}

func (a *App) handleIndex() string {
	return "Hello, World!"
}

func (a *App) handleWarmup(log *log.Logger) {
	log.Printf("warmup done")
}

func (a *App) authCron(res http.ResponseWriter, r *http.Request) {
	if head := r.Header.Get("X-Appengine-Cron"); head != "true" {
		res.WriteHeader(http.StatusUnauthorized)
	}
}

func (a *App) handleCronFetch() {
	log.Printf("fetch done")
}

func (a *App) handleBotUpdate(res http.ResponseWriter, r *http.Request) {
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

func (a *App) handleBotMessage(m *telebot.Message) {
	log.Printf("[bot] Incoiming message from %s: %q", bot.GetUserName(m.Sender), m.Text)
	if m.Text == "/start" {
		if _, err := a.Bot.Send(m.Sender, "Welcome 🎉"); err != nil {
			log.Printf("[bot] Failed to reply: %v", err)
		}
	}
}

func NewApp(httpServer *Server) *App {
	app := &App{
		HttpServer: httpServer,
	}

	app.HttpServer.Get("/", app.handleIndex)
	app.HttpServer.Get("/_ah/warmup", app.handleWarmup)
	app.HttpServer.Group("/cron", func(r martini.Router) {
		r.Get("/fetch", app.handleCronFetch)
	}, app.authCron)

	return app
}
