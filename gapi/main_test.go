package gapi

import (
	"testing"
	"time"

	db "github.com/juker1141/simplebank/db/sqlc"
	"github.com/juker1141/simplebank/util"
	"github.com/juker1141/simplebank/worker"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store, taskDistributor worker.TaskDistribtor) *Server {
	config := util.Config{
		TokenSymmetricKey: 	 util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store, taskDistributor)
	require.NoError(t, err)

	return server
}