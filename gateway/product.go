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
	"strconv"
)

type CreateProductRequestBody struct {
	Name  string `json:"name"`
	Stock int64  `json:"stock"`
	Price int64  `json:"price"`
}

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

func (p *ProductServiceClient) CreateProduct(ctx *gin.Context) {
	body := CreateProductRequestBody{}

	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	res, err := p.Client.CreateProduct(context.Background(), &pb.CreateProductRequest{
		Name:  body.Name,
		Stock: body.Stock,
		Price: body.Price,
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusCreated, &res)
}

func (p *ProductServiceClient) FindOne(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 32)

	res, err := p.Client.FindOne(context.Background(), &pb.FindOneRequest{
		Id: int64(id),
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusCreated, &res)
}
