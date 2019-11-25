package api

import (
	"github.com/Els-y/coupons/server/models"
	"github.com/Els-y/coupons/server/pkgs/redis"
	"github.com/Els-y/coupons/server/pkgs/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"strconv"
)

type AddCouponsReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Stock       int    `json:"stock"`
	Amount      int    `json:"amount"`
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
		ctx.JSON(400, gin.H{
			"errMsg": "authorization error",
		})
		return
	}

	if username != tokenUsername || tokenKindStr == models.KindCustomerStr {
		ctx.JSON(401, gin.H{
			"errMsg": "permission denied",
		})
		return
	}

	err = models.AddCoupon(username, req.Name, req.Description, req.Stock, req.Amount)
	if err != nil {
		logrus.Infof("[api.AddCoupons] models.AddCoupon db error, username: %v, name: %v, err: %v", username, req.Name, err)
		ctx.JSON(400, gin.H{
			"errMsg": "db error",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"errMsg": "",
	})
}

func GetCouponsInfo(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil {
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
		ctx.JSON(400, gin.H{
			"errMsg": "authorization error",
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
			logrus.Infof("[api.GetCouponsInfo] models.GetCouponsWithPage db error, username: %v, page: %v, err: %v", username, page, err)
			ctx.JSON(400, coupons)
			return
		}
		ctx.JSON(200, coupons)
		return
	}

	user, err := models.GetUser(username)
	if err != nil {
		logrus.Infof("[api.GetCouponsInfo] models.GetUser db error, username: %v, err: %v", username, err)
		ctx.JSON(400, gin.H{
			"errMsg": "db error",
		})
		return
	}
	if user.Kind == models.KindCustomerInt {
		ctx.JSON(400, gin.H{
			"errMsg": "user is not a saler",
		})
		return
	}

	coupons, err = models.GetCouponsWithPage(username, true, page)
	if err != nil {
		logrus.Infof("[api.GetCouponsInfo] models.GetCouponsWithPage db error, username: %v, page: %v, err: %v", username, page, err)
		ctx.JSON(400, coupons)
		return
	}
	ctx.JSON(200, coupons)
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
		ctx.JSON(400, gin.H{
			"errMsg": "authorization error",
		})
		return
	}

	user, err := models.GetUser(salerName)
	if err != nil {
		logrus.Infof("[api.AssignCoupon] models.GetUser db error, salerName: %v, err: %v", salerName, err)
		ctx.JSON(400, gin.H{
			"errMsg": "db error",
		})
		return
	}
	if user.Kind == models.KindCustomerInt {
		ctx.JSON(400, gin.H{
			"errMsg": "user is not a saler",
		})
		return
	}

	coupon, err := GetCouponWithCache(tokenUsername, couponName)
	if err != nil {
		logrus.Infof("[api.AssignCoupon] GetCouponWithCache db error, customerName: %v, couponName: %v, err: %v", tokenUsername, couponName, err)
		ctx.JSON(400, gin.H{
			"errMsg": "db error",
		})
		return
	}
	if coupon != nil {
		logrus.Infof("[api.AssignCoupon] customer alread have the coupon, customerName: %v, couponName: %v, err: %v", tokenUsername, couponName, err)
		ctx.JSON(400, gin.H{
			"errMsg": "user already have this coupon",
		})
		return
	}

	coupon, err = GetCouponWithCache(salerName, couponName)
	if err != nil {
		logrus.Infof("[api.AssignCoupon] GetCouponWithCache db error, salerName: %v, couponName: %v, err: %v", salerName, couponName, err)
		ctx.JSON(400, gin.H{
			"errMsg": "db error",
		})
		return
	}
	if coupon == nil {
		logrus.Infof("[api.AssignCoupon] coupon not exists, salerName: %v, couponName: %v, err: %v", salerName, couponName, err)
		ctx.JSON(400, gin.H{
			"errMsg": "coupon not exists",
		})
		return
	}

	ok, err := models.AssignCoupon(salerName, tokenUsername, coupon)
	if err != nil {
		logrus.Infof("[api.AssignCoupon] models.AssignCoupon db error, salerName: %v, customerName: %v, couponName: %v, err: %v",
			salerName, tokenUsername, couponName, err)
		ctx.JSON(400, gin.H{
			"errMsg": "assign db error",
		})
		return
	}
	if !ok {
		ctx.JSON(200, gin.H{
			"errMsg": "no coupons left",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"errMsg": "",
	})
}
