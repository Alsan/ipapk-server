package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var orm *gorm.DB

func InitDB() {
	var err error
	orm, err = gorm.Open("sqlite3", "app.db?loc=Asia/Shanghai")
	if err != nil {
		log.Fatal(err)
	}
	orm.AutoMigrate(&Bundle{})
}
