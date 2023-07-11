package gapi

import (
	"fmt"

	db "github.com/juker1141/simplebank/db/sqlc"
	"github.com/juker1141/simplebank/pb"
	"github.com/juker1141/simplebank/token"
	"github.com/juker1141/simplebank/util"
	"github.com/juker1141/simplebank/worker"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	config util.Config
	store db.Store
	tokenMaker token.Maker
	taskDistributor worker.TaskDistribtor
}

// NewServer creates a new gRPC server.
func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistribtor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config: config,
		store: store,
		tokenMaker: tokenMaker,
		taskDistributor: taskDistributor,
	}

	return server, nil
}