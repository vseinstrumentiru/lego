package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
)

func combine(next http.Handler, m ...mux.MiddlewareFunc) http.Handler {
	if len(m) == 0 {
		return next
	}

	return m[0](combine(next, m[1:]...))
}

func Combine(mv ...mux.MiddlewareFunc) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return combine(next, mv...)
	}
}
