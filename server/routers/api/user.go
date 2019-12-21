package api

import (
	"github.com/Els-y/coupons/server/models"
	"github.com/Els-y/coupons/server/pkgs/redis"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AddUserReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Kind     string `json:"kind"`
}

func AddUser(ctx *gin.Context) {
	var req AddUserReq
	err := ctx.BindJSON(&req)
	if err != nil {
		logrus.WithError(err).Warn("[api.AddUser] ctx.BindJSON error")
		ctx.JSON(400, gin.H{
			"errMsg": "params error",
		})
		return
	}

	logger := logrus.WithFields(logrus.Fields{"func": "api.AddUser", "username": req.Username, "kind": req.Kind})
	if req.Kind != models.KindCustomerStr && req.Kind != models.KindSalerStr {
		logger.Info("kind is not customer or saler")
		ctx.JSON(400, gin.H{
			"errMsg": "params kind error",
		})
		return
	}

	user, err := GetUserWithCache(req.Username)
	if err != nil {
		logger.WithError(err).Warn("GetUserKindWithCache db error")
		ctx.JSON(400, gin.H{
			"errMsg": "db error",
		})
		return
	}

	if user != nil {
		ctx.JSON(400, gin.H{
			"errMsg": "username exists",
		})
		return
	}

	user, err = models.AddUser(req.Username, req.Password, req.Kind)
	if err != nil {
		logger.WithError(err).Warn("models.AddUser db error")
		ctx.JSON(400, gin.H{
			"errMsg": "db error",
		})
		return
	}

	err = redis.Set(redis.GenUserKey(req.Username), user, 5*60)
	if err != nil {
		logger.WithError(err).Warn("cache user fail")
	} else {
		logger.Info("cache user success")
	}

	ctx.JSON(201, gin.H{
		"errMsg": "",
	})
}

func ExistsUser(ctx *gin.Context) {
	username := ctx.Query("username")

	user, err := GetUserWithCache(username)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"username": username,
			"err":      err,
		}).Warn("[api.ExistsUser] GetUserWithCache db error")
		ctx.JSON(400, gin.H{
			"errMsg": "db error",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"exists": user != nil,
		"errMsg": "",
	})
}
