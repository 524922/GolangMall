package model

import "github.com/jinzhu/gorm"

//UserAuth 用户权限模型
type UserAuth struct {
	gorm.Model
	UserID        uint
	IndentityType string //第三方应用的名称  微信、微博
	Indentifier   string `gorm:"unique"` //标识(第三方应用的唯一标识)
	Token         string //token凭证(保存 token)
	RefreshToken  string
}
