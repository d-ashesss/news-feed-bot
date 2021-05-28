package main

import (
	"github.com/go-martini/martini"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	HttpServer *Server
}

func (a *App) Run() {
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
