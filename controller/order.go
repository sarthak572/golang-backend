package controller

import (
	"context"
	"general-shop/database"
	"general-shop/models"
	"general-shop/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetCart(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	var cart models.Cart
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := database.CartCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&cart)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(200, gin.H{"items": []models.CartItem{}})
			return
		}
		c.JSON(500, gin.H{"error": "Failed to get cart"})
		return
	}

	c.JSON(200, cart)
}

func AddItemToCart(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	var item models.CartItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Upsert cart: if not exists, create new
	filter := bson.M{"user_id": userID}
	update := bson.M{
		"$push": bson.M{"items": item},
	}
	opts := options.Update().SetUpsert(true)

	_, err := database.CartCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to add item to cart"})
		return
	}

	c.JSON(200, gin.H{"message": "Item added to cart"})
}

func RemoveItemFromCart(c *gin.Context) {
	var request struct {
		ProductID string `json:"product_id"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find the user's cart
	cart, err := database.FindCartByUserID(userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get cart"})
		return
	}

	// Remove the item from the cart
	updatedItems := []models.CartItem{}
	for _, item := range cart.Items {
		if item.ProductID != request.ProductID {
			updatedItems = append(updatedItems, item)
		}
	}

	// Update the cart in the database
	_, err = database.CartCollection.UpdateOne(
		ctx,
		bson.M{"user_id": userID},
		bson.M{"$set": bson.M{"items": updatedItems}},
	)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to remove item from cart"})
		return
	}

	c.JSON(200, gin.H{"message": "Item removed from cart", "items": updatedItems})
}

func Checkout(c *gin.Context) {
	userID := c.GetString("userID") // from AuthMiddleware

	// Fetch cart from DB
	cart, err := database.FindCartByUserID(userID)
	if err != nil || len(cart.Items) == 0 {
		c.JSON(400, gin.H{"error": "Cart is empty or not found"})
		return
	}

	order := models.Order{
		ID:        uuid.New().String(),
		Token:     uuid.New().String(),
		UserID:    userID,
		Items:     cart.Items,
		Total:     utils.CalculateTotal(cart.Items),
		Status:    "pending",
		CreatedAt: time.Now(),
	}

	if err := database.SaveOrder(order); err != nil {
		c.JSON(500, gin.H{"error": "Could not save order"})
		return
	}

	qrCode, err := utils.GenerateQRCode(order.Token)
	if err != nil {
		c.JSON(500, gin.H{"error": "QR generation failed"})
		return
	}

	// Clear cart after checkout
	if err := database.DeleteCartByUserID(userID); err != nil {
		c.JSON(500, gin.H{"error": "Failed to clear cart"})
		return
	}

	c.Data(200, "image/png", qrCode)
}

func VerifyOrderByQR(c *gin.Context) {
	var request struct {
		Token string `json:"token"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	order, err := database.FindOrderByToken(request.Token)
	if err != nil {
		c.JSON(404, gin.H{"error": "Order not found"})
		return
	}

	err = database.DeleteOrderByToken(request.Token)
	if err != nil {
		c.JSON(500, gin.H{"error": "Could not complete order"})
		return
	}
	order.Status = "completed"
	c.JSON(200, gin.H{"message": "Order verified and completed", "order": order})
}
