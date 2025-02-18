package save_test

import (
	"bytes"
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
		expectedStatus int
		expectedBody   string
		respError      string
		mockError      error
	}{
		{
			result:         "Success",
			inputAmount:    1000,
			expectedStatus: http.StatusOK,
			expectedBody:   "{\"status\":\"OK\",\"wallet_id\":\"ok\"}\n",
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
			walletSaverMock.EXPECT().SaveWallet(testCase.inputAmount).Return("ok", testCase.mockError)
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
