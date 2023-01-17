package entity

import "errors"

var (
	ErrNotFound       = errors.New("not found")
	ErrInvalidEntity  = errors.New("invalid entity")
	ErrNotEnoughSpace = errors.New("not enough space")
	// ErrTooManyRequests = errors.New("too many requests")
)
