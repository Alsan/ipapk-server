package models

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/phinexdaz/ipapk-server/conf"
)

var orm *gorm.DB

func InitDB() error {
	var err error
	orm, err = gorm.Open("mysql", conf.AppConfig.Database)
	if err != nil {
		return err
	}
	if gin.Mode() != "release" {
		orm.LogMode(true)
	}
	orm.AutoMigrate(&Bundle{})
	return nil
}
