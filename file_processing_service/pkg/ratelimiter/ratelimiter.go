package ratelimiter

import (
	"errors"
	"sync"
)

var ErrLimitExceeded = errors.New("request limit exceeded")

const DefaultMaxLimit = 10

type limiter struct {
	mu         sync.Mutex
	maxLimit   int
	concurrent int
}

func (l *limiter) Acquire() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.concurrent >= l.maxLimit {
		return ErrLimitExceeded
	}

	l.concurrent++
	return nil
}

func (l *limiter) Release() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.concurrent--
}

func NewLimiter(options ...Option) *limiter {
	limiter := &limiter{
		maxLimit: DefaultMaxLimit,
	}
	for _, option := range options {
		option(limiter)
	}
	return limiter
}
