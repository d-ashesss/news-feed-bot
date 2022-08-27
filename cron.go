package main

import (
	"context"
	"github.com/d-ashesss/news-feed-bot/pkg/feed/fetcher"
	"github.com/d-ashesss/news-feed-bot/pkg/model"
	"log"
	"net/http"
	"time"
)

func (a *App) authCron(res http.ResponseWriter, r *http.Request) {
	if head := r.Header.Get("X-Appengine-Cron"); head != "true" {
		res.WriteHeader(http.StatusUnauthorized)
	}
}

func (a *App) handleCronFetch(res http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cats, err := a.CategoryModel.GetAll(ctx)
	if err != nil {
		log.Printf("get categories: %v", err)
		res.WriteHeader(500)
		return
	}
	for _, cat := range cats {
		a.helperFetchCategory(ctx, &cat)
	}
}

func (a *App) helperFetchCategory(ctx context.Context, cat *model.Category) {
	feeds, err := a.FeedModel.GetAll(ctx, cat)
	if err != nil {
		log.Printf("get feeds in %q: %s", cat.Name, err)
		return
	}

	f := fetcher.New(a.SubscriptionModel)
	for _, feed := range feeds {
		if err := f.Fetch(ctx, &feed, cat); err != nil {
			log.Printf("fetch updates from %q: %v", feed.Title, err)
		}
		if err := a.FeedModel.SetUpdated(ctx, &feed, time.Now().UTC()); err != nil {
			log.Printf("update feed %q: %v", feed.Title, err)
		}
	}
}
