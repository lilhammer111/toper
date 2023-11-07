package model

import (
	"gorm.io/gorm"
	"time"
)

type Toper struct {
	gorm.Model
	Description string `gorm:"not null;type:varchar(100);" json:"description"`
	Acronym     string `gorm:"type:varchar(10);index:idx_user_acronym,unique" json:"acronym"`
	DueDate     string `gorm:"not null" json:"due-date"`
	Period      string `gorm:"not null" json:"period"`

	UserID uint `gorm:"index:idx_user_acronym,unique" json:"user-id"`
	User   User
}

type DoneHistory struct {
	ID       int       `gorm:"primaryKey;type:int" json:"-"`
	DoneTime time.Time `gorm:"autoCreateTime;not null" json:"done-time"`
	ToperID  uint      `gorm:"not null" json:"toper-id"`
	Acronym  string    `gorm:"type:varchar(10);not null" json:"acronym"`
	Done     string    `gorm:"not null" json:"done"`
}
