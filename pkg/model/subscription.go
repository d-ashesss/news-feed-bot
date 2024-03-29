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
	// Subscribe subsribes the Subscriber to a Category
	Subscribe(ctx context.Context, s *Subscriber, cat Category) error
	// Unsubscribe unsubscribes the Subscriber from a Category
	Unsubscribe(ctx context.Context, s *Subscriber, cat Category) error
	// GetCategorySubscription returns a subscription status of a given Category for a given Subscriber.
	GetCategorySubscription(ctx context.Context, s *Subscriber, cat Category) (*Subscription, error)
	// GetSubscriptionStatus returns a list of all categories and their subscription status for a given Subscriber.
	GetSubscriptionStatus(ctx context.Context, s *Subscriber) ([]Subscription, error)
	// AddUpdate adds and update to each subscriber of a category. Update has to have its Category property set.
	AddUpdate(ctx context.Context, up Update) error
	// ShiftUpdate retrieves an Update for selected Category removing it from Subscriber's list of unread updates.
	ShiftUpdate(ctx context.Context, s *Subscriber, cat Category) (*Update, error)
}
