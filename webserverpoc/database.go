package gowebserver

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
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

// ConnectToPostgresWithConfig connects to Postgres using values from the config.
// Legacy mysql config is still accepted as a fallback during migration.
func ConnectToPostgresWithConfig(config *Config) (*gorm.DB, error) {
	sqlConfig := config.Postgres
	if sqlConfig.Host == "" {
		sqlConfig.Host = config.MySQL.Host
		sqlConfig.Port = config.MySQL.Port
		sqlConfig.User = config.MySQL.User
		sqlConfig.Password = config.MySQL.Password
		sqlConfig.Database = config.MySQL.Database
	}

	sslMode := sqlConfig.SSLMode
	if sslMode == "" {
		sslMode = "disable"
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC",
		sqlConfig.Host,
		sqlConfig.Port,
		sqlConfig.User,
		sqlConfig.Password,
		sqlConfig.Database,
		sslMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Postgres: %w", err)
	}

	return db, nil
}
