package change

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
	WalletID      string `json:"wallet_id"`
	OperationType string `json:"operation_type"`
	Amount        int    `json:"amount"`
}

type Response struct {
	response.Response
	Balance int `json:"balance"`
}

//go:generate mockgen -source=change.go -destination=mokcs/ckengerMock.go
type ChangerWallet interface {
	ChangeWallet(ctx context.Context, amount int, uuid string, action string) (int, error)
}

func New(log *slog.Logger, walletChanger ChangerWallet) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With("request_id", middleware.GetReqID(r.Context()))

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to parse request", errlog.Err(err))

			render.JSON(w, r, response.Error("failed to parse request"))

			return
		}

		log.Info("received request", slog.Any("request", req))

		ctx := r.Context()
		balance, err := walletChanger.ChangeWallet(ctx, req.Amount, req.WalletID, req.OperationType)
		if err != nil {
			log.Error("failed to change balance", errlog.Err(err))

			render.JSON(w, r, response.Error("failed to change balance"))

			return
		}

		render.JSON(w, r, Response{
			Response: response.OK(),
			Balance:  balance,
		})
	}
}
