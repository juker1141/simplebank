package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/juker1141/simplebank/db/mock"
	db "github.com/juker1141/simplebank/db/sqlc"
	"github.com/juker1141/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestCreateTransfer(t *testing.T) {
	amount := int64(10)

	account1 := randomAccount()
	account2 := randomAccount()
	account3 := randomAccount()

	account1.Currency = util.USD
	account2.Currency = util.USD
	account3.Currency = util.EUR

	testCases := []struct{
		name string
		body gin.H
		buildStubs func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id": account2.ID,
				"amount": amount,
				"currency": util.USD,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
					Times(1).
					Return(account1, nil)
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account2.ID)).
					Times(1).
					Return(account2, nil)

				arg := db.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID: account2.ID,
					Amount: amount,
				}
				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Eq(arg)).
					Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		store := mockdb.NewMockStore(ctrl)
		tc.buildStubs(store)

		// start test server and send request
		server := NewServer(store)
		recorder := httptest.NewRecorder()

		jsonData, err := json.Marshal(tc.body)
		require.NoError(t, err)

		url := "/transfers"
		result, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
		require.NoError(t, err)

		server.router.ServeHTTP(recorder, result)
		tc.checkResponse(t, recorder)
	}
}