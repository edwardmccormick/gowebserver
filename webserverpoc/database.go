package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func LoadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(bytes, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// ConnectToMongoDBWithConfig connects to MongoDB using values from the config
func ConnectToMongoDBWithConfig(config *Config) (*mongo.Client, error) {
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d/?authSource=admin",
		//  uri := fmt.Sprintf("mongodb://%s:%s@%s:%d/",
		config.Mongo.User,
		config.Mongo.Password,
		config.Mongo.Host,
		config.Mongo.Port,
		// config.Mongo.Database,
	)

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}
	// Test the connection and authentication
	if err := client.Ping(context.TODO(), nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	return client, nil
}

// ConnectToMySQLWithConfig connects to MySQL using values from the config
func ConnectToMySQLWithConfig(config *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.MySQL.User,
		config.MySQL.Password,
		config.MySQL.Host,
		config.MySQL.Port,
		config.MySQL.Database,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL: %w", err)
	}

	return db, nil
}

func PopulateDatabase(db *gorm.DB, mongoClient *mongo.Client) error {
	// Check if users exist in MySQL
	var userCount int64
	if err := db.Model(&User{}).Count(&userCount).Error; err != nil {
		return fmt.Errorf("failed to check user count: %w", err)
	}

	if userCount == 0 {
		// Populate users
		if err := db.Create(&users).Error; err != nil {
			return fmt.Errorf("failed to populate users: %w", err)
		}
		fmt.Println("Users populated in MySQL.")
	}

	// Check if people exist in MySQL
	var peopleCount int64
	if err := db.Model(&Person{}).Count(&peopleCount).Error; err != nil {
		return fmt.Errorf("failed to check people count: %w", err)
	}

	if peopleCount == 0 {
		// Populate people
		if err := db.Create(&people).Error; err != nil {
			return fmt.Errorf("failed to populate people: %w", err)
		}
		fmt.Println("People populated in MySQL.")
	}

	// // Check if matches exist in MySQL
	var matchCount int64
	if err := db.Model(&Match{}).Count(&matchCount).Error; err != nil {
		return fmt.Errorf("failed to check match count: %w", err)
	}

	if matchCount == 0 {
		// Populate matches
		if err := db.Create(&Matches).Error; err != nil {
			return fmt.Errorf("failed to populate matches: %w", err)
		}
		fmt.Println("Matches populated in MySQL.")
	}

	// Check if PhotoArray1 exists in MongoDB
	photoCollection := mongoClient.Database("urmid").Collection("photos")
	photoCount, err := photoCollection.CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		return fmt.Errorf("failed to check photo count: %w", err)
	}

	if photoCount == 0 {
		// Populate PhotoArray1 and PhotoArray2
		// photos := append(PhotoArray, PhotoArray2...)
		var photoDocs []interface{}
		for _, photo := range albums {
			photoDocs = append(photoDocs, photo)
		}

		if _, err := photoCollection.InsertMany(context.TODO(), photoDocs); err != nil {
			return fmt.Errorf("failed to populate photos: %w", err)
		}
		fmt.Println("Photos populated in MongoDB.")
	}

	return nil
}
