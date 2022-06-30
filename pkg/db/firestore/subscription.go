package firestore

import (
	fst "cloud.google.com/go/firestore"
	"context"
	"github.com/d-ashesss/news-feed-bot/pkg/model"
	"github.com/jschoedt/go-firestorm"
)

type subscriptionModel struct {
	fsc           *firestorm.FSClient
	categoryModel model.CategoryModel
}

func NewSubscriptionModel(c *fst.Client, categoryModel model.CategoryModel) model.SubscriptionModel {
	return subscriptionModel{fsc: firestorm.New(c, "ID", ""), categoryModel: categoryModel}
}

func (m subscriptionModel) CreateSubscriber(ctx context.Context, s *model.Subscriber) (string, error) {
	if err := m.req().CreateEntities(ctx, s)(); err != nil {
		return "", err
	}
	return m.req().GetID(s), nil
}

func (m subscriptionModel) GetSubscriber(ctx context.Context, id string) (*model.Subscriber, error) {
	s := &model.Subscriber{ID: id}
	if _, err := m.req().SetLoadPaths(firestorm.AllEntities).GetEntities(ctx, s)(); err != nil {
		return nil, err
	}
	return s, nil
}

func (m subscriptionModel) Subscribe(ctx context.Context, s *model.Subscriber, c model.Category) error {
	sub, err := m.GetSubscriber(ctx, s.ID)
	if err != nil {
		return err
	}
	sub.AddCategory(c)
	if err := m.req().UpdateEntities(ctx, sub)(); err != nil {
		return err
	}
	s.Categories = sub.Categories
	return nil
}

func (m subscriptionModel) Unsubscribe(ctx context.Context, s *model.Subscriber, c model.Category) error {
	sub, err := m.GetSubscriber(ctx, s.ID)
	if err != nil {
		return err
	}
	sub.RemoveCategory(c)
	if err := m.req().UpdateEntities(ctx, sub)(); err != nil {
		return err
	}
	s.Categories = sub.Categories
	return nil
}

func (m subscriptionModel) Subscribed(ctx context.Context, s *model.Subscriber, c model.Category) (bool, error) {
	sub, err := m.GetSubscriber(ctx, s.ID)
	if err != nil {
		return false, err
	}
	for _, cat := range sub.Categories {
		if cat.ID == c.ID {
			return true, nil
		}
	}
	return false, nil
}

func (m subscriptionModel) GetSubscriptionStatus(ctx context.Context, s *model.Subscriber) ([]model.Subscription, error) {
	cats, err := m.categoryModel.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	subs := make([]model.Subscription, len(cats))
	for i := range cats {
		subscribed, _ := m.Subscribed(ctx, s, cats[i])
		subs[i] = model.Subscription{
			Category:   cats[i],
			Subscribed: subscribed,
		}
	}
	return subs, nil
}

func (m subscriptionModel) req() *firestorm.Request {
	return m.fsc.NewRequest()
}
