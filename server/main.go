package main

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"to-persist/server/initialize"
)

var router *gin.Engine

func init() {
	initialize.Logger()
	initialize.Config()
	initialize.DB()
	router = initialize.Routers()

}

func main() {

	go func() {
		err := router.Run(":8520")
		if err != nil {
			zap.S().Panic("failed to run server, because ", err.Error())
		}
	}()

	// receive quit signal
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.S().Info("Bye ~")
}
