package scheduler

import (
	"time"

	"example.com/initializer"
	"example.com/model"
	"github.com/go-co-op/gocron"
	"gorm.io/gorm"
)

var Orders_queue []model.Order

func Cron_scheduler() {
	s := gocron.NewScheduler(time.UTC)

	// Wait for prev cron job to finish before starting a new one
	s.SingletonMode()

	// Run cron in every minute
	s.Every(1).Minute().Do(func() {
		if len(Orders_queue) > 0 {
			processOrder()
		}
	})
	s.StartBlocking()
}

func processOrder() {
	var orders []model.Order
	if len(Orders_queue) <= 10 {
		orders = append(orders, Orders_queue...)
		Orders_queue = nil
	} else {
		for i := 0; i < 10; i++ {
			orders = append(orders, Orders_queue[i])
		}
		Orders_queue = append(Orders_queue[:0], Orders_queue[10:]...)
	}

	for i := 0; i < len(orders); i++ {
		initializer.DB.Transaction(func(tx *gorm.DB) error {
			var merchant model.Merchant
			if err := tx.First(&merchant, 1).Error; err != nil {
				return err
			}
			if merchant.Stock > 0 {
				confirmOrder(orders[i], merchant, tx)
			} else {
				refundOrder(orders[i], merchant, tx)
			}
			return nil
		})

	}
}

func confirmOrder(order model.Order, merchant model.Merchant, tx *gorm.DB) {
	initializer.DB.Transaction(func(tx2 *gorm.DB) error {

		if err := tx2.Model(&order).Update("Status", "CONFIRMED").Error; err != nil {
			return err
		}

		if err := tx2.Model(&merchant).Update("Balance", merchant.Balance+order.Amount).Error; err != nil {
			return err
		}

		if err := tx2.Model(&merchant).Update("Stock", merchant.Stock-1).Error; err != nil {
			return err
		}
		return nil

	})
}

func refundOrder(order model.Order, merchant model.Merchant, tx *gorm.DB) {
	initializer.DB.Transaction(func(tx2 *gorm.DB) error {

		var user model.User

		if err := tx2.First(&user, order.UserId).Error; err != nil {
			return err
		}

		if err := tx2.Model(&user).Update("Balance", user.Balance+order.Amount).Error; err != nil {
			return err
		}

		if err := tx2.Model(&order).Update("Status", "REFUNDED").Error; err != nil {
			return err
		}

		// return nil will commit the whole transaction
		return nil

	})
}
