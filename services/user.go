package services

import (
	"context"
	"errors"
	"general-shop/database"
	"general-shop/models"
	"general-shop/utils"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateUser inserts a new user into MongoDB
func CreateUser(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if the collection is initialized
	if database.UserCollection == nil {
		log.Fatal("UserCollection is nil")
		return errors.New("database connection error")
	}

	// Check if user already exists (by email)
	var existing models.User
	err := database.UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&existing)
	if err == nil {
		log.Println("User already exists:", existing)
		return errors.New("user already exists")
	}
	if err != mongo.ErrNoDocuments {
		log.Println("Error checking user existence:", err)
		return err
	}

	// Insert the new user into the database
	_, err = database.UserCollection.InsertOne(ctx, user)
	if err != nil {
		log.Println("Error inserting user into DB:", err)
		return err
	}

	return nil
}

// LoginUser checks credentials and returns a JWT token
func LoginUser(user models.User) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if the collection is initialized
	if database.UserCollection == nil {
		log.Fatal("UserCollection is nil")
		return "", errors.New("database connection error")
	}

	var dbUser models.User
	err := database.UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&dbUser)
	if err != nil {
		log.Println("User not found:", err)
		return "", errors.New("user not found")
	}

	// Check if the password matches the stored hash
	if !utils.CheckPasswordHash(user.Password, dbUser.Password) {
		log.Println("Invalid password attempt for email:", user.Email)
		return "", errors.New("invalid password")
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(dbUser.Email)
	if err != nil {
		log.Println("Failed to generate JWT token:", err)
		return "", errors.New("failed to generate token")
	}

	return token, nil
}
