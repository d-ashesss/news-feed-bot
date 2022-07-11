package firestore

import (
	fst "cloud.google.com/go/firestore"
	"context"
	"github.com/d-ashesss/news-feed-bot/pkg/model"
	"github.com/jschoedt/go-firestorm"
)

// subscriptionModel is a Firestore implementation of model.SubscriptionModel.
type subscriptionModel struct {
	fsc             *firestorm.FSClient   // fsc is a Firestore client.
	categoryModel   model.CategoryModel   // categoryModel is an implementation of model.CategoryModel.
	subscriberModel model.SubscriberModel // subscriberModel  is an implementation of model.SubscriberModel.
	updateModel     model.UpdateModel     // updateModel is an implementation of model.UpdateModel.
}

// NewSubscriptionModel initializes Firestore implementation of model.SubscriptionModel.
func NewSubscriptionModel(
	c *fst.Client,
	categoryModel model.CategoryModel,
	subscriberModel model.SubscriberModel,
	updateModel model.UpdateModel,
) model.SubscriptionModel {
	return subscriptionModel{
		fsc:             firestorm.New(c, "ID", ""),
		categoryModel:   categoryModel,
		subscriberModel: subscriberModel,
		updateModel:     updateModel,
	}
}

func (m subscriptionModel) Subscribe(ctx context.Context, s *model.Subscriber, cat model.Category) error {
	if s == nil || s.ID == "" {
		return model.ErrInvalidSubscriber
	}
	if cat.ID == "" {
		return model.ErrInvalidCategory
	}
	sub, err := m.subscriberModel.Get(ctx, s.UserID)
	if err != nil {
		return model.ErrInvalidSubscriber
	}
	sub.AddCategory(cat)
	if err := m.req().UpdateEntities(ctx, sub)(); err != nil {
		return err
	}
	s.Categories = sub.Categories
	return nil
}

func (m subscriptionModel) Unsubscribe(ctx context.Context, s *model.Subscriber, cat model.Category) error {
	if s == nil || s.ID == "" {
		return model.ErrInvalidSubscriber
	}
	if cat.ID == "" {
		return model.ErrInvalidCategory
	}
	sub, err := m.subscriberModel.Get(ctx, s.UserID)
	if err != nil {
		return model.ErrInvalidSubscriber
	}
	sub.RemoveCategory(cat)
	if err := m.req().UpdateEntities(ctx, sub)(); err != nil {
		return err
	}
	s.Categories = sub.Categories
	return nil
}

func (m subscriptionModel) GetSubscriptionStatus(ctx context.Context, s *model.Subscriber) ([]model.Subscription, error) {
	if s == nil {
		return nil, model.ErrInvalidSubscriber
	}
	cats, err := m.categoryModel.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	subs := make([]model.Subscription, len(cats))
	for i := range cats {
		subscribed := s.HasCategory(cats[i])
		unread, _ := m.updateModel.GetCountInCategory(ctx, s, &cats[i])

		subs[i] = model.Subscription{
			Category:   cats[i],
			Subscribed: subscribed,
			Unread:     unread,
		}
	}
	return subs, nil
}

func (m subscriptionModel) AddUpdate(ctx context.Context, up model.Update) error {
	if up.Category == nil {
		return model.ErrInvalidCategory
	}
	catRef := m.req().ToRef(up.Category)
	q := m.req().ToCollection(model.Subscriber{}).Where("categories", "array-contains", catRef)
	var ss []model.Subscriber
	if err := m.req().QueryEntities(ctx, q, &ss)(); err != nil {
		return err
	}
	for _, s := range ss {
		sup := up
		sup.Subscriber = &s
		_, _ = m.updateModel.Create(ctx, &sup)
	}
	return nil
}

func (m subscriptionModel) ShiftUpdate(ctx context.Context, s *model.Subscriber, cat model.Category) (*model.Update, error) {
	up, err := m.updateModel.GetFromCategory(ctx, s, &cat)
	if err != nil {
		return nil, err
	}
	if err := m.updateModel.Delete(ctx, up); err != nil {
		return nil, err
	}
	return up, nil
}

// req is a shortcut to firestorm.FSClient.NewRequest().
func (m subscriptionModel) req() *firestorm.Request {
	return m.fsc.NewRequest()
}
