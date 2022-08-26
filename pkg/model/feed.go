package model

import (
	"context"
	"time"
)

type Feed struct {
	ID         string    // ID is an internal ID.
	Category   *Category // Category is the category of the feed.
	Title      string    // Title is the title of the feed.
	URL        string    // URL is a http link to the feed.
	LastUpdate time.Time // LastUpdate is the published time of the last update fetched from the feed.
}

type FeedModel interface {
	// Create saves a Feed entity into the DB.
	Create(ctx context.Context, f *Feed) (string, error)
	// Get retrieves a Feed entity from the DB.
	Get(ctx context.Context, cat *Category, id string) (*Feed, error)
	// GetAll retrieves all Feed entities for provided Category from the DB.
	GetAll(ctx context.Context, cat *Category) ([]Feed, error)
	SetUpdated(ctx context.Context, f *Feed, u time.Time) error
	// Delete deletes a Feed entity from the DB. Category property has to be set on Feed entity.
	Delete(ctx context.Context, f *Feed) error
}
