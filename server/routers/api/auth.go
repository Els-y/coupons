package api

import (
	"github.com/Els-y/coupons/server/models"
	"github.com/Els-y/coupons/server/pkgs/redis"
	"github.com/Els-y/coupons/server/pkgs/utils"
	"github.com/gin-gonic/gin"
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
		ctx.JSON(400, gin.H{
			"kind":   "",
			"errMsg": "params error",
		})
		return
	}

	user, err := models.GetUserWithPwd(req.Username, req.Password)
	if err != nil {
		ctx.JSON(400, gin.H{
			"kind":   "",
			"errMsg": "username or password wrong",
		})
		return
	}

	kindStr := models.KindInt2Str[user.Kind]
	token := utils.EncodeToken(user.Username, kindStr)
	err = redis.Set(redis.GenAuthorizationKey(token), 1)
	if err != nil {
		ctx.JSON(400, gin.H{
			"kind":   "",
			"errMsg": "redis error",
		})
		return
	}

	ctx.Header("Authorization", token)
	ctx.JSON(200, gin.H{
		"kind":   kindStr,
		"errMsg": "",
	})
}
