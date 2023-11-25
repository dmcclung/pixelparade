package models

import (
	"errors"
	"fmt"
)

var (
	ErrNoGalleryFound    = errors.New("models: no gallery found")
	ErrNoGalleries       = errors.New("models: user has no galleries")
	ErrEmailTaken        = errors.New("models: email address is already in use")
	ErrImageNotFound     = errors.New("models: image not found")
	ErrImageMetaNotFound = errors.New("models: image meta not found")
)

type FileError struct {
	Issue string
}

func (fe FileError) Error() string {
	return fmt.Sprintf("invalid file: %s", fe.Issue)
}
