package gateway

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/zakariawahyu/go-grpc-microservices-gin-gateway/config"
	"github.com/zakariawahyu/go-grpc-microservices-gin-gateway/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
)

type CreateOrderRequestBody struct {
	ProductId int64 `json:"productId"`
	Quantity  int64 `json:"quantity"`
}

type OrderServiceClient struct {
	Client pb.OrderServiceClient
}

func NewOrderServiceClient(cfg *config.Config) OrderServiceClient {
	conn, err := grpc.Dial(cfg.App.ServiceOrderPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect port %s : %v", cfg.App.ServiceProductPort, err)
	}

	return OrderServiceClient{
		Client: pb.NewOrderServiceClient(conn),
	}
}

func (o *OrderServiceClient) CreateOrder(ctx *gin.Context) {
	body := CreateOrderRequestBody{}
	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	userId, _ := ctx.Get("userId")

	res, err := o.Client.CreateOrder(context.Background(), &pb.CreateOrderRequest{
		ProductId: body.ProductId,
		Quantity:  body.Quantity,
		UserId:    userId.(int64),
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusCreated, &res)
}
