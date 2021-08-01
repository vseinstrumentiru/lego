package middleware

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/vseinstrumentiru/lego/v2/log"
)

func RecoverHandlerMiddleware(handler log.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					handler.WithContext(r.Context()).
						WithFields(map[string]interface{}{"status": "panic"}).
						Notify(err)

					w.WriteHeader(http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
