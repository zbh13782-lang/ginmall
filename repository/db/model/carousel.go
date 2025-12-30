package model

import "github.com/jinzhu/gorm"
//轮播图
type Carousel struct {
	gorm.Model
	ImgPath   string
	ProductID uint `gorm:"not null"`
}