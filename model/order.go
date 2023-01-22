package model

import (
	"gorm.io/gorm"
)

// status = CREATED, CONFIRMED, CANCELLED, REFUNDED

type Order struct {
	gorm.Model
	Status string
	Amount uint
	UserId uint 
	User User `json:"-" gorm:"foreignKey:UserId"` 
}