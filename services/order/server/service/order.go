package service

import (
	"context"
	"github.com/zakariawahyu/go-grpc-microservices-gin-gateway/entity"
	"github.com/zakariawahyu/go-grpc-microservices-gin-gateway/pb"
	"github.com/zakariawahyu/go-grpc-microservices-gin-gateway/services/order/client"
	"gorm.io/gorm"
	"net/http"
)

type OrderService struct {
	DB             *gorm.DB
	ProductService client.ProductServiceClient
	pb.UnimplementedOrderServiceServer
}

func (o *OrderService) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	product, err := o.ProductService.FindOne(req.GetProductId())

	if err != nil {
		return &pb.CreateOrderResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}, nil
	} else if product.GetStatus() >= http.StatusNotFound {
		return &pb.CreateOrderResponse{
			Status: product.GetStatus(),
			Error:  product.GetError(),
		}, nil
	} else if product.Data.GetStock() < req.GetQuantity() {
		return &pb.CreateOrderResponse{
			Status: http.StatusConflict,
			Error:  "Stock too less",
		}, nil
	}

	order := entity.Order{
		Price:     product.Data.GetPrice() * req.GetQuantity(),
		ProductId: product.Data.Id,
		UserId:    req.UserId,
		Quantity:  req.Quantity,
	}

	o.DB.Create(&order)

	decreaseStock, err := o.ProductService.DecreaseStock(req.GetProductId(), order.Id, order.Quantity)
	if err != nil {
		return &pb.CreateOrderResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}, nil
	} else if decreaseStock.GetStatus() == http.StatusConflict {
		o.DB.Delete(&entity.Order{}, order.Id)

		return &pb.CreateOrderResponse{
			Status: http.StatusConflict,
			Error:  decreaseStock.GetError(),
		}, nil
	}

	return &pb.CreateOrderResponse{
		Status: http.StatusCreated,
		Id:     order.Id,
	}, nil
}
