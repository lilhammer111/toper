package initialize

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func Logger() {
	// 创建一个EncoderConfig用于定义日志的格式
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zap.InfoLevel)

	// 创建并设置日志文件
	file, err := os.OpenFile("server/zap.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	writerSyncer := zapcore.AddSync(file)

	// 创建日志的Core部分
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig), // 或 NewConsoleEncoder
		writerSyncer,
		atomicLevel,
	)

	// 初始化Logger
	logger := zap.New(core)
	zap.ReplaceGlobals(logger)
}
