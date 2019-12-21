package models

import (
	"github.com/Els-y/coupons/server/pkgs/utils"
	"github.com/jinzhu/gorm"
)

type User struct {
	ID       uint   `gorm:"primary_key"`
	Username string `gorm:"type:varchar(20)" json:"username"`
	Password string `gorm:"type:varchar(32)" json:"password"`
	Kind     string `json:"kind"`
}

func ExistsUsername(username string) (bool, error) {
	var user User
	err := db.Where(&User{Username: username}).First(&user).Error
	if gorm.IsRecordNotFoundError(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func AddUser(username, password, kind string) (*User, error) {
	user := User{
		Username: username,
		Password: utils.MD5Encode(password),
		Kind:     kind,
	}

	if err := db.Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUser(username string) (*User, error) {
	var user User
	err := db.Where(&User{Username: username}).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserWithPwd(username, password string) (*User, error) {
	var user User
	err := db.Where(&User{Username: username, Password: utils.MD5Encode(password)}).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func CheckPwd(password, target string) bool {
	return utils.MD5Encode(password) == target
}
