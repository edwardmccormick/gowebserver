package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectToMongoDB() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017") // Or your Atlas connection string
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Errorf("Failed to connect to MongoDB: " + err.Error())
	}
	return client, nil
}

func ConnectToMySQL() (*gorm.DB, error) {
	dsn := "user:password@tcp(localhost:3306)/database?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Errorf("Failed to connect to MySQL: " + err.Error())
	}

	return db, nil
}
