package ratelimiter

type Option func(*limiter)

func WithMaxConcurrentConnections(maxConcurrentConnections int) Option {
	return func(l *limiter) {
		l.maxLimit = maxConcurrentConnections
	}
}
