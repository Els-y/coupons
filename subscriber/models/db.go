package models

import (
	"fmt"
	"github.com/sirupsen/logrus"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/Els-y/coupons/subscriber/pkgs/setting"
)

var db *gorm.DB

// Setup initializes the database instance
func Setup() {
	var err error
	db, err = gorm.Open(setting.DatabaseSetting.Type, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.Name))

	if err != nil {
		logrus.Fatalf("models.Setup err: %v", err)
	}

	db.SingularTable(true)
	db.DB().SetMaxIdleConns(1024)
	db.DB().SetMaxOpenConns(4096)

	db.AutoMigrate(&User{}, &Coupon{})
}

// CloseDB closes database connection (unnecessary)
func CloseDB() {
	defer db.Close()
}
