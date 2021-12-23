package model

import (
	"context"
	"fmt"
)

// User represents a user entity
type User struct {
	Id         string
	TelegramId int
}

// NewUser initializes new User
func NewUser(TelegramId int) *User {
	return &User{
		TelegramId: TelegramId,
	}
}

// UserStore is an interface wrapper for a DB engine
type UserStore interface {
	Create(ctx context.Context, u interface{}) (string, error)
	Get(ctx context.Context, id string, u interface{}) error
	GetByTelegramId(ctx context.Context, telegramId int, u interface{}) (string, error)
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
	if len(u.Id) > 0 {
		return fmt.Errorf("create user: provided user is not new")
	}
	u.Id, err = m.db.Create(ctx, u)
	return
}

// Get retrieves User from the DB by ID.
func (m UserModel) Get(ctx context.Context, id string) (*User, error) {
	var u User
	if err := m.db.Get(ctx, id, &u); err != nil {
		return nil, err
	}
	u.Id = id
	return &u, nil
}

// GetByTelegramId retrieves User from the DB by Telegram ID.
func (m UserModel) GetByTelegramId(ctx context.Context, telegramId int) (*User, error) {
	var u User
	id, err := m.db.GetByTelegramId(ctx, telegramId, &u)
	if err != nil {
		return nil, err
	}
	u.Id = id
	return &u, nil
}

// Delete deletes User from the DB.
func (m UserModel) Delete(ctx context.Context, u *User) error {
	return m.db.Delete(ctx, u.Id)
}
