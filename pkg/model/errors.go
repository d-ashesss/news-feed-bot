package model

import "errors"

var ErrInvalidSubscriber = errors.New("invalid subscriber")
var ErrInvalidSubscriberID = errors.New("invalid subscriber ID")
var ErrInvalidUpdate = errors.New("invalid update")
var ErrInvalidFeed = errors.New("invalid feed")
var ErrInvalidCategory = errors.New("invalid category")
var ErrInvalidCategoryName = errors.New("invalid category name")
var ErrNotFound = errors.New("not found")
var ErrNoUpdates = errors.New("no update Available")
