package main

import (
	// "fmt"
	"log"
	"os"

	"example.com/controller"
	"example.com/initializer"
	"example.com/model"
	"example.com/scheduler"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)


func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dsn := os.Getenv("DSN")
	initializer.Connect(dsn)
}

func main() {
	go scheduler.Cron_scheduler()
	r := gin.Default()
	initializer.DB.AutoMigrate(&model.User{})
	initializer.DB.AutoMigrate(&model.Order{})
	initializer.DB.AutoMigrate(&model.Merchant{})

	initializer.DB.Exec("set session transaction isolation level read committed")

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.POST("/createUser", controller.CreateUser)
	r.POST("/createOrder", controller.CreateOrder)
	r.PUT("/cancelOrder", controller.CancelOrder)
	r.PUT("/refundOrder", controller.RefundOrder)
	r.GET("/orderDetails", controller.OrderDetails)

	r.Run() // listen and serve on 0.0.0.0:8080
	
}
