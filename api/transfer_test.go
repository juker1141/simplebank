package api

import (
	"testing"

	"github.com/golang/mock/gomock"
	mockdb "github.com/juker1141/simplebank/db/mock"
)

func TestCreateTransfer(t *testing.T) {
	account1 := randomAccount()
	account2 := randomAccount()
	account3 := randomAccount()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	store.EXPECT().CreateTransfer(gomock.Any(), ).Times(1).Return()
}