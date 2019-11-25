package middlewares

import (
	"fmt"
	"github.com/Els-y/coupons/server/pkgs/file"
	"github.com/Els-y/coupons/server/pkgs/setting"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"time"
)

// getLogFilePath get the log file save path
func getLogFilePath() string {
	return path.Join(setting.AppSetting.RuntimeRootPath, setting.AppSetting.LogSavePath)
}

// getLogFileName get the save name of the log file
func getLogFileName() string {
	return fmt.Sprintf("%s%s.%s",
		setting.AppSetting.LogSaveName,
		time.Now().Format(setting.AppSetting.TimeFormat),
		setting.AppSetting.LogFileExt,
	)
}

func LoggerToFile() gin.HandlerFunc {
	logFilePath := getLogFilePath()
	logFileName := getLogFileName()

	if file.CheckNotExist(logFilePath) {
		err := file.MkDir(logFilePath)
		if err != nil {
			logrus.WithError(err).Errorf("LoggerToFile middleware create folder error")
		}
	}

	fileName := path.Join(logFilePath, logFileName)

	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("err", err)
	}

	logger := logrus.New()
	logger.Out = src
	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&logrus.TextFormatter{})

	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		endTime := time.Now()
		latencyTime := endTime.Sub(startTime)

		reqMethod := c.Request.Method
		reqUri := c.Request.RequestURI
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()

		logger.Infof("| %3d | %13v | %15s | %s | %s |",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri,
		)
	}
}
