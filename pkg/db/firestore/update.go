package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/d-ashesss/news-feed-bot/pkg/model"
	"github.com/jschoedt/go-firestorm"
)

// updateModel is a Firestore implementation of model.UpdateModel.
type updateModel struct {
	fsc *firestorm.FSClient
}

// NewUpdateModel initializes Firestore implementation of model.UpdateModel.
func NewUpdateModel(c *firestore.Client) model.UpdateModel {
	return updateModel{fsc: firestorm.New(c, "ID", "Subscriber")}
}

func (m updateModel) Create(ctx context.Context, up *model.Update) (string, error) {
	if up == nil {
		return "", model.ErrInvalidUpdate
	}
	if up.Subscriber == nil || len(up.Subscriber.ID) == 0 {
		return "", model.ErrInvalidSubscriber
	}
	if up.Category == nil || len(up.Category.ID) == 0 {
		return "", model.ErrInvalidCategory
	}
	if err := m.req().CreateEntities(ctx, up)(); err != nil {
		return "", err
	}
	return m.req().GetID(up), nil
}

func (m updateModel) GetFromCategory(ctx context.Context, s *model.Subscriber, cat *model.Category) (*model.Update, error) {
	if s == nil || len(s.ID) == 0 {
		return nil, model.ErrInvalidSubscriber
	}
	if cat == nil || len(cat.ID) == 0 {
		return nil, model.ErrInvalidCategory
	}
	var ups []model.Update
	catRef := m.req().ToRef(cat)
	q := m.req().ToCollection(model.Update{Subscriber: s}).
		Where("category", "==", catRef).
		OrderBy("date", firestore.Asc).
		Limit(1)
	if err := m.req().SetLoadPaths(firestorm.AllEntities).QueryEntities(ctx, q, &ups)(); err != nil {
		return nil, err
	}
	if len(ups) == 0 {
		return nil, model.ErrNoUpdates
	}
	return &ups[0], nil
}

func (m updateModel) GetCountInCategory(ctx context.Context, s *model.Subscriber, cat *model.Category) (int, error) {
	if s == nil || len(s.ID) == 0 {
		return 0, model.ErrInvalidSubscriber
	}
	if cat == nil || len(cat.ID) == 0 {
		return 0, model.ErrInvalidCategory
	}
	var ups []model.Update
	catRef := m.req().ToRef(cat)
	q := m.req().ToCollection(model.Update{Subscriber: s}).
		Where("category", "==", catRef).
		OrderBy("date", firestore.Asc)
	if err := m.req().QueryEntities(ctx, q, &ups)(); err != nil {
		return 0, err
	}
	return len(ups), nil
}

func (m updateModel) Delete(ctx context.Context, up *model.Update) error {
	if up == nil || len(up.ID) == 0 {
		return model.ErrInvalidUpdate
	}
	if up.Subscriber == nil || len(up.Subscriber.ID) == 0 {
		return model.ErrInvalidSubscriber
	}
	if err := m.req().DeleteEntities(ctx, up)(); err != nil {
		return err
	}
	return nil
}

func (m updateModel) DeleteForSubscriber(ctx context.Context, s *model.Subscriber) error {
	if s == nil || len(s.ID) == 0 {
		return model.ErrInvalidSubscriber
	}
	var ups []model.Update
	q := m.req().ToCollection(model.Update{Subscriber: s}).Query
	if err := m.req().SetLoadPaths(firestorm.AllEntities).QueryEntities(ctx, q, &ups)(); err != nil {
		return err
	}
	_ = m.req().DeleteEntities(ctx, ups)()
	return nil
}

// req is a shortcut to firestorm.FSClient.NewRequest().
func (m updateModel) req() *firestorm.Request {
	return m.fsc.NewRequest()
}
