package routes

import (
	"general-shop/controller"
	"general-shop/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// Public routes
	router.POST("/register", controller.Register)
	router.POST("/login", controller.Login)

	// Product routes (Admin only)
	admin := router.Group("/admin")
	admin.Use(middleware.AuthMiddleware())
	admin.POST("/product", controller.AddProduct)
	admin.GET("/products", controller.GetProducts)
}
