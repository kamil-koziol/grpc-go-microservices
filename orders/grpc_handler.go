package main

import (
	"context"
	"log"

	pb "github.com/kamil-koziol/common/api"
	"google.golang.org/grpc"
)

type grpcHandler struct {
	pb.UnimplementedOrderServiceServer

	service OrdersService
}

func NewGRPCHandler(grpcServer *grpc.Server, service OrdersService) {
	handler := &grpcHandler{
		service: service,
	}
	pb.RegisterOrderServiceServer(grpcServer, handler)
}

func (h *grpcHandler) CreateOrder(ctx context.Context, payload *pb.CreateOrderRequest) (*pb.Order, error) {
	log.Printf("New order received! %v", payload)
	return &pb.Order{
		ID:         "42",
		CustomerID: payload.CustomerId,
	}, nil
}
