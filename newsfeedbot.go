package main

import (
	"NewsFeedBot/bot"
	"log"
	"os"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(0)
	log.SetPrefix("[main] ")
}

func main() {
	httpServer := NewHttpServer()
	app := NewApp(httpServer)

	b, err := bot.New()
	if err != nil {
		log.Printf("Was not able to init TG bot: %v", err)
	} else {
		if err = app.SetBot(b); err != nil {
			log.Printf("Failed to init TG bot for the app: %v", err)
		}
	}

	app.Run()
}
