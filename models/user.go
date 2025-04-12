package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username string             `bson:"username" json:"username"`
	Password string             `bson:"password" json:"password"` // Password will be hashed
	Email    string             `bson:"email" json:"email"`
	Role     string             `bson:"role" json:"role"` // "admin" or "customer"
}
