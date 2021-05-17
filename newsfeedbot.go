package main

import (
	"log"
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
	server.Get("/_ah/warmupHandler", warmupHandler)

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
