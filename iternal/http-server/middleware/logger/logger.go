package logger

import (
	"log/slog"
	"net/http"
)

func New(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log = log.With(slog.String("middleware", "logger"))

		log.Info("logger middleware created")

		fn := func(w http.ResponseWriter, r *http.Request) {
			log.Info("request started")

			defer func() {
				log.Info("request finished")
			}()

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
