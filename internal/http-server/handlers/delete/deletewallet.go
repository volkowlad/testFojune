package deletewallet

import (
	"context"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"testFojune/internal/errlog"
	"testFojune/internal/http-server/api/response"
)

type Request struct {
	WalletID string `json:"wallet_id"`
}

//go:generate mockgen -source=deletewallet.go -destination=mocks/deletemock.go
type DeleterWallet interface {
	DeleteWallet(ctx context.Context, uuid string) error
}

func New(log *slog.Logger, walletDeleter DeleterWallet) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to parse request: %w", errlog.Err(err))

			render.JSON(w, r, response.Error("failed to parse request"))

			return
		}

		ctx := r.Context()
		err = walletDeleter.DeleteWallet(ctx, req.WalletID)
		if err != nil {
			log.Error("failed to delete wallet: %w", errlog.Err(err))
		}

		render.JSON(w, r, response.OK())
	}
}
