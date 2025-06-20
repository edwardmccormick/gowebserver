package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"

	"github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	// Connect to MongoDB
	// mongoClient, err := ConnectToMongoDB()
	// if err != nil {
	//     // Handle error
	// }
	// defer mongoClient.Disconnect(context.TODO())

	// // Connect to MySQL
	// mysqlDB, err := ConnectToMySQL()
	// if err != nil {
	//     // Handle error
	// }
	// defer mysqlDB.Close()

	// // Read JSON file
	// file, err := os.Open("data.json")
	// if err != nil {
	//     // Handle error
	// }
	// defer file.Close()

	// // Decode JSON
	// var data map[string]interface{}
	// err = json.NewDecoder(file).Decode(&data)
	// if err != nil {
	//     // Handle error
	// }

	// // Insert data into MongoDB
	// collection := mongoClient.Database("mydb").Collection("mycollection")
	// _, err = collection.InsertOne(context.TODO(), data)
	// if err != nil {
	//     // Handle error
	// }

	// // Insert data into MySQL
	// _, err = mysqlDB.Exec("INSERT INTO mytable (data) VALUES (?)", data)
	// if err != nil {
	//     // Handle error
	// }
}

func ConnectToMongoDB() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017") // Or your Atlas connection string
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		// Handle error
	}

}

func ConnectToMySQL() (*sql.DB, error) {
	dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	return db, nil
}

func ParseConfigJson() {
	jsonFile, err := os.Open("users.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened users.json")

	byteValue, _ := io.ReadAll(jsonFile)

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
}
