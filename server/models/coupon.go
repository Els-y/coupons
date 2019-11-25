package models

import (
	"github.com/Els-y/coupons/server/pkgs/setting"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type Coupon struct {
	ID          uint   `gorm:"primary_key" json:"-"`
	Username    string `gorm:"type:varchar(20)" json:"username"`
	Coupons     string `gorm:"type:varchar(60)" json:"coupons"`
	Description string `gorm:"type:varchar(60)" json:"description"`
	Stock       int    `json:"stock"`
	Amount      int    `json:"amount"`
	Left        int    `json:"left"`
}

func GetCoupon(username, name string) (*Coupon, error) {
	var coupon Coupon
	err := db.Where(&Coupon{Username: username, Coupons: name}).First(&coupon).Error
	return &coupon, err
}

func GetCouponsWithPage(username string, isSaler bool, page int) ([]Coupon, error) {
	var coupons = make([]Coupon, setting.AppSetting.PageSize)
	query := db.Where(&Coupon{Username: username})
	if isSaler {
		query = query.Where("`left` > 0")
	}

	err := query.Offset(setting.AppSetting.PageSize * (page - 1)).Limit(setting.AppSetting.PageSize).Find(&coupons).Error
	return coupons, err
}

func AddCoupon(username, name, description string, stock, amount int) error {
	coupon, err := GetCoupon(username, name)
	if gorm.IsRecordNotFoundError(err) {
		coupon = &Coupon{
			Username:    username,
			Coupons:     name,
			Description: description,
			Stock:       stock,
			Amount:      amount,
			Left:        amount,
		}
		return db.Create(coupon).Error
	}

	if err != nil {
		return err
	}

	err = db.Exec("UPDATE coupon SET `amount`=`amount`+?, `left`=`left`+? WHERE `username`=? AND `coupons`=?",
		amount, amount, username, name).Error
	return err
}

func AssignCoupon(salerName, customerName string, coupon *Coupon) (bool, error) {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return false, err
	}

	query := tx.Exec("UPDATE coupon SET `left`=`left`-1 WHERE `username`=? AND `coupons`=? AND `left`>1", salerName, coupon.Coupons)
	err := query.Error
	rowsAffected := query.RowsAffected
	if err != nil || rowsAffected != 1 {
		logrus.Infof("[models.AssignCoupon] salerName: %v, customerName: %v, couponName: %v, rowsAffected: %v",
			salerName, customerName, coupon.Coupons, rowsAffected)
		tx.Rollback()
		return false, err
	}

	err = tx.Create(&Coupon{
		Username:    customerName,
		Coupons:     coupon.Coupons,
		Description: "",
		Stock:       coupon.Stock,
		Amount:      1,
		Left:        1,
	}).Error
	if err != nil {
		tx.Rollback()
		return false, err
	}

	tx.Commit()
	return true, err
}
