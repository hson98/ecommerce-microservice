package postgres

import (
	"fmt"
	"github.com/hson98/ecommerce-microservice/src/auth-service/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

const (
	maxOpenConns    = 60
	connMaxLifetime = 120
	maxIdleConns    = 30
	connMaxIdleTime = 20
)

func NewPostgresDB(config *config.Config) *gorm.DB {
	username := config.Server.DBUser
	password := config.Server.DBPass
	dbName := config.Server.DBName
	dbHost := config.Server.DBHost
	port := config.Server.DBPort
	//dsn = "host=localhost auth=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"require
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable  TimeZone=Asia/Ho_Chi_Minh", dbHost, username, password, dbName, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetConnMaxLifetime(connMaxLifetime * time.Second)
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetConnMaxIdleTime(connMaxIdleTime * time.Second)

	return db
}
