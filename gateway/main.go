package main

import (
	"log"
	"net/http"

	_ "github.com/joho/godotenv/autoload"
	"github.com/kamil-koziol/common"
	pb "github.com/kamil-koziol/common/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	httpAddr         = common.EnvString("HTTP_ADDR", ":3000")
	orderServiceAddr = "localhost:3000"
)

func main() {
	conn, err := grpc.NewClient(orderServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to make new client connection")
	}
	defer conn.Close()

	log.Println("Dialong orders service at", orderServiceAddr)
	c := pb.NewOrderServiceClient(conn)

	mux := http.NewServeMux()
	handler := NewHandler(c)
	handler.registerRoutes(mux)

	log.Printf("Running on %s", httpAddr)

	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		log.Fatal("failed to start http server")
	}
}
