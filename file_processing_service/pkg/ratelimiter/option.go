package ratelimiter

import "time"

type Option func(*limiter)

// WithMaxConcurrentConnections is an option that sets the maximum number of concurrent connections.
func WithMaxConcurrentConnections(maxConcurrentConnections int) Option {
	return func(l *limiter) {
		l.maxLimit = maxConcurrentConnections
	}
}

// WithLimitingByUserID is an option that enables limiting by user ID.
func WithLimitingByUserID() Option {
	return func(l *limiter) {
		l.isLimitedByID = true
	}
}

// WithDuration is an option that sets the duration for which the limit is applied.
func WithDuration(duration time.Duration) Option {
	return func(l *limiter) {
		l.duration = duration
	}
}

// WithLimiterType is an option that sets the type of limiter.
func WithLimiterType(limiterType limiterType) Option {
	return func(l *limiter) {
		l.limiterType = limiterType
	}
}
