package db

import (
	"github.com/zakariawahyu/go-grpc-microservices-gin-gateway/config"
	"github.com/zakariawahyu/go-grpc-microservices-gin-gateway/entity"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

func InitDatabase(cfg *config.Config) *gorm.DB {
	db, err := gorm.Open(mysql.Open(cfg.Db.DSN), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&entity.Product{})
	db.AutoMigrate(&entity.StockDecreaseLog{})

	return db
}
