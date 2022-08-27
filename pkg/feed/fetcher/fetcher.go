package fetcher

import (
	"context"
	"github.com/d-ashesss/news-feed-bot/pkg/model"
	"github.com/mmcdole/gofeed"
	"log"
)

// Fetcher reads the feed and extracts posts from it.
type Fetcher struct {
	subscriptionModel model.SubscriptionModel
}

// New instantiates new Fetcher.
func New(subscriptionModel model.SubscriptionModel) *Fetcher {
	return &Fetcher{subscriptionModel: subscriptionModel}
}

func (f Fetcher) GetTitle(ctx context.Context, URL string) (string, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURLWithContext(URL, ctx)
	if err != nil {
		return "", err
	}
	return feed.Title, nil
}

// Fetch reads posts from the feed into subscriptions.
func (f Fetcher) Fetch(ctx context.Context, fd *model.Feed, cat *model.Category) error {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURLWithContext(fd.URL, ctx)
	if err != nil {
		return err
	}
	log.Printf("[fetcher] fetching updates from feed %q [%s] for category %q", feed.Title, feed.Language, cat.Name)

	for _, i := range feed.Items {
		up := model.Update{
			Category: cat,
			FeedID:   i.GUID,
			Title:    i.Title,
			Date:     i.PublishedParsed.UTC(),
			URL:      i.Link,
		}
		if fd.LastUpdate.After(up.Date) {
			continue
		}
		if err := f.subscriptionModel.AddUpdate(ctx, up); err != nil {
			log.Printf("[fetcher] failed to save update: %v", err)
		}
	}
	return nil
}
