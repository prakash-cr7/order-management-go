package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name string
	Balance uint 
}