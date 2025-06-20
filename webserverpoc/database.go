package main

// import (
// 	"context"
// 	"database/sql"

// 	"github.com/go-sql-driver/mysql"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// 	"gorm.io/driver/mysql"
// 	"gorm.io/gorm"
// )

// func ConnectToMongoDB() (*mongo.Client, error) {
// 	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017") // Or your Atlas connection string
// 	client, err := mongo.Connect(context.TODO(), clientOptions)
// 	if err != nil {
// 		// Handle error
// 	}

// }

// func ConnectToMySQL() (*gorm.DB, error) {
// 	dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
// 	db, err := gorm.Open(gorm.Open(dsn), &gorm.Config{})

// 	return db, nil
// }
