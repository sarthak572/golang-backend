package controller

import (
	"general-shop/database"
	"general-shop/models"
	"general-shop/utils"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetCart(c *gin.Context) {
	session := sessions.Default(c)
	cart := session.Get("cart")
	if cart == nil {
		cart = models.Cart{Items: []models.CartItem{}}
	}
	c.JSON(200, cart)
}

func AddItemToCart(c *gin.Context) {
	var item models.CartItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	session := sessions.Default(c)
	cart := session.Get("cart")
	if cart == nil {
		cart = models.Cart{Items: []models.CartItem{}}
	}

	currentCart := cart.(models.Cart)
	found := false

	for i, cartItem := range currentCart.Items {
		if cartItem.ProductID == item.ProductID {
			currentCart.Items[i].Quantity += item.Quantity
			found = true
			break
		}
	}

	if !found {
		currentCart.Items = append(currentCart.Items, item)
	}

	session.Set("cart", currentCart)
	session.Save()
	c.JSON(200, currentCart)
}
func RemoveItemFromCart(c *gin.Context) {
	var request struct {
		ProductID string `json:"product_id"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	session := sessions.Default(c)
	cart := session.Get("cart")
	if cart == nil {
		c.JSON(404, gin.H{"error": "Cart not found"})
		return
	}

	currentCart := cart.(models.Cart)
	for i, item := range currentCart.Items {
		if item.ProductID == request.ProductID {
			currentCart.Items = append(currentCart.Items[:i], currentCart.Items[i+1:]...)
			break
		}
	}

	session.Set("cart", currentCart)
	session.Save()
	c.JSON(200, currentCart)
}

func Checkout(c *gin.Context) {
	// This assumes session already has the cart
	session := sessions.Default(c)
	cart := session.Get("cart")
	if cart == nil {
		c.JSON(400, gin.H{"error": "Cart is empty"})
		return
	}
	currentCart := cart.(models.Cart) // Make sure your Cart struct is in models

	order := models.Order{
		ID:        uuid.New().String(),
		Token:     uuid.New().String(),
		UserID:    "demo_user", // Replace with actual user ID if logged in
		Items:     currentCart.Items,
		Total:     utils.CalculateTotal(currentCart.Items),
		Status:    "pending",
		CreatedAt: time.Now(),
	}

	err := database.SaveOrder(order)
	if err != nil {
		c.JSON(500, gin.H{"error": "Could not save order"})
		return
	}

	qrCode, err := utils.GenerateQRCode(order.Token)
	if err != nil {
		c.JSON(500, gin.H{"error": "QR generation failed"})
		return
	}

	// Delete cart
	session.Delete("cart")
	session.Save()

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

	c.JSON(200, gin.H{"message": "Order verified and completed", "order": order})
}
