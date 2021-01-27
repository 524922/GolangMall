package model

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

//Boss 店家模型
type Boss struct {
	Product []Product `gorm:"FOREIGNKEY:ProductID;ASSOCIATION_FOREIGNKEY:ID"`
	gorm.Model
	UserName       string `gorm:"unique"`
	Email          string //`gorm:"unique"`
	PasswordDigest string
	Nickname       string `gorm:"unique"`
	Status         string
	Limit          int
	Avatar         string `gorm:"size:1000"`
}

//GetBoss 用ID获取Boss
func GetBoss(ID interface{}) (Boss, error) {
	var Boss Boss
	result := DB.First(&Boss, ID)
	return Boss, result.Error
}

//SetBossPassword 设置密码
func (Boss *Boss) SetBossPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), PassWordCost)
	if err != nil {
		return err
	}
	Boss.PasswordDigest = string(bytes)
	return nil
}

//CheckPassword 校验密码
func (Boss *Boss) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(Boss.PasswordDigest), []byte(password))
	return err == nil
}

//AvatarUrl 封面地址
func (Boss *Boss) AvatarURL() string {
	signedGetURL := "http://q1.qlogo.cn/g?b=qq&nk=294350394&s=640"
	return signedGetURL
}
