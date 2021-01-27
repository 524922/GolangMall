package service

import (
	"cmall/model"
	"cmall/pkg/e"
	"cmall/pkg/logging"
	"cmall/pkg/util"
	"cmall/serializer"
	"github.com/jinzhu/gorm"
)

type BossRegisterService struct {
	UserName        string `form:"username" json:"user_name" binding:"required,min=5,max=30"`
	PassWord        string `form:"password" json:"password" binding:"required,min=8,max=40"`
	PassWordConfirm string `form:"password_confirm" json:"password_confirm" binding:"required,min=8,max=40"`
}
type BossLoginService struct {
	BossName  string `form:"boss_name" json:"boss_name" binding:"required,min=5,max=15"`
	PassWord  string `form:"password" json:"password" binding:"required,min=8,max=15"`
	Challenge string `form:"challenge" json:"challenge"`
	Validate  string `form:"validate" json:"validate"`
	Seccode   string `form:"seccode" json:"seccode"`
}
type BossUpdateService struct {
	ID        uint   `form:"id" json:"id"`
	NickName  string `form:"nickname" json:"nickname" binding:"required,min=2,max=14"`
	UserName  string `form:"username" json:"user_name" binding:"required,min=2,max=13"`
	Avatar    string `form:"avator" json:"avatar"`
	ProductID uint   `form:"productID" json:"product_id"`
}

func (service *BossRegisterService) Valid(userID, status interface{}) *serializer.Response {
	var code int
	count := 0
	err := model.DB.Model(&model.Boss{}).Where("user_name=?", service.UserName).Count(&count).Error
	if err != nil {
		code := e.ERROR_DATABASE
		return &serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	if count > 0 {
		code = e.ERROR_EXIST_USER
		return &serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	count = 0
	err = model.DB.Model(&model.Boss{}).Where("user_name=?", service.UserName).Count(&count).Error
	if err != nil {
		code = e.ERROR_DATABASE
		return &serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	if count > 0 {
		code = e.ERROR_EXIST_USER
		return &serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	return nil
}

func (service *BossRegisterService) Register(userID, status interface{}) *serializer.Response {
	boss := model.Boss{
		UserName: service.UserName,
	}
	code := e.SUCCESS
	if res := service.Valid(userID, status); res != nil {
		return res
	}
	if err := boss.SetBossPassword(service.PassWord); err != nil {
		logging.Info(err)
		code = e.ERROR_FAIL_ENCRYPTION
		return &serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	boss.Avatar = "static\\img\\avatar\\3.jpg"
	if err := model.DB.Create(&boss).Error; err != nil {
		logging.Info(err)
		code = e.ERROR_DATABASE
		return &serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	return &serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}

func (service *BossLoginService) Login(BossID, status interface{}) serializer.Response {
	var boss model.Boss
	code := e.SUCCESS
	if err := model.DB.Where("boss_name=?", service.BossName).First(&boss).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			logging.Info(err)
			code = e.ERROR_NOT_EXIST_USER
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
			}
		}
		logging.Info(err)
		code = e.ERROR_DATABASE
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	if boss.CheckPassword(service.PassWord) == false {
		code = e.ERROR_NOT_COMPARE
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	token, err := util.GenerateToken(service.BossName, service.PassWord, 1)
	if err != nil {
		logging.Info(err)
		code = e.ERROR_AUTH_TOKEN
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	return serializer.Response{
		Status: code,
		Data:   serializer.TokenData{User: serializer.BuildBoss(boss), Token: token},
		Msg:    e.GetMsg(code),
	}
}

func (service *BossUpdateService) Update() serializer.Response {
	var boss model.Boss
	code := e.SUCCESS
	err := model.DB.First(&boss, service.ID).Error
	if err != nil {
		logging.Info(err)
		code = e.ERROR_DATABASE
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	boss.Nickname = service.NickName
	boss.UserName = service.UserName
	if service.Avatar != "" {
		boss.Avatar = service.Avatar
	}
	err = model.DB.Save(&boss).Error
	if err != nil {
		logging.Info(err)
		code = e.ERROR_DATABASE
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.Response{
		Status: code,
		Data:   serializer.BuildBoss(boss),
		Msg:    e.GetMsg(code),
	}
}
