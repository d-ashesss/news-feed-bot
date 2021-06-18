package main

import (
	"log"
	"net/http"
)

func (a *App) authCron(res http.ResponseWriter, r *http.Request) {
	if head := r.Header.Get("X-Appengine-Cron"); head != "true" {
		res.WriteHeader(http.StatusUnauthorized)
	}
}

func (a *App) handleCronFetch() {
	log.Printf("fetch done")
}
