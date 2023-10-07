package models

import "errors"

var (
	ErrNoGalleryFound = errors.New("models: no gallery found")
	ErrNoGalleries    = errors.New("models: user has no galleries")
	ErrEmailTaken     = errors.New("models: email address is already in use")
)
