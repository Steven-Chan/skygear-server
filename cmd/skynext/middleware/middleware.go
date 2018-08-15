package middleware

import "net/http"

type MiddlewareFunc func(http.Handler) http.Handler
type Middleware interface {
	Handle(http.Handler) http.Handler
}

func ChainMiddleware(ms ...MiddlewareFunc) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		for i := len(ms) - 1; i >= 0; i-- {
			next = ms[i](next)
		}
		return next
	}
}

type NoOpMiddleware struct{}

func (m NoOpMiddleware) Handle(next http.Handler) http.Handler {
	return next
}
