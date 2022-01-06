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
	testUserID         = "test-ID"
	testUserTelegramID = 1111
)

func TestUserModel_Create(t *testing.T) {
	t.Run("user exists", func(t *testing.T) {
		db := &mocks.UserStore{}

		m := NewUserModel(db)
		u := NewUser(testUserTelegramID)
		u.ID = "new-ID"

		if err := m.Create(context.Background(), u); err == nil {
			t.Errorf("Create() expected error")
		}
	})

	t.Run("create error", func(t *testing.T) {
		db := &mocks.UserStore{}
		db.On("Create", mock.Anything, mock.Anything).Return("", testUserError)

		m := NewUserModel(db)
		u := NewUser(testUserTelegramID)

		if err := m.Create(context.Background(), u); err != testUserError {
			t.Errorf("Create() got error %q, want %q", err, testUserError)
		}
	})

	t.Run("user created", func(t *testing.T) {
		db := &mocks.UserStore{}
		db.On("Create", mock.Anything, mock.Anything).Return(func(ctx context.Context, o interface{}) string {
			return testUserID
		}, nil)

		m := NewUserModel(db)
		u := NewUser(testUserTelegramID)

		if err := m.Create(context.Background(), u); err != nil {
			t.Errorf("Create() unexpected error = %v", err)
		}
		if u.ID != testUserID {
			t.Errorf("Create() got user ID %q, want %q", u.ID, testUserID)
		}
	})
}

func TestUserModel_Get(t *testing.T) {
	t.Run("user not found", func(t *testing.T) {
		db := &mocks.UserStore{}
		db.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(testUserError)
		m := NewUserModel(db)

		u, err := m.Get(context.Background(), testUserID)
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
			return id == testUserID
		}), mock.Anything).Return(nil)
		m := NewUserModel(db)

		u, err := m.Get(context.Background(), testUserID)
		if err != nil {
			t.Errorf("Get() unexpected error = %v", err)
		}
		if u.ID != testUserID {
			t.Errorf("Get() got user with ID %q, want %q", u.ID, testUserID)
		}
	})
}

func TestUserModel_GetByTelegramID(t *testing.T) {
	t.Run("user found", func(t *testing.T) {
		db := &mocks.UserStore{}
		db.On("GetByTelegramID", mock.Anything, mock.MatchedBy(func(id int) bool {
			return id == testUserTelegramID
		}), mock.Anything).Return(testUserID, nil)
		m := NewUserModel(db)

		u, err := m.GetByTelegramID(context.Background(), testUserTelegramID)
		if err != nil {
			t.Errorf("GetByTelegramID() unexpected error = %v", err)
		}
		if u.ID != testUserID {
			t.Errorf("GetByTelegramID() got user with ID %q, want %q", u.ID, testUserID)
		}
	})
}

func TestUserModel_Delete(t *testing.T) {
	t.Run("user deleted", func(t *testing.T) {
		db := &mocks.UserStore{}
		db.On("Delete", mock.Anything, mock.MatchedBy(func(id string) bool {
			return id == testUserID
		})).Return(nil)
		m := NewUserModel(db)
		u := NewUser(testUserTelegramID)
		u.ID = testUserID
		if err := m.Delete(context.Background(), u); err != nil {
			t.Errorf("Delete() unexpected error = %v", err)
		}
	})
}
