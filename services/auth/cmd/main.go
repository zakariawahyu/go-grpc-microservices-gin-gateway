package main

import (
	"github.com/zakariawahyu/go-grpc-microservices-gin-gateway/config"
	"github.com/zakariawahyu/go-grpc-microservices-gin-gateway/pb"
	"github.com/zakariawahyu/go-grpc-microservices-gin-gateway/pkg/db"
	"github.com/zakariawahyu/go-grpc-microservices-gin-gateway/services/auth/server/service"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	cfg := config.NewConfig()
	conn := db.InitDatabase(cfg)

	l, err := net.Listen("tcp", cfg.App.ServiceAuthPort)
	if err != nil {
		log.Fatalf("could not listen to %s: %v", cfg.App.ServiceAuthPort, err)
	}

	grpcServer := grpc.NewServer()
	authServer := service.AuthService{
		DB: conn,
	}

	pb.RegisterAuthServiceServer(grpcServer, &authServer)
	log.Fatal(grpcServer.Serve(l))
}
