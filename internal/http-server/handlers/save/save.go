package save

import (
	"context"
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

//go:generate mockgen -source=save.go -destination=mocks/saverMock.go
type WalletSaver interface {
	SaveWallet(ctx context.Context, amount int) (string, error)
}

func New(log *slog.Logger, saverWallet WalletSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log = log.With("request_id", middleware.GetReqID(ctx))

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

		uuid, err := saverWallet.SaveWallet(ctx, req.Amount)
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
