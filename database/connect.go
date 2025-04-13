package database

import (
	"context"
	"fmt"
	"general-shop/config"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client            *mongo.Client
	UserCollection    *mongo.Collection
	ProductCollection *mongo.Collection
	CartCollection    *mongo.Collection
)

func Connect() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	clientOptions := options.Client().ApplyURI(config.MONGO_URI)
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("MongoDB Ping Failed:", err)
	}

	fmt.Println("Connected to MongoDB!")

	// Initialize collections after client is connected
	UserCollection = client.Database("general-shop").Collection("users")
	ProductCollection = client.Database("general-shop").Collection("products")
	CartCollection = client.Database("general-shop").Collection("carts")

}

func GetClient() *mongo.Client {
	return client
}
