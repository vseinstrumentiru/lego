package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type app struct{}

func (app) Providers() []interface{} {
	return []interface{}{
		// put here your constructors
	}
}

func (app) ConfigureHTTP(router *mux.Router) {
	router.Path("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Hello LeGo!"))
	})
}
