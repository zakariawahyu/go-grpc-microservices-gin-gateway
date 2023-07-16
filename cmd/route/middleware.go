package route

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/zakariawahyu/go-grpc-microservices-gin-gateway/gateway"
	"github.com/zakariawahyu/go-grpc-microservices-gin-gateway/pb"
	"net/http"
	"strings"
)

type AuthMiddlewareConfig struct {
	authService *gateway.AuthServiceClient
}

func InitAuthMiddleware(authService *gateway.AuthServiceClient) AuthMiddlewareConfig {
	return AuthMiddlewareConfig{
		authService: authService,
	}
}

func (mw *AuthMiddlewareConfig) AuthRequired(ctx *gin.Context) {

	authorization := ctx.Request.Header.Get("authorization")

	if authorization == "" {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	token := strings.Split(authorization, "Bearer ")

	if len(token) < 2 {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	res, err := mw.authService.Client.Validate(context.Background(), &pb.ValidateRequest{Token: token[1]})
	if err != nil || res.Status != http.StatusOK {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ctx.Set("userId", res.UserId)

	ctx.Next()
}
