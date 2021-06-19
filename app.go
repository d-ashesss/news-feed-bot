package main

import (
	"github.com/d-ashesss/news-feed-bot/bot"
	"github.com/d-ashesss/news-feed-bot/http"
	"github.com/go-martini/martini"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	Config     Config
	HttpServer *http.Server
	Bot        *bot.Bot
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

	sig := <-signals
	log.Printf("[app] Received signal %s", sig)
	log.Printf("[app] Stopping HTTP server")
	a.HttpServer.Shutdown()
	log.Printf("[app] Gracefully exiting")
}

func NewApp(config Config, httpServer *http.Server) *App {
	app := &App{
		Config:     config,
		HttpServer: httpServer,
	}

	app.HttpServer.Get("/", app.handleIndex)
	app.HttpServer.Get("/_ah/warmup", app.handleWarmup)
	app.HttpServer.Group("/cron", func(r martini.Router) {
		r.Get("/fetch", app.handleCronFetch)
	}, app.authCron)

	return app
}
