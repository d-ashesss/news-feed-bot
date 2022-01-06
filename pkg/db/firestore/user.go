package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"google.golang.org/api/iterator"
)

const UserCollection = "Users"

const UserFieldTelegramID = "TelegramID"

var UserNotFound = errors.New("user not found")

type UserStore struct {
	client *firestore.Client
}

func NewUserStore(client *firestore.Client) *UserStore {
	return &UserStore{client: client}
}

func (us *UserStore) Create(ctx context.Context, u interface{}) (string, error) {
	doc, _, err := us.client.Collection(UserCollection).Add(ctx, u)
	if err != nil {
		return "", err
	}
	return doc.ID, nil
}

func (us *UserStore) Get(ctx context.Context, id string, u interface{}) error {
	snap, err := us.client.Collection(UserCollection).Doc(id).Get(ctx)
	if err != nil {
		return err
	}
	if err := snap.DataTo(u); err != nil {
		return err
	}
	return nil
}

func (us *UserStore) GetByTelegramID(ctx context.Context, telegramID int, u interface{}) (string, error) {
	q := us.client.Collection(UserCollection).
		Where(UserFieldTelegramID, "==", telegramID).
		Limit(1)
	iter := q.Documents(ctx)
	snap, err := iter.Next()
	if err == iterator.Done {
		return "", UserNotFound
	}
	if err != nil {
		return "", err
	}
	if err := snap.DataTo(u); err != nil {
		return "", err
	}
	return snap.Ref.ID, nil
}

func (us *UserStore) Delete(ctx context.Context, id string) error {
	if _, err := us.client.Collection(UserCollection).Doc(id).Delete(ctx); err != nil {
		return err
	}
	return nil
}
