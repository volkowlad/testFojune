package patch

import (
	"context"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"testFojune/internal/errlog"
	"testFojune/internal/http-server/api/response"
)

type Request struct {
	WalletID string `json:"wallet_id"`
	Amount   int    `json:"amount"`
}

type Response struct {
	response.Response
	Balance int `json:"balance"`
}

type UpdaterWallet interface {
	UpdateWallet(ctx context.Context, uuid string, amount int) (int, error)
}

func New(log *slog.Logger, updaterWallet UpdaterWallet) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With("request_id", middleware.GetReqID(r.Context()))

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.With("failed to parse request", errlog.Err(err))

			render.JSON(w, r, response.Error("failed to parse request"))

			return
		}

		log.Info("received request", slog.Any("request", req))

		ctx := r.Context()
		balance, err := updaterWallet.UpdateWallet(ctx, req.WalletID, req.Amount)
		if err != nil {
			log.With("failed to update wallet", errlog.Err(err))

			render.JSON(w, r, response.Error("failed to update wallet"))

			return
		}

		render.JSON(w, r, Response{
			Response: response.OK(),
			Balance:  balance,
		})
	}
}
