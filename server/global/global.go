package global

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"os"
	"strconv"
	"to-persist/server/config"
)

var (
	DB *gorm.DB

	Config = &config.ServerConfig{}

	Debugging       bool
	AccessKeyId     string
	AccessKeySecret string
)

func init() {
	var err error
	Debugging, err = strconv.ParseBool(os.Getenv("TOPER_DEBUG"))
	if err != nil {
		zap.S().Panic("failed to convert string to bool type, because ", err.Error())
	}

	AccessKeyId = os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID")

	AccessKeySecret = os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET")
}
