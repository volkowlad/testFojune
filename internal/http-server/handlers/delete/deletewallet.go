package deletewallet

import (
	"github.com/go-chi/render"
	"github.com/gofrs/uuid"
	"log/slog"
	"net/http"
	"testFojune/internal/errlog"
	"testFojune/internal/http-server/api/response"
)

type Request struct {
	WalletID uuid.UUID `json:"wallet_id"`
}

type DeleterWallet interface {
	DeleteWallet(uuid uuid.UUID) error
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

		err = walletDeleter.DeleteWallet(req.WalletID)
		if err != nil {
			log.Error("failed to delete wallet: %w", errlog.Err(err))
		}

		render.JSON(w, r, response.OK())
	}
}
