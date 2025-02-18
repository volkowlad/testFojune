package change

import (
	"bytes"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	mock_change "testFojune/internal/http-server/handlers/change/mokcs"
	"testing"
)

func TestChange(t *testing.T) {
	cases := []struct {
		result         string
		inputAction    string
		inputAmount    int
		inputWalletid  string
		expectedStatus int
		expectedBody   string
		respError      string
		mockError      error
	}{
		{
			result:         "success",
			inputAction:    "dePosit",
			inputAmount:    666,
			inputWalletid:  uuid.New().String(),
			expectedStatus: http.StatusOK,
			expectedBody:   "{\"status\":\"OK\",\"balance\":1}\n",
		},
	}

	for _, tt := range cases {
		t.Run(tt.respError, func(t *testing.T) {
			t.Parallel()

			c := gomock.NewController(t)
			defer c.Finish()

			changerMock := mock_change.NewMockChangerWallet(c)
			changerMock.EXPECT().ChangeWallet(tt.inputAmount, tt.inputWalletid, tt.inputAction).Return(1, tt.mockError)
			changeHandler := New(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{})), changerMock)

			inputBody := fmt.Sprintf(`{"wallet_id": "%s", "operation_type": "%s", "amount": %d}`, tt.inputWalletid, tt.inputAction, tt.inputAmount)

			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/wallet", bytes.NewReader([]byte(inputBody)))
			require.NoError(t, err)

			changeHandler.ServeHTTP(w, req)

			require.Equal(t, tt.expectedStatus, w.Code)
			require.Equal(t, tt.expectedBody, w.Body.String())

		})
	}
}
