package service

import (
	"context"
	"github.com/zakariawahyu/go-grpc-microservices-gin-gateway/entity"
	"github.com/zakariawahyu/go-grpc-microservices-gin-gateway/pb"
	"github.com/zakariawahyu/go-grpc-microservices-gin-gateway/pkg/helpers"
	"gorm.io/gorm"
	"net/http"
)

type AuthService struct {
	DB  *gorm.DB
	Jwt helpers.JwtWrapper
	pb.UnimplementedAuthServiceServer
}

func (a *AuthService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	user := entity.User{}

	if err := a.DB.Where(&entity.User{Email: req.GetEmail()}).First(&user).Error; err != nil {
		return &pb.RegisterResponse{
			Status: http.StatusConflict,
			Error:  "E-Mail already exists",
		}, nil
	}

	user.Email = req.GetEmail()
	user.Password = helpers.HashPassword(req.GetPassword())

	a.DB.Create(&user)

	return &pb.RegisterResponse{
		Status: http.StatusCreated,
	}, nil
}

func (a *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user := entity.User{}

	if err := a.DB.Where(&entity.User{Email: req.GetEmail()}).First(&user).Error; err != nil {
		return &pb.LoginResponse{
			Status: http.StatusNotFound,
			Error:  "User not found",
		}, nil
	}

	match := helpers.CheckPasswordHash(req.GetPassword(), user.Password)

	if !match {
		return &pb.LoginResponse{
			Status: http.StatusBadRequest,
			Error:  "Wrong User Password",
		}, nil
	}

	token, _ := a.Jwt.GenerateToken(user)

	return &pb.LoginResponse{
		Status: http.StatusOK,
		Token:  token,
	}, nil
}

func (a *AuthService) Validate(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	user := entity.User{}

	claims, err := a.Jwt.ValidateToken(req.GetToken())
	if err != nil {
		return &pb.ValidateResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}, nil
	}

	if err := a.DB.Where(&entity.User{Email: claims.Email}).First(&user).Error; err != nil {
		return &pb.ValidateResponse{
			Status: http.StatusNotFound,
			Error:  "User not found",
		}, nil
	}

	return &pb.ValidateResponse{
		Status: http.StatusOK,
		UserId: user.Id,
	}, nil
}
