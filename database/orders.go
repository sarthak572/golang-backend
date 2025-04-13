package database

import (
	"context"
	"general-shop/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func SaveOrder(order models.Order) error {
	collection := client.Database("shop").Collection("orders")
	_, err := collection.InsertOne(context.TODO(), order)
	return err
}

func FindCartByUserID(userID string) (models.Cart, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var cart models.Cart
	err := CartCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&cart)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Cart{}, nil // Return an empty cart if not found
		}
		return models.Cart{}, err
	}
	return cart, nil
}

// DeleteCartByUserID deletes the cart for a user by their userID
func DeleteCartByUserID(userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := CartCollection.DeleteOne(ctx, bson.M{"user_id": userID})
	return err
}

func FindOrderByToken(token string) (*models.Order, error) {
	collection := client.Database("shop").Collection("orders")
	var order models.Order
	err := collection.FindOne(context.TODO(), bson.M{"token": token}).Decode(&order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func DeleteOrderByToken(token string) error {
	collection := client.Database("shop").Collection("orders")
	_, err := collection.DeleteOne(context.TODO(), bson.M{"token": token})
	return err
}
