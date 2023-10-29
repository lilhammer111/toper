package initialize

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
	"to-persist/server/global"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func DB() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		global.Config.MysqlConfig.Username,
		global.Config.MysqlConfig.Password,
		global.Config.MysqlConfig.AddrConfig.Host,
		global.Config.MysqlConfig.AddrConfig.Port,
		global.Config.MysqlConfig.DBName,
	)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  true,        // Enable color
		},
	)

	var err error
	global.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
		Logger:         newLogger,
	})

	if err != nil {
		zap.S().Panic("failed to init db, because ", err.Error())
	}
}