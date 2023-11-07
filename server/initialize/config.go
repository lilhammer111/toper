package initialize

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"to-persist/server/global"
)

func Config() {

	configFilePath := fmt.Sprintf("server/config-pro.yaml")
	if global.Debugging {
		configFilePath = fmt.Sprintf("server/config-debug.yaml")
	}

	v := viper.New()

	v.SetConfigFile(configFilePath)

	err := v.ReadInConfig()
	if err != nil {
		zap.S().Panic("failed to read config file, because ", err.Error())
	}

	err = v.Unmarshal(&global.ServerConfig)
	if err != nil {
		zap.S().Panic("failed to read config file, because ", err.Error())
	}

	zap.S().Infof("server config is %+v", global.ServerConfig)

}
