package model

import "errors"

var ErrInvalidSubscriber = errors.New("invalid subscriber")
var ErrInvalidCategory = errors.New("invalid category")
var ErrInvalidCategoryName = errors.New("invalid category name")
var ErrCategoryNotFound = errors.New("category no found")
var ErrNoUpdatesAvailable = errors.New("no update Available")
