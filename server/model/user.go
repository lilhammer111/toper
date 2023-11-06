package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(20);unique"`
	Mobile   string `gorm:"index:idx_mobile;unique;type:varchar(11);not null"`
	Password string `gorm:"type:varchar(100);not null;"`
}
