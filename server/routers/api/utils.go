package api

import (
	"encoding/json"
	"github.com/Els-y/coupons/server/models"
	"github.com/Els-y/coupons/server/pkgs/redis"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

func GetUserWithCache(username string) (*models.User, error) {
	key := redis.GenUserKey(username)
	userBytes, err := redis.Get(key)
	if err == nil {
		user, err := UserBytesToStruct(userBytes)
		if err == nil {
			logrus.Infof("[utils.GetUser] user exists in redis, user: %+v, err: %v", user, err)
			return user, nil
		}
	}

	user, err := models.GetUser(username)
	if gorm.IsRecordNotFoundError(err) {
		logrus.Infof("[utils.GetUser] models.GetUser user not exists, username: %v, err: %v", username, err)
		return nil, nil
	}
	if err != nil {
		logrus.Infof("[utils.GetUser] models.GetUser db error, username: %v, err: %v", username, err)
		return nil, err
	}

	err = redis.Set(key, user)
	if err != nil {
		logrus.Infof("[utils.GetUser] redis.Set error, username: %v, err: %v", username, err)
	} else {
		logrus.Infof("[utils.GetUser] redis.Set success, username: %v", username)
	}

	return user, nil
}

func UserBytesToStruct(userBytes []byte) (*models.User, error) {
	var user models.User
	err := json.Unmarshal(userBytes, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetCouponWithCache(username, couponName string) (*models.Coupon, error) {
	key := redis.GenCouponKey(username, couponName)
	couponBytes, err := redis.Get(key)
	if err == nil {
		coupon, err := CouponBytesToStruct(couponBytes)
		if err == nil {
			logrus.Infof("[utils.GetCouponWithCache] coupon exists in redis, coupon: %+v, err: %v", coupon, err)
			return coupon, nil
		}
	}

	coupon, err := models.GetCoupon(username, couponName)
	if gorm.IsRecordNotFoundError(err) {
		logrus.Infof("[utils.GetCouponWithCache] models.GetCoupon coupon not exists, username: %v, couponName: %v, err: %v",
			username, couponName, err)
		return nil, nil
	}
	if err != nil {
		logrus.Infof("[utils.GetCouponWithCache] models.GetCoupon db error, username: %v, couponName: %v, err: %v",
			username, couponName, err)
		return nil, err
	}

	err = redis.Set(key, coupon)
	if err != nil {
		logrus.Infof("[utils.GetCouponWithCache] redis.Set error, username: %v, couponName: %v, err: %v", username, couponName, err)
	} else {
		logrus.Infof("[utils.GetCouponWithCache] redis.Set success, username: %v, couponName: %v", username, couponName)
	}

	return coupon, nil
}

func CouponBytesToStruct(couponBytes []byte) (*models.Coupon, error) {
	var coupon models.Coupon
	err := json.Unmarshal(couponBytes, &coupon)
	if err != nil {
		return nil, err
	}
	return &coupon, nil
}
