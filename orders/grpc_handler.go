package main

import (
	"context"
	"encoding/json"
	"log"

	pb "github.com/kamil-koziol/common/api"
	"github.com/kamil-koziol/common/broker"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
)

type grpcHandler struct {
	pb.UnimplementedOrderServiceServer

	service OrdersService
	channel *amqp.Channel
}

func NewGRPCHandler(grpcServer *grpc.Server, service OrdersService, channel *amqp.Channel) {
	handler := &grpcHandler{
		service: service,
		channel: channel,
	}
	pb.RegisterOrderServiceServer(grpcServer, handler)
}

func (h *grpcHandler) CreateOrder(ctx context.Context, payload *pb.CreateOrderRequest) (*pb.Order, error) {
	log.Printf("New order received! %v", payload)
    o := &pb.Order{
		ID:         "42",
		CustomerID: payload.CustomerId,
	}

    marshalledOrder, err := json.Marshal(o)
    if err != nil {
        return nil, err
    }

    q, err := h.channel.QueueDeclare(broker.OrderCreatedEvent,true, false, false, false, nil)
    if err != nil {
        log.Fatal(err)
    }

    h.channel.PublishWithContext(ctx, "", q.Name, false, false, amqp.Publishing{
        ContentType: "application/json",
        Body: marshalledOrder,
        DeliveryMode: amqp.Persistent,
    })

    return o, nil
}
