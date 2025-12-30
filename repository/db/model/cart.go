package model

import "github.com/jinzhu/gorm"

type Cart struct {
	gorm.Model
	UserID    uint
	//Product   Product `gorm:"ForeignKey:ProductID"`
	ProductID uint    `gorm:"not null"`
	BossID    uint
	Num       uint
	MaxNum    uint
	Check     bool
}
