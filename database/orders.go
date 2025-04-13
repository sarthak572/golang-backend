package database

import (
	"context"
	"general-shop/models"

	"go.mongodb.org/mongo-driver/bson"
)

func SaveOrder(order models.Order) error {
	collection := client.Database("shop").Collection("orders")
	_, err := collection.InsertOne(context.TODO(), order)
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
