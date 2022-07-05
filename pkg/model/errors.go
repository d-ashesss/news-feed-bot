package model

import "errors"

var ErrInvalidSubscriber = errors.New("invalid subscriber")
var ErrInvalidCategory = errors.New("invalid category")
var ErrNoUpdatesAvailable = errors.New("no update Available")
