package model

import (
	"gorm.io/gorm"
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
	gorm.Model
	ToperID uint   `gorm:"not null"`
	UserID  uint   `gorm:"not null"`
	Acronym string `gorm:"type:varchar(10)"`
	Desc    string
	Done    bool
}
