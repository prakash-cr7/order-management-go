package model

import "gorm.io/gorm"

type Merchant struct {
	gorm.Model 
	Balance uint
	Stock uint
}