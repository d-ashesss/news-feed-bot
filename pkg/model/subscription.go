package model

import "context"

// Subscription represents a status of subscription for a single Category.
type Subscription struct {
	Category   Category // Category is a Category in question.
	Subscribed bool     // Subscribed shows if a Subscriber is subscribed to current Category.
	Unread     int      // Unread shows count of unread updates from current Category.
}

// SubscriptionModel is a data model for Subscription.
type SubscriptionModel interface {
	// CreateSubscriber saves a Subscriber entity into the DB.
	CreateSubscriber(ctx context.Context, s *Subscriber) (string, error)
	// GetSubscriber retrieves a Subscriber entity from the DB.
	GetSubscriber(ctx context.Context, id string) (*Subscriber, error)
	// Subscribe subsribes the Subscriber to a Category
	Subscribe(ctx context.Context, s *Subscriber, c Category) error
	// Unsubscribe unsubscribes the Subscriber from a Category
	Unsubscribe(ctx context.Context, s *Subscriber, c Category) error
	// GetSubscriptionStatus returns a list of all categories and their subscription status for a given Subscriber.
	GetSubscriptionStatus(ctx context.Context, s *Subscriber) ([]Subscription, error)
}
