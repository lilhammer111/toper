package initialize

import (
	"github.com/gin-gonic/gin"
	"to-persist/server/api"
	"to-persist/server/routers"
)

func Routers() *gin.Engine {

	r := gin.Default()

	ApiGroup := r.Group("/v1")

	ApiGroup.GET("/ping", api.Pong)

	ApiGroup.GET("/sms", api.SendSms)

	routers.InitUserRoutes(ApiGroup)

	return r
}
