package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/kamil-koziol/common"
	"github.com/kamil-koziol/common/broker"
	"github.com/kamil-koziol/common/discovery"
	"github.com/kamil-koziol/common/discovery/consul"
	"google.golang.org/grpc"
)

var (
	serviceName  = "orders"
	grpcAddr     = common.EnvString("GRPC_ADDR", "localhost:2000")
	consulAddr   = common.EnvString("CONSUL_ADDR", "localhost:8500")
	amqpUser     = common.EnvString("AMQP_USER", "guest")
	amqpPassword = common.EnvString("AMQP_PASSWORD", "guest")
	amqpHost     = common.EnvString("AMQP_HOST", "localhost")
	amqpPort     = common.EnvString("AMQP_PORT", "5672")
)

func main() {
	registry, err := consul.NewRegistry(consulAddr, serviceName)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)

	if err := registry.Register(ctx, instanceID, serviceName, grpcAddr); err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registry.HealthCheck(instanceID, serviceName); err != nil {
				log.Fatal("failed to healthcheck")
				time.Sleep(1 * time.Second)
			}
		}
	}()
	defer registry.Deregister(ctx, instanceID, serviceName)

	ch, close := broker.Connect(amqpUser, amqpPassword, amqpHost, amqpPort)
	defer func() {
		close()
		ch.Close()
	}()

	grpcServer := grpc.NewServer()
	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	store := NewStore()
	svc := NewService(store)
	NewGRPCHandler(grpcServer, svc, ch)

	svc.CreateOrder(ctx)

	log.Printf("GRPC Server Started at %s", grpcAddr)
	if err := grpcServer.Serve(l); err != nil {
		log.Fatal(err.Error())
	}
}
