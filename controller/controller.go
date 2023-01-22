package controller

import (
	"fmt"

	"example.com/initializer"
	"example.com/model"
	"example.com/scheduler"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateUser(c *gin.Context) {
	var body struct {
		Name    string
		Balance int
	}
	c.Bind(&body)
	usr := model.User{Name: body.Name, Balance: uint(body.Balance)}
	result := initializer.DB.Create(&usr)
	if result.Error != nil {
		c.Status(400)
		return
	}
	c.JSON(200, gin.H{
		"user": usr,
	})
}

func CreateOrder(c *gin.Context) {
	var body struct {
		UserId uint
		Amount uint
	}
	c.Bind(&body)
	order := model.Order{UserId: body.UserId, Amount: body.Amount, Status: "CREATED"}

	initializer.DB.Transaction(func(tx *gorm.DB) error {
		var user = model.User{}
		err := tx.First(&user, body.UserId).Error

		if err != nil {
			c.JSON(400, ErrRes{Error: "User doesn't exist"})
			return err
		}

		if user.Balance >= body.Amount {

			if err := tx.Create(&order).Error; err != nil {
				return err
			}

			if err := tx.Model(&user).Update("Balance", user.Balance-body.Amount).Error; err != nil {
				fmt.Println("error updating here", err)
				return err
			}

			c.JSON(200, order)

			// push order to the queue for further preocessing
			scheduler.Orders_queue = append(scheduler.Orders_queue, order)

			// return nil will commit the whole transaction
			return nil

		} else {
			c.JSON(400, ErrRes{Error: "Not enough balance"})
			return nil
		}
	})
}

func CancelOrder(c *gin.Context) {
	var body struct {
		OrderId uint
	}
	c.Bind(&body)
	var order = model.Order{}

	if err := initializer.DB.First(&order, body.OrderId).Error; err != nil {
		c.JSON(400, ErrRes{Error: "Order doesn't exist"})
		return
	}

	if order.Status == "CREATED" {
		if err := initializer.DB.Model(&order).Update("Status", "CANCELLED").Error; err != nil {
			c.JSON(400, ErrRes{Error: "Error in cancelling"})
			return
		}
	} else {
		c.JSON(400, ErrRes{Error: "order not confirmed"})
		return
	}

	c.JSON(200, order)
}

func RefundOrder(c *gin.Context) {
	var body struct {
		OrderId uint
	}
	c.Bind(&body)
	var order = model.Order{}

	if err := initializer.DB.First(&order, body.OrderId).Error; err != nil {
		c.JSON(400, ErrRes{Error: "Order doesn't exist"})
		return
	}

	if order.Status != "CANCELLED" {
		c.JSON(400, ErrRes{Error: "order not cancelled"})
		return
	}

	initializer.DB.Transaction(func(tx *gorm.DB) error {

		var user = model.User{}
		err := tx.First(&user, order.UserId).Error

		if err != nil {
			c.JSON(400, ErrRes{Error: "User doesn't exist"})
			return err
		}

		if err := tx.Model(&order).Update("Status", "REFUNDED").Error; err != nil {
			return err
		}

		if err := tx.Model(&user).Update("Balance", user.Balance+order.Amount).Error; err != nil {
			return err
		}

		c.JSON(200, order)

		// return nil will commit the whole transaction
		return nil

	})

}

func OrderDetails(c *gin.Context) {
	var body struct {
		OrderId uint
	}
	c.BindQuery(&body)
	var order = model.Order{}

	if err := initializer.DB.First(&order, body.OrderId).Error; err != nil {
		c.JSON(400, ErrRes{Error: "Order doesn't exist"})
		return
	}

	c.JSON(200, order)
}

type ErrRes struct {
	Error string
}

type Res struct {
	Response string
}
