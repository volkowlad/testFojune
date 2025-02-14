package change

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/gofrs/uuid"
	"log/slog"
	"net/http"
	"testFojune/internal/errlog"
	"testFojune/internal/http-server/api/response"
)

type Request struct {
	WalletID      uuid.UUID `json:"wallet_id"`
	OperationType string    `json:"operation_type"`
	Amount        int       `json:"amount"`
}

type Response struct {
	response.Response
	Balance int `json:"balance"`
}

type ChangerWallet interface {
	ChangeWallet(amount int, uuid uuid.UUID, action string) (int, error)
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

		balance, err := walletChanger.ChangeWallet(req.Amount, req.WalletID, req.OperationType)
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
