package firestore

import (
	fst "cloud.google.com/go/firestore"
	"context"
	"github.com/d-ashesss/news-feed-bot/pkg/model"
	"github.com/jschoedt/go-firestorm"
)

// FeedModel is a Firestore implementation of model.FeedModel.
type FeedModel struct {
	fsc *firestorm.FSClient // fsc is a Firestore client.
}

// NewFeedModel initializes Firestore implementation of model.FeedModel.
func NewFeedModel(c *fst.Client) model.FeedModel {
	return FeedModel{fsc: firestorm.New(c, "ID", "Category")}
}

func (m FeedModel) Create(ctx context.Context, f *model.Feed) (string, error) {
	if f == nil {
		return "", model.ErrInvalidFeed
	}
	if f.Category == nil || len(f.Category.ID) == 0 {
		return "", model.ErrInvalidCategory
	}
	if err := m.req().CreateEntities(ctx, f)(); err != nil {
		return "", err
	}
	return m.req().GetID(f), nil
}

func (m FeedModel) Get(ctx context.Context, cat *model.Category, id string) (*model.Feed, error) {
	if id == "" {
		return nil, model.ErrNotFound
	}
	if cat == nil || len(cat.ID) == 0 {
		return nil, model.ErrInvalidCategory
	}
	f := &model.Feed{ID: id, Category: cat}
	_, err := m.req().GetEntities(ctx, f)()
	if _, ok := err.(firestorm.NotFoundError); ok {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (m FeedModel) GetAll(ctx context.Context, cat *model.Category) ([]model.Feed, error) {
	if cat == nil || len(cat.ID) == 0 {
		return nil, model.ErrInvalidCategory
	}
	var feeds []model.Feed
	q := m.req().ToCollection(model.Feed{Category: cat}).Query
	if err := m.req().QueryEntities(ctx, q, &feeds)(); err != nil {
		return nil, err
	}
	return feeds, nil
}

func (m FeedModel) Delete(ctx context.Context, f *model.Feed) error {
	if f == nil {
		return model.ErrInvalidFeed
	}
	if f.Category == nil || len(f.Category.ID) == 0 {
		return model.ErrInvalidCategory
	}
	return m.req().DeleteEntities(ctx, f)()
}

// req is a shortcut to firestorm.FSClient.NewRequest().
func (m FeedModel) req() *firestorm.Request {
	return m.fsc.NewRequest()
}
