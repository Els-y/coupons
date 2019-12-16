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

	user, err := GetUserWithCache(req.Username)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"username": req.Username,
			"err":      err,
		}).Warn("[api.AddUser] GetUserWithCache db error")
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

	err = models.AddUser(req.Username, req.Password, req.Kind)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"username": req.Username,
			"err":      err,
		}).Warn("[api.AddUser] models.AddUser db error")
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
