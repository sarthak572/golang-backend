package main

import (
	"encoding/gob"
	"general-shop/database"
	"general-shop/models"
	"general-shop/routes"

	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func init() {
	gob.Register(models.Cart{})
	gob.Register([]models.CartItem{})
}

func main() {
	database.Connect()
	log.Println("databse connected")

	r := gin.Default()
	store := cookie.NewStore([]byte("secret"))

	r.Use(sessions.Sessions("mysession", store))

	routes.SetupRoutes(r)

	if err := r.Run(":8080"); err != nil {
		log.Fatal("server failed to start ", err)
	}
	log.Println("server started on 8080")
}
