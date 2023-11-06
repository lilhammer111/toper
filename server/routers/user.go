package routers

import (
	"github.com/gin-gonic/gin"
	"to-persist/server/api"
)

func InitUserRoutes(APIGroup *gin.RouterGroup) {
	UserGroup := APIGroup.Group("user")
	{
		UserGroup.GET("list", api.GetUserList)
		UserGroup.POST("login", api.Login)
		UserGroup.POST("register", api.Register)
	}
}
