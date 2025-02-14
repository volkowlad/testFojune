package save

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"testFojune/internal/errlog"
	"testFojune/internal/http-server/api/response"
)

type Request struct {
	Amount int `json:"amount"`
}

type Response struct {
	response.Response
	Uuid string `json:"wallet_id,omitempty"`
}

type WalletSaver interface {
	SaveWallet(amount int) (string, error)
}

func New(log *slog.Logger, saverWallet WalletSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With("request_id", middleware.GetReqID(r.Context()))

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to parse request: %s", errlog.Err(err))

			render.JSON(w, r, response.Error("failed to parse request"))

			return
		}

		log.Info("received request", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			log.Error("failed to validate request: %s", errlog.Err(err))

			render.JSON(w, r, response.Error("failed to validate request"))

			return
		}

		uuid, err := saverWallet.SaveWallet(req.Amount)
		if err != nil {
			log.Error("failed to save wallet: %s", errlog.Err(err))

			render.JSON(w, r, response.Error("failed to save wallet"))

			return
		}

		render.JSON(w, r, Response{
			Response: response.OK(),
			Uuid:     uuid,
		})
	}
}
