package main

import (
	"fmt"
	"github.com/Els-y/coupons/server/models"
	"github.com/Els-y/coupons/server/pkgs/redis"
	"github.com/Els-y/coupons/server/pkgs/setting"
	"github.com/Els-y/coupons/server/routers"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func init() {
	setting.Setup()
	models.Setup()
	redis.Setup()
}

func main() {
	gin.SetMode(setting.ServerSetting.RunMode)
	if setting.ServerSetting.RunMode == "release" {
		logrus.SetLevel(logrus.WarnLevel)
	}

	routersInit := routers.InitRouter()
	readTimeout := setting.ServerSetting.ReadTimeout
	writeTimeout := setting.ServerSetting.WriteTimeout
	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)
	maxHeaderBytes := 1 << 20

	server := &http.Server{
		Addr:           endPoint,
		Handler:        routersInit,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	logrus.Infof("start http server listening %s", endPoint)

	server.ListenAndServe()
}
