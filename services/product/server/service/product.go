package service

import (
	"context"
	"github.com/zakariawahyu/go-grpc-microservices-gin-gateway/entity"
	"github.com/zakariawahyu/go-grpc-microservices-gin-gateway/pb"
	"gorm.io/gorm"
	"net/http"
)

type ProductService struct {
	DB *gorm.DB
	pb.UnimplementedProductServiceServer
}

func (p *ProductService) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	product := entity.Product{}

	product.Name = req.GetName()
	product.Stock = req.GetStock()
	product.Price = req.GetPrice()

	if err := p.DB.Create(&product); err.Error != nil {
		return &pb.CreateProductResponse{
			Status: http.StatusConflict,
			Error:  err.Error.Error(),
		}, nil
	}

	return &pb.CreateProductResponse{
		Status: http.StatusCreated,
		Id:     product.Id,
	}, nil
}

func (p *ProductService) FindOne(ctx context.Context, req *pb.FindOneRequest) (*pb.FindOneResponse, error) {
	product := entity.Product{}

	if err := p.DB.First(&product, req.GetId()); err.Error != nil {
		return &pb.FindOneResponse{
			Status: http.StatusNotFound,
			Error:  err.Error.Error(),
		}, nil
	}

	data := &pb.FindOneData{
		Id:    product.Id,
		Name:  product.Name,
		Stock: product.Stock,
		Price: product.Price,
	}

	return &pb.FindOneResponse{
		Status: http.StatusOK,
		Data:   data,
	}, nil
}

func (p *ProductService) DecreaseStock(ctx context.Context, req *pb.DecreaseStockRequest) (*pb.DecreaseStockResponse, error) {
	product := entity.Product{}

	if err := p.DB.First(&product, req.GetId()).Error; err != nil {
		return &pb.DecreaseStockResponse{
			Status: http.StatusNotFound,
			Error:  err.Error(),
		}, nil
	}

	if product.Stock <= 0 {
		return &pb.DecreaseStockResponse{
			Status: http.StatusConflict,
			Error:  "Stock too low",
		}, nil
	}

	decreaseLog := entity.StockDecreaseLog{}
	if err := p.DB.Where(&entity.StockDecreaseLog{OrderId: req.GetOrderId()}).First(&decreaseLog).Error; err != nil {
		return &pb.DecreaseStockResponse{
			Status: http.StatusConflict,
			Error:  "Stock already decreased",
		}, nil
	}

	product.Stock = product.Stock - req.Quantity
	p.DB.Save(&product)

	decreaseLog.OrderId = req.GetOrderId()
	decreaseLog.ProductRefer = product.Id
	p.DB.Create(&decreaseLog)

	return &pb.DecreaseStockResponse{
		Status: http.StatusOK,
	}, nil
}
