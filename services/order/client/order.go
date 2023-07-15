package client

import (
	"context"
	"github.com/zakariawahyu/go-grpc-microservices-gin-gateway/config"
	"github.com/zakariawahyu/go-grpc-microservices-gin-gateway/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

type ProductServiceClient struct {
	Client pb.ProductServiceClient
}

func NewProductServiceClient(cfg *config.Config) ProductServiceClient {
	conn, err := grpc.Dial(cfg.App.ServiceProductPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect port %s : %v", cfg.App.ServiceProductPort, err)
	}

	return ProductServiceClient{
		Client: pb.NewProductServiceClient(conn),
	}
}

func (p *ProductServiceClient) FindOne(productID int64) (*pb.FindOneResponse, error) {
	req := &pb.FindOneRequest{
		Id: productID,
	}

	return p.Client.FindOne(context.Background(), req)
}

func (p *ProductServiceClient) DecreaseStock(productId int64, orderId int64, quantity int64) (*pb.DecreaseStockResponse, error) {
	req := &pb.DecreaseStockRequest{
		Id:       productId,
		OrderId:  orderId,
		Quantity: quantity,
	}

	return p.Client.DecreaseStock(context.Background(), req)
}
