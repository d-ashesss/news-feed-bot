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
	HttpServer *http.Server
	Bot        *bot.Bot
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

func NewApp(httpServer *http.Server) *App {
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
