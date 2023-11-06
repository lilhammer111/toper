package main

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"to-persist/server/global"
	"to-persist/server/initialize"
	"to-persist/server/util/scheduler"
)

var router *gin.Engine

func init() {
	initialize.Logger()
	initialize.Config()
	initialize.MysqlDB()
	initialize.RedisClient()
	initialize.Scheduler()
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
	closeAll()
	zap.S().Info("Bye #############################################################################")
}

func closeAll() {
	err := global.RedisClient.Close()
	if err != nil {
		zap.S().Error("failed to close redis conn")
	}

	taskScheduler := scheduler.NewTaskScheduler()
	taskScheduler.Stop()

}
