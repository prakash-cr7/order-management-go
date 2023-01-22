package initializer

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect (dsn string)  {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if(err != nil){
		panic(err)
	}

	DB = db
}