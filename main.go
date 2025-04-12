package main

import (
	"general-shop/database"
	"general-shop/routes"

	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()
	log.Println("databse connected")

	r := gin.Default()

	routes.SetupRoutes(r)

	if err := r.Run(":8080"); err != nil {
		log.Fatal("server failed to start ", err)
	}
	log.Println("server started on 8080")
}
