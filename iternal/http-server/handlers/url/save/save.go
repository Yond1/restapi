package save

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"restapi/iternal/lib/api/response"
	"restapi/iternal/storage"
	"strings"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}

type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = slog.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request", err)
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.String("url", req.URL), slog.String("alias", req.Alias))

		if err := validator.New().Struct(req); err != nil {
			log.Error("failed to validate request", err)

			render.JSON(w, r, response.Error("failed to validate request"))

			return
		}

		alias := req.Alias
		if alias == "" {
			alias = strings.TrimPrefix(req.URL, "https://")
		}

		id, err := urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrorExists) {
			log.Info("url already exists", err)

			render.JSON(w, r, response.Error("url already exists"))

			return
		}
		if err != nil {
			log.Error("failed to save url", err)

			render.JSON(w, r, response.Error("failed to save url"))

			return
		}
		log.Info("url saved", slog.Int64("id", id))

		render.JSON(w, r, Response{
			Response: response.Ok(),
			Alias:    alias,
		})
	}
}
