package get

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"testFojune/internal/errlog"
	"testFojune/internal/http-server/api/response"
)

type Response struct {
	response.Response
	Balance int `json:"balance"`
}

type BalanceGetter interface {
	GetWallet(uuid string) (int, error)
}

func New(log *slog.Logger, getterBalance BalanceGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With("request_id", middleware.GetReqID(r.Context()))

		walletid := chi.URLParam(r, "uuid")
		if walletid == "" {
			log.Info("no uuid provided")

			render.JSON(w, r, response.Error("invalid request"))

			return
		}

		balance, err := getterBalance.GetWallet(walletid)
		if errors.Is(err, errors.New("wallet not found")) {
			log.Info("wallet not found", "uuid", walletid)

			render.JSON(w, r, response.Error("wallet not found"))

			return
		}
		if err != nil {
			log.Error("failed to get wallet", errlog.Err(err))

			render.JSON(w, r, response.Error("failed to get wallet"))

			return
		}

		log.Info("got balance", slog.Int("balance", balance))

		render.JSON(w, r, Response{
			Response: response.OK(),
			Balance:  balance,
		})
	}
}
