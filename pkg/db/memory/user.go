package memory

import (
	"context"
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"sync"
	"time"
)

// UserStore is a dummy implementation of model.UserStore interface for testing purposes.
type UserStore struct {
	sync.Mutex
	users map[string]interface{}
}

func NewUserStore() *UserStore {
	return &UserStore{users: map[string]interface{}{}}
}

func (us *UserStore) Create(_ context.Context, u interface{}) (string, error) {
	us.Lock()
	defer us.Unlock()

	rand.Seed(time.Now().UnixNano())
	id := rand.Int63()
	idStr := strconv.FormatInt(id, 16)

	us.users[idStr] = cloneObject(u)
	return idStr, nil
}

func (us *UserStore) Get(_ context.Context, id string, u interface{}) error {
	if o, ok := us.users[id]; ok {
		copyObject(o, u)
		return nil
	}
	return fmt.Errorf("user not found")
}

func (us *UserStore) GetByTelegramId(_ context.Context, telegramId int, u interface{}) (string, error) {
	for id, o := range us.users {
		if getField(o, "TelegramId").(int) == telegramId {
			copyObject(o, u)
			return id, nil
		}
	}
	return "", fmt.Errorf("user not found")
}

func (us *UserStore) Delete(_ context.Context, id string) error {
	us.Lock()
	defer us.Unlock()

	delete(us.users, id)
	return nil
}

func getField(o interface{}, f string) interface{} {
	return reflect.ValueOf(o).Elem().FieldByName(f).Interface()
}

func cloneObject(src interface{}) interface{} {
	dst := reflect.New(reflect.TypeOf(src).Elem()).Interface()
	copyObject(src, dst)
	return dst
}

func copyObject(src, dst interface{}) {
	dstVal := reflect.ValueOf(dst).Elem()
	srcVal := reflect.ValueOf(src).Elem()
	for i := 0; i < srcVal.NumField(); i++ {
		if dstField := dstVal.Field(i); dstField.CanSet() {
			dstField.Set(srcVal.Field(i))
		}
	}
	return
}
