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
			logrus.Infof("[utils.GetUser] user exists in redis, user: %+v", user)
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

	err = redis.Set(key, user, 5*60)
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
			logrus.Infof("[utils.GetCouponWithCache] coupon exists in redis, coupon: %+v", coupon)
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

	err = redis.Set(key, coupon, 5*60)
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

func CheckIfUserHasCoupon(username, couponName string) (bool, error) {
	key := redis.GenCouponOwnersKey(couponName)
	exist, err := redis.SIsmember(key, username)
	if err == nil {
		if exist == true {
			logrus.Infof("[utils.CheckIfUserHasCoupon] user has the coupon, username: %v, coupon: %v", username, couponName)
			return true, nil
		}
	}

	_, err = models.GetCoupon(username, couponName)
	if gorm.IsRecordNotFoundError(err) {
		logrus.Infof("[utils.CheckIfUserHasCoupon] models.GetCoupon coupon not exists, username: %v, couponName: %v, err: %v",
			username, couponName, err)
		return false, nil
	}
	if err != nil {
		logrus.Infof("[utils.CheckIfUserHasCoupon] models.GetCoupon db error, username: %v, couponName: %v, err: %v",
			username, couponName, err)
		return false, err
	}

	_, err = redis.SAdd(key, username)
	if err != nil {
		logrus.Infof("[utils.CheckIfUserHasCoupon] redis.SAdd error, username: %v, couponName: %v, err: %v", username, couponName, err)
	} else {
		logrus.Infof("[utils.CheckIfUserHasCoupon] redis.SAdd success, username: %v, couponName: %v", username, couponName)
	}

	return true, nil
}
