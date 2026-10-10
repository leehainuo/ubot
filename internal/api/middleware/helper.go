package middleware

import (
	"net/http"
	"slices"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func Chain(handler http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, mw := range slices.Backward(middlewares) {
		handler = mw(handler)
	}
	return handler
}
