package api

import (
	"encoding/json"
	"github.com/Els-y/coupons/server/models"
	"github.com/Els-y/coupons/server/pkgs/redis"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

func GetUserWithCache(username string) (*models.User, error) {
	logger := logrus.WithFields(logrus.Fields{"func": "utils.GetUserWithCache", "username": username})
	key := redis.GenUserKey(username)
	userBytes, err := redis.Get(key)
	if err == nil {
		user, err := UserBytesToStruct(userBytes)
		if err == nil {
			logger.Info("user exists in redis")
			return user, nil
		} else {
			str, err := BytesToString(userBytes)
			if err == nil && str == redis.NULL {
				logger.Info("user not exists, get null from redis")
				return nil, nil
			}
		}
	}

	user, err := models.GetUser(username)
	if gorm.IsRecordNotFoundError(err) {
		logger.Info("models.GetUser user not exists")
		err = redis.Set(key, redis.NULL, 60)
		if err != nil {
			logger.WithError(err).Warn("cache null fail")
		} else {
			logger.Info("cache null success")
		}
		return nil, nil
	}
	if err != nil {
		logger.WithError(err).Warn("models.GetUser db error")
		return nil, err
	}

	err = redis.Set(key, user, 5*60)
	if err != nil {
		logger.WithError(err).Warn("cache user fail")
	} else {
		logger.Info("cache user success")
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

func BytesToString(bytes []byte) (string, error) {
	var str string
	err := json.Unmarshal(bytes, &str)
	if err != nil {
		return "", err
	}
	return str, nil
}

func GetCouponWithCache(username, couponName string) (*models.Coupon, error) {
	logger := logrus.WithFields(logrus.Fields{
		"func":     "utils.GetCouponWithCache",
		"username": username,
		"coupon":   couponName,
	})
	key := redis.GenCouponKey(username, couponName)
	couponBytes, err := redis.Get(key)
	if err == nil {
		coupon, err := CouponBytesToStruct(couponBytes)
		if err == nil {
			logger.Info("coupon exists in redis")
			return coupon, nil
		} else {
			str, err := BytesToString(couponBytes)
			if err == nil && str == redis.NULL {
				logger.Info("coupon not exists, get null from redis")
				return nil, nil
			}
		}
	}

	coupon, err := models.GetCoupon(username, couponName)
	if gorm.IsRecordNotFoundError(err) {
		logger.Info("models.GetCoupon coupon not exists")
		err = redis.Set(key, redis.NULL, 60)
		if err != nil {
			logger.WithError(err).Warn("cache null fail")
		} else {
			logger.Info("cache null success")
		}
		return nil, nil
	}
	if err != nil {
		logger.WithError(err).Warn("models.GetCoupon db error")
		return nil, err
	}

	err = redis.Set(key, coupon, 5*60)
	if err != nil {
		logger.WithError(err).Warn("cache coupon fail")
	} else {
		logger.Info("cache coupon success")
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

func CheckIfUserHasCoupon(username, couponOwnerName, couponName string) (bool, error) {
	logger := logrus.WithFields(logrus.Fields{
		"func":        "utils.CheckIfUserHasCoupon",
		"username":    username,
		"couponOwner": couponOwnerName,
		"coupon":      couponName,
	})
	key := redis.GenCouponOwnersKey(couponOwnerName, couponName)
	exist, err := redis.SIsmember(key, username)
	if err == nil && exist == true {
		logger.Info("user has the coupon")
		return true, nil
	}

	coupon, err := GetCouponWithCache(username, couponName)
	if err != nil {
		logger.WithError(err).Warn("get coupon with cache fail")
		return false, err
	}
	if coupon == nil {
		return false, nil
	}

	_, err = redis.SAdd(key, username)
	if err != nil {
		logger.WithError(err).Warn("update coupon owners fail")
	} else {
		logger.Info("update coupon owners success")
	}

	return true, nil
}
