package api

import (
	"encoding/json"
	"github.com/Els-y/coupons/server/models"
	"github.com/Els-y/coupons/server/pkgs/redis"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

func GetUserWithCache(username string) (*models.User, error) {
	logger := logrus.WithFields(logrus.Fields{"username": username})
	key := redis.GenUserKey(username)
	userBytes, err := redis.Get(key)
	if err == nil {
		user, err := UserBytesToStruct(userBytes)
		if err == nil {
			logger.Info("[utils.GetUserWithCache] user exists in redis")
			return user, nil
		}
	}

	user, err := models.GetUser(username)
	if gorm.IsRecordNotFoundError(err) {
		logger.Info("[utils.GetUserWithCache] models.GetUser user not exists")
		return nil, nil
	}
	if err != nil {
		logger.WithError(err).Warn("[utils.GetUserWithCache] models.GetUser db error")
		return nil, err
	}

	err = redis.Set(key, user, 5*60)
	if err != nil {
		logger.WithError(err).Warn("[utils.GetUserWithCache] redis.Set error")
	} else {
		logger.Info("[utils.GetUserWithCache] redis.Set success")
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
	logger := logrus.WithFields(logrus.Fields{
		"username": username,
		"coupon":   couponName,
	})
	key := redis.GenCouponKey(username, couponName)
	couponBytes, err := redis.Get(key)
	if err == nil {
		coupon, err := CouponBytesToStruct(couponBytes)
		if err == nil {
			logger.Info("[utils.GetCouponWithCache] coupon exists in redis")
			return coupon, nil
		}
	}

	coupon, err := models.GetCoupon(username, couponName)
	if gorm.IsRecordNotFoundError(err) {
		logger.Info("[utils.GetCouponWithCache] models.GetCoupon coupon not exists")
		return nil, nil
	}
	if err != nil {
		logger.WithError(err).Warn("[utils.GetCouponWithCache] models.GetCoupon db error")
		return nil, err
	}

	err = redis.Set(key, coupon, 5*60)
	if err != nil {
		logger.WithError(err).Warn("[utils.GetCouponWithCache] redis.Set error")
	} else {
		logger.Info("[utils.GetCouponWithCache] redis.Set success")
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
	logger := logrus.WithFields(logrus.Fields{
		"username": username,
		"coupon":   couponName,
	})
	key := redis.GenCouponOwnersKey(couponName)
	exist, err := redis.SIsmember(key, username)
	if err == nil {
		if exist == true {
			logger.Info("[utils.CheckIfUserHasCoupon] user has the coupon")
			return true, nil
		}
	}

	_, err = models.GetCoupon(username, couponName)
	if gorm.IsRecordNotFoundError(err) {
		logger.Info("[utils.CheckIfUserHasCoupon] models.GetCoupon coupon not exists")
		return false, nil
	}
	if err != nil {
		logger.WithError(err).Warn("[utils.CheckIfUserHasCoupon] models.GetCoupon db error")
		return false, err
	}

	_, err = redis.SAdd(key, username)
	if err != nil {
		logger.WithError(err).Warn("[utils.CheckIfUserHasCoupon] redis.SAdd error")
	} else {
		logger.Info("[utils.CheckIfUserHasCoupon] redis.SAdd success")
	}

	return true, nil
}

func GetUserKindWithCache(username string) (string, error) {
	logger := logrus.WithFields(logrus.Fields{"func": "utils.GetUserKindWithCache", "username": username})
	key := redis.GenUserKindKey(username)
	kindBytes, err := redis.Get(key)
	if err == nil && kindBytes != nil {
		kindStr, err := UserKindBytesToString(kindBytes)
		if err == nil {
			logger.WithFields(logrus.Fields{"kind": kindStr}).Info("get user kind in redis success")
			return kindStr, nil
		}
	}

	user, err := models.GetUser(username)
	if gorm.IsRecordNotFoundError(err) {
		logger.Info("models.GetUser user not exists")
		return "", nil
	}
	if err != nil {
		logger.WithError(err).Warn("models.GetUser db error")
		return "", err
	}

	kindStr := models.KindInt2Str[user.Kind]
	err = redis.Set(key, kindStr, 5*60)
	if err != nil {
		logger.WithError(err).Warn("redis.Set error")
	} else {
		logger.Info("redis.Set success")
	}

	return kindStr, nil
}

func UserKindBytesToString(kindBytes []byte) (string, error) {
	var kindStr string
	err := json.Unmarshal(kindBytes, &kindStr)
	if err != nil {
		return "", err
	}
	return kindStr, nil
}
