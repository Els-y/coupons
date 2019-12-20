package api

import (
	"github.com/Els-y/coupons/server/models"
	"github.com/Els-y/coupons/server/pkgs/redis"
	"github.com/Els-y/coupons/server/pkgs/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/nats-io/nats.go"
	"strconv"
)

type AddCouponsReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Stock       int    `json:"stock"`
	Amount      int    `json:"amount"`
}

type SubscribeAssignCoupon struct {
	SalerName string
	TokenUsername string
	CouponName string
	CouponStock int
}

func init() {
	models.Nc, models.Nat_err = nats.Connect(models.NatsUrl)
	if models.Nat_err != nil {
		logrus.Infof("[queue.init] subscribe nats.Connect url error, url: %v, err: %v", models.NatsUrl, models.Nat_err.Error())
		return
	}

	models.NatsEncodedConn, models.Nat_err = nats.NewEncodedConn(models.Nc, nats.JSON_ENCODER)
	if models.Nat_err != nil {
		logrus.Infof("[queue.init] subscribe nats.NewEncodedConn error, err: %v", models.Nat_err.Error())
		return
	}
}


func AddCoupons(ctx *gin.Context) {
	var req AddCouponsReq
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{
			"errMsg": "params error",
		})
		return
	}

	username := ctx.Param("username")
	token := ctx.GetHeader("Authorization")
	if !redis.Exists(redis.GenAuthorizationKey(token)) {
		ctx.JSON(401, gin.H{
			"errMsg": "token not exists",
		})
		return
	}

	tokenUsername, tokenKindStr, err := utils.DecodeToken(ctx.GetHeader("Authorization"))
	if err != nil {
		ctx.JSON(401, gin.H{
			"errMsg": "authorization error",
		})
		return
	}

	if username != tokenUsername || tokenKindStr != models.KindSalerStr {
		ctx.JSON(401, gin.H{
			"errMsg": "permission denied",
		})
		return
	}

	err = models.AddCoupon(username, req.Name, req.Description, req.Stock, req.Amount)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"saler":  username,
			"coupon": req.Name,
			"err":    err,
		}).Warn("[api.AddCoupons] models.AddCoupon db error")
		ctx.JSON(400, gin.H{
			"errMsg": "db error",
		})
		return
	}

	ok, err := redis.SAdd(redis.GenCouponOwnersKey(req.Name), username)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"saler":  username,
			"coupon": req.Name,
			"ok":     ok,
			"err":    err,
		}).Warn("[api.AddCoupons] redis.SAdd error")
		ctx.JSON(400, gin.H{
			"errMsg": "redis error",
		})
	}

	_, err = redis.IncrBy(redis.GenCouponLeftKey(username, req.Name), req.Amount)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"saler":  username,
			"coupon": req.Name,
			"err":    err,
		}).Warn("[api.AddCoupons] redis set left error")
		ctx.JSON(400, gin.H{
			"errMsg": "redis error",
		})
		return
	}

	ctx.JSON(201, gin.H{
		"errMsg": "",
	})
}

func GetCouponsInfo(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil {
		ctx.JSON(400, gin.H{
			"errMsg": "params error",
			"data":   []models.Coupon{},
		})
		return
	}

	username := ctx.Param("username")
	token := ctx.GetHeader("Authorization")
	if !redis.Exists(redis.GenAuthorizationKey(token)) {
		ctx.JSON(401, gin.H{
			"errMsg": "token not exists",
			"data":   []models.Coupon{},
		})
		return
	}

	tokenUsername, tokenKindStr, err := utils.DecodeToken(ctx.GetHeader("Authorization"))
	if err != nil {
		ctx.JSON(401, gin.H{
			"errMsg": "authorization error",
			"data":   []models.Coupon{},
		})
		return
	}

	var coupons []models.Coupon
	if username == tokenUsername {
		if tokenKindStr == models.KindCustomerStr {
			coupons, err = models.GetCouponsWithPage(username, false, page)
		} else {
			coupons, err = models.GetCouponsWithPage(username, true, page)
		}
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"username": username,
				"page":     page,
				"err":      err,
			}).Warn("[api.GetCouponsInfo] models.GetCouponsWithPage db error")
			ctx.JSON(400, gin.H{
				"errMsg": "",
				"data":   []models.Coupon{},
			})
			return
		}
		if coupons == nil || len(coupons) == 0 {
			ctx.JSON(204, gin.H{
				"errMsg": "query null",
				"data": []models.Coupon{},
			})
			return
		}
		ctx.JSON(200, gin.H{
			"errMsg": "",
			"data": coupons,
		})
		return
	}

	userKindStr, err := GetUserKindWithCache(username)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"username": username,
			"err":      err,
		}).Warn("[api.GetCouponsInfo] GetUserKindWithCache db error")
		ctx.JSON(400, gin.H{
			"errMsg": "db error",
			"data":   []models.Coupon{},
		})
		return
	}

	if userKindStr != models.KindSalerStr {
		logrus.WithFields(logrus.Fields{
			"username": username,
			"kind":     userKindStr,
		}).Info("[api.GetCouponsInfo] GetUserKindWithCache user is not a saler")
		ctx.JSON(401, gin.H{
			"errMsg": "user is not a saler",
			"data":   []models.Coupon{},
		})
		return
	}

	coupons, err = models.GetCouponsWithPage(username, true, page)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"username": username,
			"page":     page,
			"err":      err,
		}).Warn("[api.GetCouponsInfo] models.GetCouponsWithPage db error")
		ctx.JSON(400, gin.H{
			"errMsg": "db error",
			"data":   []models.Coupon{},
		})
		return
	}
	if coupons == nil || len(coupons) == 0 {
		ctx.JSON(204, gin.H{
			"errMsg": "query null",
			"data": []models.Coupon{},
		})
		return
	}
	ctx.JSON(200, gin.H{
		"errMsg": "",
		"data": coupons,
	})
}

