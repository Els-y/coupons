package api

import (
	"github.com/Els-y/coupons/server/models"
	"github.com/Els-y/coupons/server/pkgs/redis"
	"github.com/Els-y/coupons/server/pkgs/utils"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"strconv"
)

type AddCouponsReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Stock       int    `json:"stock"`
	Amount      int    `json:"amount"`
}

type SubscribeAssignCoupon struct {
	SalerName     string
	TokenUsername string
	CouponName    string
	CouponStock   int
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

	logger := logrus.WithFields(logrus.Fields{"func": "api.AddCoupons", "user": username, "coupon": req.Name})
	coupon, err := models.AddCoupon(username, req.Name, req.Description, req.Stock, req.Amount)
	if err != nil {
		logger.WithError(err).Warn("models.AddCoupon db error")
		ctx.JSON(400, gin.H{
			"errMsg": "db error",
		})
		return
	}

	err = redis.Set(redis.GenCouponKey(username, req.Name), coupon, -1)
	if err != nil {
		logger.WithError(err).Warn("cache coupon fail")
	}

	ok, err := redis.SAdd(redis.GenCouponOwnersKey(username, req.Name), username)
	if err != nil {
		logger.WithFields(logrus.Fields{"ok": ok}).WithError(err).Warn("cache coupon owners fail")
	}

	_, err = redis.IncrBy(redis.GenCouponLeftKey(username, req.Name), req.Amount)
	if err != nil {
		logger.WithError(err).Warn("redis incr left error")
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

	logger := logrus.WithFields(logrus.Fields{"func": "api.GetCouponsInfo", "username": username, "page": page})
	user, err := GetUserWithCache(username)
	if err != nil {
		logger.WithError(err).Warn("GetUserWithCache db error")
		ctx.JSON(400, gin.H{
			"errMsg": "db error",
			"data":   []models.Coupon{},
		})
		return
	}
	if user == nil {
		ctx.JSON(204, gin.H{
			"errMsg": "query null",
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
			logger.WithError(err).Warn("models.GetCouponsWithPage db error")
			ctx.JSON(400, gin.H{
				"errMsg": "",
				"data":   []models.Coupon{},
			})
			return
		}
		if coupons == nil || len(coupons) == 0 {
			ctx.JSON(204, gin.H{
				"errMsg": "query null",
				"data":   []models.Coupon{},
			})
			return
		}
		ctx.JSON(200, gin.H{
			"errMsg": "",
			"data":   coupons,
		})
		return
	}

	if user.Kind != models.KindSalerStr {
		ctx.JSON(401, gin.H{
			"errMsg": "user is not a saler",
			"data":   []models.Coupon{},
		})
		return
	}

	coupons, err = models.GetCouponsWithPage(username, true, page)
	if err != nil {
		logger.WithError(err).Warn("models.GetCouponsWithPage db error")
		ctx.JSON(400, gin.H{
			"errMsg": "db error",
			"data":   []models.Coupon{},
		})
		return
	}
	if coupons == nil || len(coupons) == 0 {
		ctx.JSON(204, gin.H{
			"errMsg": "query null",
			"data":   []models.Coupon{},
		})
		return
	}
	ctx.JSON(200, gin.H{
		"errMsg": "",
		"data":   coupons,
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

	logger := logrus.WithFields(
		logrus.Fields{"func": "api.AssignCoupon", "saler": salerName, "customer": tokenUsername, "coupon": couponName})
	saler, err := GetUserWithCache(salerName)
	if err != nil {
		logger.WithError(err).Warn("get saler error")
		ctx.JSON(204, gin.H{
			"errMsg": "db error",
		})
		return
	}
	if saler == nil || saler.Kind != models.KindSalerStr {
		ctx.JSON(204, gin.H{
			"errMsg": "user is not a saler",
		})
		return
	}

	exist, err := CheckIfUserHasCoupon(tokenUsername, salerName, couponName)
	if err != nil {
		logger.WithError(err).Warn("check whether customer has coupon fail")
		ctx.JSON(204, gin.H{
			"errMsg": "db error",
		})
		return
	}
	if exist == true {
		ctx.JSON(204, gin.H{
			"errMsg": "user already have this coupon",
		})
		return
	}

	coupon, err := GetCouponWithCache(salerName, couponName)
	if err != nil {
		logger.WithError(err).Warn("get coupon fail")
		ctx.JSON(204, gin.H{
			"errMsg": "db error",
		})
		return
	}
	if coupon == nil {
		ctx.JSON(204, gin.H{
			"errMsg": "coupon not exists",
		})
		return
	}

	left, err := redis.Decr(redis.GenCouponLeftKey(salerName, couponName))
	if err != nil {
		logger.WithError(err).Warn("decr coupon left error")
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

	ok, err := redis.SAdd(redis.GenCouponOwnersKey(salerName, couponName), tokenUsername)
	if err != nil {
		logger.WithFields(logrus.Fields{"ok": ok}).WithError(err).Warn("update coupon owners fail")
	}

	// push msg to mq
	subscribeAssignCoupon := &(SubscribeAssignCoupon{
		SalerName: salerName, TokenUsername: tokenUsername, CouponName: couponName, CouponStock: coupon.Stock})
	_ = models.NatsEncodedConn.Publish(models.AssignCoupon_Subj, subscribeAssignCoupon)

	ctx.JSON(201, gin.H{
		"errMsg": "",
	})
}
