package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
)

const CategoryCollection = "Categories"

//var CategoryNotFound = errors.New("user not found")

type CategoryStore struct {
	client *firestore.Client
}

func NewCategoryStore(client *firestore.Client) *CategoryStore {
	return &CategoryStore{client: client}
}

func (us *CategoryStore) Create(ctx context.Context, u interface{}) (string, error) {
	doc, _, err := us.client.Collection(CategoryCollection).Add(ctx, u)
	if err != nil {
		return "", err
	}
	return doc.ID, nil
}

func (us *CategoryStore) Get(ctx context.Context, id string, u interface{}) error {
	snap, err := us.client.Collection(CategoryCollection).Doc(id).Get(ctx)
	if err != nil {
		return err
	}
	if err := snap.DataTo(u); err != nil {
		return err
	}
	return nil
}

func (us *CategoryStore) Delete(ctx context.Context, id string) error {
	if _, err := us.client.Collection(CategoryCollection).Doc(id).Delete(ctx); err != nil {
		return err
	}
	return nil
}
