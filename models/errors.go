package models

import "errors"

var (
	ErrEmailTaken = errors.New("models: email address is already taken")
	ErrNotFound   = errors.New("models: resource could not be found")
	ErrTitleTaken = errors.New("models: gallery title is already taken")
)
