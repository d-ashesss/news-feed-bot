package firestore

import (
	"cloud.google.com/go/firestore"
	fst "cloud.google.com/go/firestore"
	"context"
	"github.com/d-ashesss/news-feed-bot/pkg/model"
	"github.com/jschoedt/go-firestorm"
)

type CategoryModel struct {
	client *firestorm.FSClient
}

// NewCategoryModel initializes Firestore implementation of model.CategoryModel.
func NewCategoryModel(c *fst.Client) model.CategoryModel {
	return &CategoryModel{client: firestorm.New(c, "ID", "")}
}

func (m CategoryModel) Create(ctx context.Context, c *model.Category) (string, error) {
	if err := m.req().CreateEntities(ctx, c)(); err != nil {
		return "", err
	}
	return m.req().GetID(c), nil
}

func (m CategoryModel) Get(ctx context.Context, id string) (*model.Category, error) {
	c := &model.Category{ID: id}
	if _, err := m.req().GetEntities(ctx, c)(); err != nil {
		return nil, err
	}
	return c, nil
}

func (m CategoryModel) GetAll(ctx context.Context) ([]model.Category, error) {
	cats := make([]model.Category, 0)
	q := m.req().ToCollection(model.Category{}).OrderBy("name", firestore.Asc)
	if err := m.req().QueryEntities(ctx, q, &cats)(); err != nil {
		return nil, err
	}
	return cats, nil
}

func (m CategoryModel) Delete(ctx context.Context, c *model.Category) error {
	return m.req().DeleteEntities(ctx, c)()
}

func (m CategoryModel) req() *firestorm.Request {
	return m.client.NewRequest()
}
