package routes

import (
	"general-shop/controller"
	"general-shop/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// Public auth routes
	router.POST("/register", controller.Register)
	router.POST("/login", controller.Login)

	// Admin-only product routes
	admin := router.Group("/admin")
	admin.Use(middleware.AuthMiddleware())
	{
		admin.POST("/product", controller.AddProduct)
		admin.GET("/products", controller.GetProducts)
	}

	// Cart routes (accessible to any logged-in user if you want to protect it)
	cart := router.Group("/cart")
	{
		cart.GET("/", controller.GetCart)
		cart.POST("/add", controller.AddItemToCart)
		cart.POST("/remove", controller.RemoveItemFromCart)
	}

	// Order routes (checkout and verify)
	order := router.Group("/order")
	{
		order.POST("/checkout", controller.Checkout)
		order.POST("/verify", controller.VerifyOrderByQR)
	}
}
