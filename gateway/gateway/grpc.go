package gateway

import (
	"context"
	"log"

	pb "github.com/kamil-koziol/common/api"
	"github.com/kamil-koziol/common/discovery"
)

type gateway struct {
	registry discovery.Registry
}

func NewGRPCGateway(registry discovery.Registry) *gateway {
	return &gateway{registry}
}

func (g *gateway) CreateOrder(ctx context.Context, p *pb.CreateOrderRequest) (*pb.Order, error) {
	conn, err := discovery.ServiceConnection(ctx, "orders", g.registry)
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	c := pb.NewOrderServiceClient(conn)
	return c.CreateOrder(ctx, p)
}
