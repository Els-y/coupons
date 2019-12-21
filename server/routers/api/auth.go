package api

import (
	"github.com/Els-y/coupons/server/models"
	"github.com/Els-y/coupons/server/pkgs/redis"
	"github.com/Els-y/coupons/server/pkgs/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AuthReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Auth(ctx *gin.Context) {
	var req AuthReq
	err := ctx.BindJSON(&req)
	ctx.Header("Authorization", "")

	if err != nil {
		ctx.JSON(401, gin.H{
			"kind":   "",
			"errMsg": "params error",
		})
		return
	}

	logger := logrus.WithFields(logrus.Fields{"func": "api.Auth", "username": req.Username})

	user, err := GetUserWithCache(req.Username)
	if err != nil {
		logger.WithError(err).Warn("GetUserKindWithCache db error")
		ctx.JSON(401, gin.H{
			"kind":   "",
			"errMsg": "db error",
		})
		return
	}
	if user == nil || !models.CheckPwd(req.Password, user.Password) {
		ctx.JSON(401, gin.H{
			"kind":   "",
			"errMsg": "username or password wrong",
		})
		return
	}

	token := utils.EncodeToken(user.Username, user.Kind)
	err = redis.Set(redis.GenAuthorizationKey(token), 1, -1)
	if err != nil {
		logger.WithError(err).Warn("redis error")
		ctx.JSON(401, gin.H{
			"kind":   "",
			"errMsg": "redis error",
		})
		return
	}

	ctx.Header("Authorization", token)
	ctx.JSON(200, gin.H{
		"kind":   user.Kind,
		"errMsg": "",
	})
}
