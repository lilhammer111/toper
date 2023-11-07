package initialize

import (
	"go.uber.org/zap"
	"to-persist/server/util/scheduler"
)

func Scheduler() {
	taskScheduler := scheduler.NewTaskScheduler()
	taskScheduler.Start()

	// 从数据库加载并初始化任务
	if err := taskScheduler.ReInitTasksFromDB(); err != nil {
		zap.S().Panicf("failed to initialize tasks: %v", err)
	}

}
