package model

import (
	"context"
	"fmt"
)

// User represents a user entity
type User struct {
	ID         string `firestore:",omitempty"`
	TelegramID int
}

// NewUser initializes new User
func NewUser(TelegramID int) *User {
	return &User{
		TelegramID: TelegramID,
	}
}

// UserStore is an interface wrapper for a DB engine
type UserStore interface {
	Create(ctx context.Context, u interface{}) (string, error)
	Get(ctx context.Context, id string, u interface{}) error
	GetByTelegramID(ctx context.Context, telegramID int, u interface{}) (string, error)
	Delete(ctx context.Context, id string) error
}

// UserModel data model for User
type UserModel struct {
	db UserStore
}

// NewUserModel initializes UserModel
func NewUserModel(db UserStore) *UserModel {
	return &UserModel{db: db}
}

// Create saves a new User into the DB.
func (m UserModel) Create(ctx context.Context, u *User) (err error) {
	if len(u.ID) > 0 {
		return fmt.Errorf("create user: provided user is not new")
	}
	u.ID, err = m.db.Create(ctx, u)
	return
}

// Get retrieves User from the DB by ID.
func (m UserModel) Get(ctx context.Context, id string) (*User, error) {
	var u User
	if err := m.db.Get(ctx, id, &u); err != nil {
		return nil, err
	}
	u.ID = id
	return &u, nil
}

// GetByTelegramID retrieves User from the DB by Telegram ID.
func (m UserModel) GetByTelegramID(ctx context.Context, telegramID int) (*User, error) {
	var u User
	id, err := m.db.GetByTelegramID(ctx, telegramID, &u)
	if err != nil {
		return nil, err
	}
	u.ID = id
	return &u, nil
}

// Delete deletes User from the DB.
func (m UserModel) Delete(ctx context.Context, u *User) error {
	return m.db.Delete(ctx, u.ID)
}
