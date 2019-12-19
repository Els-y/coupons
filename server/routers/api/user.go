package api

import (
	"github.com/Els-y/coupons/server/models"
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
		logrus.Infof("[api.AddUser] ctx.BindJSON error, username: %v, err: %v", req.Username, err.Error())
		ctx.JSON(400, gin.H{
			"errMsg": "params error",
		})
		return
	}

	if req.Kind != models.KindCustomerStr && req.Kind != models.KindSalerStr {
		logrus.Infof("[api.AddUser] kind is not customer or saler : %v", req.Kind)
		ctx.JSON(400, gin.H{ 
			"errMsg": "params kind error",
		})
		return
	}

	user, err := GetUserWithCache(req.Username)
	if err != nil {
		logrus.Infof("[api.AddUser] GetUserWithCache db error, username: %v, err: %v", req.Username, err.Error())
		ctx.JSON(400, gin.H{
			"errMsg": "db error",
		})
		return
	}
	if user != nil {
		logrus.Infof("[api.AddUser] username exists, username: %v, err: %v", req.Username, err.Error())
		ctx.JSON(400, gin.H{
			"errMsg": "username exists",
		})
		return
	}

	err = models.AddUser(req.Username, req.Password, req.Kind)
	if err != nil {
		logrus.Errorf("[api.AddUser] models.AddUser db error, username: %v, err: %v", req.Username, err)
		ctx.JSON(400, gin.H{
			"errMsg": "db error",
		})
		return
	}

	ctx.JSON(201, gin.H{
		"errMsg": "",
	})
}

func ExistsUser(ctx *gin.Context) {
	username := ctx.Query("username")

	user, err := GetUserWithCache(username)
	if err != nil {
		logrus.Infof("[api.ExistsUser] GetUserWithCache db error, username: %v, err: %v", username, err)
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
