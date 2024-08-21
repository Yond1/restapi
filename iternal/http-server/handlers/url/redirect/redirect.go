package redirect

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"restapi/iternal/lib/api/response"
)

type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, UrlGetter URLGetter) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		op := "handlers.url.redirect.New"

		log.With(slog.String("op", op)).
			With(slog.String("request_id", middleware.GetReqID(r.Context()))).
			Info("request received")

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			render.JSON(w, r, response.Error("url not found"))
			return
		}

		log.Info("request body decoded", slog.String("url", alias))

		urlRedirect, err := UrlGetter.GetURL(alias)
		if err != nil {
			render.JSON(w, r, response.Error("internal server error"))
			return
		}
		http.Redirect(w, r, urlRedirect, http.StatusTemporaryRedirect)
	}
}
