package controller

import (
	"general-shop/database"
	"general-shop/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AddProduct inserts a new product into the MongoDB collection
func AddProduct(c *gin.Context) {
	var product models.Product

	// Bind incoming JSON to product struct
	if err := c.BindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Generate new MongoDB ObjectID
	product.ID = primitive.NewObjectID()

	// Insert the product into the database
	_, err := database.ProductCollection.InsertOne(c, product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add product", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product added successfully"})
}

// GetProducts retrieves all products from the MongoDB collection
func GetProducts(c *gin.Context) {
	var products []models.Product

	// Find all documents in the collection
	cursor, err := database.ProductCollection.Find(c, bson.D{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}
	defer cursor.Close(c)

	// Iterate over the cursor and decode each document into a Product
	for cursor.Next(c) {
		var product models.Product
		if err := cursor.Decode(&product); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode product"})
			return
		}
		products = append(products, product)
	}

	// Return the list of products as JSON
	c.JSON(http.StatusOK, products)
}
