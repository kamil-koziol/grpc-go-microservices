package main

import (
	"context"
	"log"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/kamil-koziol/common"
	"github.com/kamil-koziol/common/discovery"
	"github.com/kamil-koziol/common/discovery/consul"
	"github.com/kamil-koziol/gateway/gateway"
)

var (
	serviceName = "gateway"
	httpAddr    = common.EnvString("HTTP_ADDR", ":8080")
	consulAddr  = common.EnvString("CONSUL_ADDR", "localhost:8500")
)

func main() {
	registry, err := consul.NewRegistry(consulAddr, serviceName)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)

	if err := registry.Register(ctx, instanceID, serviceName, httpAddr); err != nil {
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

	mux := http.NewServeMux()

	ordersGateway := gateway.NewGRPCGateway(registry)

	handler := NewHandler(ordersGateway)
	handler.registerRoutes(mux)

	log.Printf("Running on %s", httpAddr)

	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		log.Fatal("failed to start http server")
	}
}
