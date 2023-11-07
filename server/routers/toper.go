package routers

import (
	"github.com/gin-gonic/gin"
	"to-persist/server/api"
	"to-persist/server/middlewares"
)

func InitToperRoutes(APIGroup *gin.RouterGroup) {
	UserGroup := APIGroup.Group("toper").Use(middlewares.JwtAuth())
	{
		UserGroup.GET("", api.List)
		UserGroup.POST("", api.Create)
		UserGroup.POST("/status", api.Done)
		UserGroup.GET("/history", api.History)
		UserGroup.PUT("/:acronym", api.Alter)
	}
}
