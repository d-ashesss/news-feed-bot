package firestore

import (
	fst "cloud.google.com/go/firestore"
	"context"
	"github.com/d-ashesss/news-feed-bot/pkg/model"
	"github.com/jschoedt/go-firestorm"
)

// subscriberModel is a Firestore implementation of model.SubscriberModel.
type subscriberModel struct {
	fsc         *firestorm.FSClient // fsc is a Firestore client.
	updateModel model.UpdateModel   // updateModel is an implementation of model.UpdateModel.
}

// NewSubscriberModel initializes Firestore implementation of model.SubscriberModel.
func NewSubscriberModel(c *fst.Client, updateModel model.UpdateModel) model.SubscriberModel {
	return subscriberModel{
		fsc:         firestorm.New(c, "ID", ""),
		updateModel: updateModel,
	}
}

func (m subscriberModel) Create(ctx context.Context, s *model.Subscriber) (string, error) {
	if s == nil {
		return "", model.ErrInvalidSubscriber
	}
	if s.UserID == "" {
		return "", model.ErrInvalidSubscriberID
	}
	if err := m.req().CreateEntities(ctx, s)(); err != nil {
		return "", err
	}
	return m.req().GetID(s), nil
}

func (m subscriberModel) Get(ctx context.Context, id string) (*model.Subscriber, error) {
	var ss []model.Subscriber
	q := m.req().ToCollection(model.Subscriber{}).Where("userid", "==", id).Limit(1)
	if err := m.req().SetLoadPaths(firestorm.AllEntities).QueryEntities(ctx, q, &ss)(); err != nil {
		return nil, err
	}
	if len(ss) == 0 {
		return nil, model.ErrNotFound
	}
	return &ss[0], nil
}

func (m subscriberModel) Delete(ctx context.Context, s *model.Subscriber) error {
	if s == nil || s.ID == "" {
		return model.ErrInvalidSubscriber
	}
	if err := m.updateModel.DeleteForSubscriber(ctx, s); err != nil {
		return err
	}
	return m.req().DeleteEntities(ctx, s)()
}

// req is a shortcut to firestorm.FSClient.NewRequest().
func (m subscriberModel) req() *firestorm.Request {
	return m.fsc.NewRequest()
}
