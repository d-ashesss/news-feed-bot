package model

import (
	"context"
	"errors"
	"github.com/d-ashesss/news-feed-bot/mocks"
	"github.com/stretchr/testify/mock"
	"testing"
)

var (
	testUserError = errors.New("test error")
)

const (
	testUserId         = "test-ID"
	testUserTelegramId = 1111
)

func TestUserModel_Create(t *testing.T) {
	t.Run("user exists", func(t *testing.T) {
		db := &mocks.UserStore{}

		m := NewUserModel(db)
		u := NewUser(testUserTelegramId)
		u.Id = "new-ID"

		if err := m.Create(context.Background(), u); err == nil {
			t.Errorf("Create() expected error")
		}
	})

	t.Run("create error", func(t *testing.T) {
		db := &mocks.UserStore{}
		db.On("Create", mock.Anything, mock.Anything).Return("", testUserError)

		m := NewUserModel(db)
		u := NewUser(testUserTelegramId)

		if err := m.Create(context.Background(), u); err != testUserError {
			t.Errorf("Create() got error %q, want %q", err, testUserError)
		}
	})

	t.Run("user created", func(t *testing.T) {
		db := &mocks.UserStore{}
		db.On("Create", mock.Anything, mock.Anything).Return(func(ctx context.Context, o interface{}) string {
			return testUserId
		}, nil)

		m := NewUserModel(db)
		u := NewUser(testUserTelegramId)

		if err := m.Create(context.Background(), u); err != nil {
			t.Errorf("Create() unexpected error = %v", err)
		}
		if u.Id != testUserId {
			t.Errorf("Create() got user ID %q, want %q", u.Id, testUserId)
		}
	})
}

func TestUserModel_Get(t *testing.T) {
	t.Run("user not found", func(t *testing.T) {
		db := &mocks.UserStore{}
		db.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(testUserError)
		m := NewUserModel(db)

		u, err := m.Get(context.Background(), testUserId)
		if err != testUserError {
			t.Errorf("Get() got error %q, want %q", err, testUserError)
		}
		if u != nil {
			t.Errorf("Get() returned unexpected user")
		}
	})

	t.Run("user found", func(t *testing.T) {
		db := &mocks.UserStore{}
		db.On("Get", mock.Anything, mock.MatchedBy(func(id string) bool {
			return id == testUserId
		}), mock.Anything).Return(nil)
		m := NewUserModel(db)

		u, err := m.Get(context.Background(), testUserId)
		if err != nil {
			t.Errorf("Get() unexpected error = %v", err)
		}
		if u.Id != testUserId {
			t.Errorf("Get() got user with ID %q, want %q", u.Id, testUserId)
		}
	})
}

func TestUserModel_GetByTelegramId(t *testing.T) {
	t.Run("user found", func(t *testing.T) {
		db := &mocks.UserStore{}
		db.On("GetByTelegramId", mock.Anything, mock.MatchedBy(func(id int) bool {
			return id == testUserTelegramId
		}), mock.Anything).Return(testUserId, nil)
		m := NewUserModel(db)

		u, err := m.GetByTelegramId(context.Background(), testUserTelegramId)
		if err != nil {
			t.Errorf("GetByTelegramId() unexpected error = %v", err)
		}
		if u.Id != testUserId {
			t.Errorf("GetByTelegramId() got user with ID %q, want %q", u.Id, testUserId)
		}
	})
}

func TestUserModel_Delete(t *testing.T) {
	t.Run("user deleted", func(t *testing.T) {
		db := &mocks.UserStore{}
		db.On("Delete", mock.Anything, mock.MatchedBy(func(id string) bool {
			return id == testUserId
		})).Return(nil)
		m := NewUserModel(db)
		u := NewUser(testUserTelegramId)
		u.Id = testUserId
		if err := m.Delete(context.Background(), u); err != nil {
			t.Errorf("Delete() unexpected error = %v", err)
		}
	})
}
