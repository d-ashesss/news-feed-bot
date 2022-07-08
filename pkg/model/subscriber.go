package model

import "context"

// Subscriber represents subscriber entitiy.
type Subscriber struct {
	ID         string     // ID is an internal DB ID of the user.
	UserID     string     // UserID is an external ID of the user. Like Telegram user ID.
	Categories []Category // Categories is a list of Category'ies the user is subscribed to.
}

// NewSubscriber initializes new Subscriber.
func NewSubscriber(userID string) *Subscriber {
	return &Subscriber{UserID: userID}
}

// AddCategory adds a Category to the list of Subscriber's subscriptions.
func (s *Subscriber) AddCategory(c Category) {
	s.Categories = append(s.Categories, c)
}

// RemoveCategory removes a Category from the list of Subscriber's subscriptions.
func (s *Subscriber) RemoveCategory(c Category) {
	subs := make([]Category, 0, len(s.Categories))
	for _, cat := range s.Categories {
		if cat.ID != c.ID {
			subs = append(subs, cat)
		}
	}
	s.Categories = subs
}

type SubscriberModel interface {
	// Create saves a Subscriber entity into the DB.
	Create(ctx context.Context, s *Subscriber) (string, error)
	// Get retrieves a Subscriber entity from the DB by external UserID.
	Get(ctx context.Context, id string) (*Subscriber, error)
	// Delete deletes a Subscriber entity from the DB.
	Delete(ctx context.Context, s *Subscriber) error
}
