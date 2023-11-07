package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
	"to-persist/server/model"
)

func main() {
	dsn := "root:root@tcp(127.0.0.1:3306)/toper?charset=utf8mb4&parseTime=True&loc=Local"

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  true,        // Enable color
		},
	)

	// Globally mode
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
		Logger:         newLogger,
	})
	if err != nil {
		log.Fatal(err)
	}
	tableExist = false
	if tableExist {

		hashedPassword, e := bcrypt.GenerateFromPassword([]byte("12345678"), bcrypt.DefaultCost)
		if e != nil {
			log.Fatal(e)
		}

		newPassword := string(hashedPassword)
		for i := 0; i < 10; i++ {
			user := model.User{
				Mobile:   fmt.Sprintf("1953587698%d", i),
				Password: newPassword,
				Username: fmt.Sprintf("lilhammer11%d", i),
			}
			db.Save(&user)
		}
	} else {
		//err = db.AutoMigrate(&model.User{})
		err = db.AutoMigrate(&model.DoneHistory{})
		if err != nil {
			log.Fatal(err)
		}
	}
}

var tableExist bool
