package delete

import (
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"restapi/iternal/lib/api/response"
)

type UrlRemover interface {
	DeleteURL(string) error
}

func New(log *slog.Logger, UrlRemover UrlRemover) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		op := "handlers.url.delete.New"

		log.Info(op, slog.String("url", chi.URLParam(r, "url")))

		err := UrlRemover.DeleteURL(chi.URLParam(r, "alias"))
		if err != nil {
			response.Error("internal server error")
			return
		}

		log.Info("url deleted", slog.String("url", chi.URLParam(r, "alias")))

		w.WriteHeader(http.StatusAccepted)
	}
}
