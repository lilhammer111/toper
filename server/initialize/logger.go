package initialize

import "go.uber.org/zap"

func Logger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return
	}
	zap.ReplaceGlobals(logger)
}
