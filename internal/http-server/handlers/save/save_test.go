package save_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testFojune/internal/http-server/handlers/save"
	mock_save "testFojune/internal/http-server/handlers/save/mocks"
	"testing"
)

func TestSave(t *testing.T) {
	cases := []struct {
		result         string
		inputAmount    int
		inputCtx       context.Context
		expectedStatus int
		expectedBody   string
		mockError      error
	}{
		{
			result:         "Success",
			inputAmount:    1000,
			inputCtx:       context.Background(),
			expectedStatus: http.StatusOK,
			expectedBody:   "{\"status\":\"OK\",\"wallet_id\":\"ok\"}\n",
		},
		{
			result:         "empty amount",
			inputCtx:       context.Background(),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "{\"status\":\"ERROR\",\"error\":\"failed to save wallet\"}\n",
			mockError:      errors.New("empty amount"),
		},
		{
			result:         "no connected",
			inputAmount:    1000,
			inputCtx:       context.Background(),
			expectedStatus: http.StatusServiceUnavailable,
			expectedBody:   "{\"status\":\"ERROR\",\"error\":\"failed to save wallet\"}\n",
			mockError:      errors.New("no connected"),
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.result, func(t *testing.T) {
			t.Parallel()

			c := gomock.NewController(t)
			defer c.Finish()

			var walletSaverMock = mock_save.NewMockWalletSaver(c)
			//if testCase.respError == "" || testCase.mockError != nil {
			//	walletSaverMock.EXPECT().SaveWallet(testCase.inputAmount).Return(1, testCase.mockError)
			//}
			walletSaverMock.EXPECT().SaveWallet(testCase.inputCtx, testCase.inputAmount).Return("ok", testCase.mockError)
			saverHandler := save.New(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{})), walletSaverMock)

			input := fmt.Sprintf(`{"amount": %d}`, testCase.inputAmount)

			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/wallet/save", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			// записываем ответ handler
			saverHandler.ServeHTTP(w, req)

			require.Equal(t, w.Code, http.StatusOK)
			require.Equal(t, w.Body.String(), testCase.expectedBody)
		})
	}
}
