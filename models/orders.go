package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CartItem struct {
	ProductID string  `json:"product_id"`
	Name      string  `json:"name"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type Order struct {
	ID        string     `json:"id" bson:"_id"`
	Token     string     `json:"token" bson:"token"`
	UserID    string     `json:"user_id" bson:"user_id"`
	Items     []CartItem `json:"items" bson:"items"`
	Total     float64    `json:"total" bson:"total"`
	Status    string     `json:"status" bson:"status"`
	CreatedAt time.Time  `json:"created_at" bson:"created_at"`
}

type Cart struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	UserID string             `bson:"user_id"`
	Items  []CartItem         `json:"items" bson:"items"`
}