func AssignCoupon(ctx *gin.Context) {
	salerName := ctx.Param("username")
	couponName := ctx.Param("name")

	token := ctx.GetHeader("Authorization")
	if !redis.Exists(redis.GenAuthorizationKey(token)) {
		ctx.JSON(401, gin.H{
			"errMsg": "token not exists",
		})
		return
	}

	tokenUsername, _, err := utils.DecodeToken(ctx.GetHeader("Authorization"))
	if err != nil {
		ctx.JSON(401, gin.H{
			"errMsg": "authorization error",
		})
		return
	}

	salerKindStr, err := GetUserKindWithCache(salerName)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"saler": salerName,
			"err":   err,
		}).Warn("[api.AssignCoupon] GetUserKindWithCache db error")
		ctx.JSON(204, gin.H{
			"errMsg": "db error",
		})
		return
	}
	if salerKindStr != models.KindSalerStr {
		logrus.WithFields(logrus.Fields{
			"saler": salerName,
			"kind":  salerKindStr,
		}).Info("[api.AssignCoupon] GetUserKindWithCache user is not a saler")
		ctx.JSON(204, gin.H{
			"errMsg": "user is not a saler",
		})
		return
	}

	exist, couponStock, err := CheckIfUserHasCoupon(tokenUsername, couponName)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"customer": tokenUsername,
			"coupon":   couponName,
			"err":      err,
		}).Warn("[api.AssignCoupon] CheckIfUserHasCoupon db error")
		ctx.JSON(204, gin.H{
			"errMsg": "db error",
		})
		return
	}
	if exist == true {
		logrus.WithFields(logrus.Fields{
			"customer": tokenUsername,
			"coupon":   couponName,
		}).Info("[api.AssignCoupon] customer already has the coupon")
		ctx.JSON(204, gin.H{
			"errMsg": "user already have this coupon",
		})
		return
	}

	exist, couponStock, err = CheckIfUserHasCoupon(salerName, couponName)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"saler":  salerName,
			"coupon": couponName,
			"err":    err,
		}).Warn("[api.AssignCoupon] CheckIfUserHasCoupon db error")
		ctx.JSON(204, gin.H{
			"errMsg": "db error",
		})
		return
	}
	if exist == false {
		logrus.WithFields(logrus.Fields{
			"saler":  salerName,
			"coupon": couponName,
		}).Info("[api.AssignCoupon] coupon not exists")
		ctx.JSON(204, gin.H{
			"errMsg": "coupon not exists",
		})
		return
	}

	left, err := redis.Decr(redis.GenCouponLeftKey(salerName, couponName))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"saler":    salerName,
			"customer": tokenUsername,
			"coupon":   couponName,
			"err":      err,
		}).Warn("[api.AssignCoupon] redis error")
		ctx.JSON(204, gin.H{
			"errMsg": "redis decr fail",
		})
		return
	}

	if left < 0 {
		ctx.JSON(204, gin.H{
			"errMsg": "no coupons left",
		})
		return
	}

	ok, err := redis.SAdd(redis.GenCouponOwnersKey(couponName), tokenUsername)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"customer": tokenUsername,
			"coupon":   couponName,
			"ok":       ok,
		}).Warn("[api.AssignCoupon] redis.SAdd error")
	}

	// push msg to mq
	subscribeAssignCoupon := &(SubscribeAssignCoupon{SalerName: salerName, TokenUsername: tokenUsername, CouponName: couponName, CouponStock: couponStock})
	models.NatsEncodedConn.Publish(models.AssignCoupon_Subj, subscribeAssignCoupon)

	ctx.JSON(201, gin.H{
		"errMsg": "",
	})
}
