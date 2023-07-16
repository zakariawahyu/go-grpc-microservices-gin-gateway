package main

import (
	"github.com/gin-gonic/gin"
	"github.com/zakariawahyu/go-grpc-microservices-gin-gateway/cmd/route"
	"github.com/zakariawahyu/go-grpc-microservices-gin-gateway/config"
)

func main() {
	cfg := config.NewConfig()
	r := gin.Default()

	authService := route.RegisterAuthRoutes(r, cfg)
	route.RegisterOrderRoutes(r, cfg, authService)
	route.RegisterProductRoutes(r, cfg)

	r.Run(cfg.App.AppGatewayPort)
}
