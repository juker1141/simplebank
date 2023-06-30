package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/juker1141/simplebank/api"
	"github.com/juker1141/simplebank/gapi"
	"github.com/juker1141/simplebank/pb"
	"github.com/juker1141/simplebank/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	db "github.com/juker1141/simplebank/db/sqlc"

	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db.", err)
	}

	store := db.NewStore(conn)
	runGrpcServer(config, store)
}

func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	lintener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("cannot create listener")
	}

	log.Printf("start gRPC server at%s", lintener.Addr().String())
	err = grpcServer.Serve(lintener)
	if err != nil {
		log.Fatal("cannot create gRPC server")
	}
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}