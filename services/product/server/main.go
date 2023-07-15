package main

import (
	"github.com/zakariawahyu/go-grpc-microservices-gin-gateway/config"
	"github.com/zakariawahyu/go-grpc-microservices-gin-gateway/pb"
	"github.com/zakariawahyu/go-grpc-microservices-gin-gateway/pkg/db"
	"github.com/zakariawahyu/go-grpc-microservices-gin-gateway/services/product/server/service"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	cfg := config.NewConfig()
	conn := db.InitDatabase(cfg)

	l, err := net.Listen("tcp", cfg.App.ServiceProductPort)
	if err != nil {
		log.Fatalf("could not listen to %s: %v", cfg.App.ServiceProductPort, err)
	}

	grpcServer := grpc.NewServer()
	productServer := service.ProductService{
		DB: conn,
	}

	pb.RegisterProductServiceServer(grpcServer, &productServer)
	log.Fatal(grpcServer.Serve(l))
}
