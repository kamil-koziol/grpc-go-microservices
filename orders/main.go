package main

import (
	"context"
	"log"
	"net"

	"github.com/kamil-koziol/common"
	"google.golang.org/grpc"
)

var (
	grpcAddr = common.EnvString("GRPC_ADDR", "localhost:3000")
)

func main() {
	grpcServer := grpc.NewServer()
	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	store := NewStore()
	svc := NewService(store)
	NewGRPCHandler(grpcServer, svc)

	svc.CreateOrder(context.Background())

	log.Printf("GRPC Server Started at %s", grpcAddr)
	if err := grpcServer.Serve(l); err != nil {
		log.Fatal(err.Error())
	}
}
