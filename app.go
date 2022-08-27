package main

import (
	"github.com/d-ashesss/news-feed-bot/bot"
	"github.com/d-ashesss/news-feed-bot/http"
	"github.com/d-ashesss/news-feed-bot/pkg/model"
	"github.com/go-martini/martini"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	Config            Config
	HttpServer        *http.Server
	Bot               *bot.Bot
	FeedModel         model.FeedModel
	CategoryModel     model.CategoryModel
	SubscriberModel   model.SubscriberModel
	SubscriptionModel model.SubscriptionModel
}

func (a *App) Run() {
	if a.Bot == nil {
		log.Printf("[app] Running in botless mode")
	} else {
		log.Printf("[app] Serving for TG bot %v", a.Bot.GetName())
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		log.Printf("[app] Starting HTTP server")
		if err := a.HttpServer.Run(); err != nil {
			log.Printf("[app] Http server error: %s", err)
		}
		signals <- syscall.SIGQUIT
	}()

	if !a.Config.BotWebhookMode {
		go func() {
			log.Printf("[app] Starting TG Bot")
			a.Bot.Start()
		}()
	}

	sig := <-signals
	log.Printf("[app] Received signal %s", sig)
	log.Printf("[app] Stopping HTTP server")
	a.HttpServer.Shutdown()
	log.Printf("[app] Gracefully exiting")
}

func NewApp(
	config Config,
	httpServer *http.Server,
	feedModel model.FeedModel,
	categoryModel model.CategoryModel,
	subscriberModel model.SubscriberModel,
	subscriptionModel model.SubscriptionModel,
) *App {
	app := &App{
		Config:            config,
		HttpServer:        httpServer,
		FeedModel:         feedModel,
		CategoryModel:     categoryModel,
		SubscriberModel:   subscriberModel,
		SubscriptionModel: subscriptionModel,
	}

	app.HttpServer.Get("/", app.handleIndex)
	app.HttpServer.Get("/_ah/warmup", app.handleWarmup)
	app.HttpServer.Group("/cron", func(r martini.Router) {
		r.Get("/fetch", app.handleCronFetch)
	}, app.authCron)

	return app
}
