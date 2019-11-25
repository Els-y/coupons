package routers

import (
	"github.com/Els-y/coupons/server/middlewares"
	"github.com/Els-y/coupons/server/routers/api"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(middlewares.LoggerToFile())
	r.Use(gin.Recovery())

	apis := r.Group("/api")
	apis.POST("/users", api.AddUser)
	apis.GET("/users", api.ExistsUser)
	apis.POST("/auth", api.Auth)
	apis.POST("/users/:username/coupons", api.AddCoupons)
	apis.GET("/users/:username/coupons", api.GetCouponsInfo)
	apis.PATCH("/users/:username/coupons/:name", api.AssignCoupon)

	return r
}
