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

type LoginRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthServiceClient struct {
	Client pb.AuthServiceClient
}

func NewAuthServiceClient(cfg *config.Config) AuthServiceClient {
	conn, err := grpc.Dial(cfg.App.ServiceAuthPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect port %s : %v", cfg.App.ServiceAuthPort, err)
	}

	return AuthServiceClient{
		Client: pb.NewAuthServiceClient(conn),
	}
}

func (a *AuthServiceClient) Login(ctx *gin.Context) {
	b := LoginRequestBody{}

	if err := ctx.BindJSON(&b); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	res, err := a.Client.Login(context.Background(), &pb.LoginRequest{
		Email:    b.Email,
		Password: b.Password,
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusCreated, &res)
}

func (a *AuthServiceClient) Register(ctx *gin.Context) {
	body := RegisterRequestBody{}

	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	res, err := a.Client.Register(context.Background(), &pb.RegisterRequest{
		Email:    body.Email,
		Password: body.Password,
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
	}

	ctx.JSON(int(res.Status), &res)
}
