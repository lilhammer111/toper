package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Pong(c *gin.Context) {
	c.Status(http.StatusOK)
}
