package model

import "github.com/jinzhu/gorm"

type Favorite struct {
	gorm.Model
	UserID    uint
	ProductID uint
	BossID    uint
}
