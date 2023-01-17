package entity

import "errors"

var (
	ErrNotFound       = errors.New("file not found")
	ErrInvalidEntity  = errors.New("invalid entity")
	ErrNotEnoughSpace = errors.New("not enough space")
	ErrAlreadyExists  = errors.New("already exists")
	// ErrTooManyRequests = errors.New("too many requests")
)
