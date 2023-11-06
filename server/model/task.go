package model

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	ToperID      string `gorm:"varchar(20);index"`
	Expression   string `gorm:"varchar(50)"`
	TaskFuncType string `gorm:"varchar(30)"`
}
