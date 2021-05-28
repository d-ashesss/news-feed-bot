package main

import (
	"log"
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("[main] ")
}

func main() {
	httpServer := NewHttpServer()
	app := NewApp(httpServer)
	app.Run()
}
