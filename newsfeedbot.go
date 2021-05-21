package main

import (
	"github.com/go-martini/martini"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("[main] ")
}

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	server := NewHttpServer()

	server.Get("/", indexHandler)
	server.Get("/_ah/warmup", warmupHandler)
	server.Group("/cron", func(r martini.Router) {
		r.Get("/fetch", cronFetchHandler)
	}, cronAuth)

	go func() {
		log.Printf("Starting HTTP server")
		if err := server.Run(); err != nil {
			log.Printf("Http server error: %s", err)
		}
		signals <- syscall.SIGQUIT
	}()

	sig := <-signals
	log.Printf("Received signal %s", sig)
	log.Printf("Stopping HTTP server")
	server.Shutdown()
	log.Printf("Gracefully exiting")
}

func indexHandler() string {
	return "Hello, World!"
}

func warmupHandler(log *log.Logger) {
	log.Printf("warmup done")
}

func cronAuth(res http.ResponseWriter, r *http.Request) {
	if head := r.Header.Get("X-Appengine-Cron"); head != "true" {
		res.WriteHeader(http.StatusUnauthorized)
	}
}

func cronFetchHandler() {
	log.Printf("fetch done")
}
