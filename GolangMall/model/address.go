package model

import "github.com/jinzhu/gorm"

type Address struct {
	gorm.Model
	UserID  uint
	Name    string
	Phone   string
	Address string
}
