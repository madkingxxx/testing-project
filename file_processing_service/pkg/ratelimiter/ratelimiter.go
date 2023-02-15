package ratelimiter

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrLimitExceeded      = errors.New("request limit exceeded")
	ErrWrongFunctionUsage = errors.New("wrong function usage by limiter option")
)

const (
	_defaultMaxLimit    = 10
	_defaultDuration    = 1 * time.Second
	_defaultLimiterType = Concurrent

	_defaultKey = "default"
)

type limiter struct {
	mu            sync.Mutex
	maxLimit      int
	concurrent    map[string]int
	isLimitedByID bool
	limiterType   limiterType
	duration      time.Duration
}

type limiterType int

const (
	_ limiterType = iota
	Concurrent
	DurationBased
)

func (l *limiter) Acquire() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.concurrent[_defaultKey] >= l.maxLimit {
		return ErrLimitExceeded
	}

	l.concurrent[_defaultKey]++

	return nil
}

func (l *limiter) AcquireWithID(id string) error {
	if !l.isLimitedByID {
		return ErrWrongFunctionUsage
	}
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.concurrent[id] >= l.maxLimit {
		return fmt.Errorf("%w for - %s", ErrLimitExceeded, id)
	}

	l.concurrent[id]++

	return nil
}

func (l *limiter) Release() {
	l.mu.Lock()
	l.concurrent[_defaultKey]--
	l.mu.Unlock()
}

func (l *limiter) ReleaseWithID(id string) {
	l.mu.Lock()
	l.concurrent[id]--
	l.mu.Unlock()
}

func (l *limiter) Start() {
	if l.limiterType == DurationBased {
		ticker := time.NewTicker(l.duration)
		go func() {
			for range ticker.C {
				l.mu.Lock()
				l.concurrent = make(map[string]int)
				l.mu.Unlock()
			}
		}()
	}
}

func NewLimiter(options ...Option) *limiter {
	limiter := &limiter{
		maxLimit:    _defaultMaxLimit,
		duration:    _defaultDuration,
		limiterType: _defaultLimiterType,
		concurrent:  make(map[string]int),
	}
	for _, option := range options {
		option(limiter)
	}
	limiter.Start()
	return limiter
}
