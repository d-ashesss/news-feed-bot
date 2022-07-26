package model

import (
	"context"
	"time"
)

// Update represents update entitiy.
type Update struct {
	ID         string      // ID is an internal ID.
	Subscriber *Subscriber // Subscriber is the receiver of the update.
	Category   *Category   // Category is the category of the update.
	Title      string      // Title is the title of the update.
	Date       time.Time   // Date is the date when the update was published.
}

// UpdateModel is a data model for Update.
type UpdateModel interface {
	// Create saves an Update entity into the DB.
	Create(ctx context.Context, c *Update) (string, error)
	// GetFromCategory retrieves the oldest available update from selected Category for the Subscriber.
	GetFromCategory(ctx context.Context, s *Subscriber, cat *Category) (*Update, error)
	// GetCountInCategory retrieves the number of updates available in selected Category for the Subscriber.
	GetCountInCategory(ctx context.Context, s *Subscriber, cat *Category) (int, error)
	// Delete deletes an Update entity from the DB.
	Delete(ctx context.Context, up *Update) error
	// DeleteForSubscriber deletes all Update's for the Subscriber.
	DeleteForSubscriber(ctx context.Context, s *Subscriber) error
}
