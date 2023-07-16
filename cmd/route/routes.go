package route

import (
	"github.com/gin-gonic/gin"
	"github.com/zakariawahyu/go-grpc-microservices-gin-gateway/config"
	"github.com/zakariawahyu/go-grpc-microservices-gin-gateway/gateway"
)

func RegisterAuthRoutes(r *gin.Engine, cfg *config.Config) *gateway.AuthServiceClient {
	authService := gateway.NewAuthServiceClient(cfg)

	routes := r.Group("/auth")
	routes.POST("/register", authService.Register)
	routes.POST("/login", authService.Login)

	return &authService
}

func RegisterOrderRoutes(r *gin.Engine, cfg *config.Config, authService *gateway.AuthServiceClient) {
	auth := InitAuthMiddleware(authService)

	orderServices := gateway.NewOrderServiceClient(cfg)

	routes := r.Group("/order")
	routes.Use(auth.AuthRequired)
	routes.POST("/", orderServices.CreateOrder)
}

func RegisterProductRoutes(r *gin.Engine, cfg *config.Config) {
	productService := gateway.NewProductServiceClient(cfg)

	routes := r.Group("/product")
	routes.POST("/", productService.CreateProduct)
	routes.GET("/:id", productService.FindOne)
}
