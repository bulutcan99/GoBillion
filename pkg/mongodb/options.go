package mongodb

import "time"

type Option func(db *mongoDB)

func ConnAttempts(attempts int) Option {
	return func(p *mongoDB) {
		p.connAttempts = attempts
	}
}

func ConnTimeout(timeout time.Duration) Option {
	return func(p *mongoDB) {
		p.connTimeout = timeout
	}
}
