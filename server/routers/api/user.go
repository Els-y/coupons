package api

import (
	"github.com/Els-y/coupons/server/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AddUserReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Kind     int    `json:"kind"`
}

func AddUser(ctx *gin.Context) {
	var req AddUserReq

	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.JSON(400, gin.H{
			"errMsg": "params error",
		})
		return
	}

	exists, err := models.ExistsUsername(req.Username)
	if err != nil {
		ctx.JSON(400, gin.H{
			"errMsg": "db error",
		})
		return
	}

	if exists {
		ctx.JSON(400, gin.H{
			"errMsg": "username exists",
		})
		return
	}

	err = models.AddUser(req.Username, req.Password, req.Kind)
	if err != nil {
		logrus.Errorf("AddUser: err= %+v", err)
		ctx.JSON(400, gin.H{
			"errMsg": "db error",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"errMsg": "",
	})
}

func ExistsUser(ctx *gin.Context) {
	username := ctx.Query("username")

	exists, err := models.ExistsUsername(username)
	if err != nil {
		ctx.JSON(400, gin.H{
			"errMsg": "db error",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"exists": exists,
		"errMsg": "",
	})
}
