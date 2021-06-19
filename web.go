package main

import "log"

func (a *App) handleIndex() string {
	return "Hello, World!"
}

func (a *App) handleWarmup(log *log.Logger) {
	log.Printf("[web] warmup done")
}
