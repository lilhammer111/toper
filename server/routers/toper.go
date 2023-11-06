package routers

import (
	"github.com/gin-gonic/gin"
	"to-persist/server/api"
	"to-persist/server/middlewares"
)

func InitToperRoutes(APIGroup *gin.RouterGroup) {
	UserGroup := APIGroup.Group("toper").Use(middlewares.JwtAuth())
	{
		UserGroup.POST("", api.Create)
		UserGroup.GET("", api.List)
		UserGroup.GET("/:id", api.History)
		UserGroup.PUT("/:id", api.Alter)
	}
}
