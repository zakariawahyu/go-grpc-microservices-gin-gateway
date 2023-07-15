package main

import (
	"github.com/zakariawahyu/go-grpc-microservices-gin-gateway/config"
	"github.com/zakariawahyu/go-grpc-microservices-gin-gateway/pb"
	"github.com/zakariawahyu/go-grpc-microservices-gin-gateway/pkg/db"
	"github.com/zakariawahyu/go-grpc-microservices-gin-gateway/services/order/client"
	"github.com/zakariawahyu/go-grpc-microservices-gin-gateway/services/order/server/service"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	cfg := config.NewConfig()
	conn := db.InitDatabase(cfg)

	l, err := net.Listen("tcp", cfg.App.ServiceOrderPort)
	if err != nil {
		log.Fatalf("could not listen to %s: %v", cfg.App.ServiceOrderPort, err)
	}

	productServices := client.NewProductServiceClient(cfg)

	grpcServer := grpc.NewServer()
	orderServer := service.OrderService{
		DB:             conn,
		ProductService: productServices,
	}

	pb.RegisterOrderServiceServer(grpcServer, &orderServer)
	log.Fatal(grpcServer.Serve(l))
}
