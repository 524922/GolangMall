package service

import (
	"cmall/model"
	"cmall/pkg/e"
	"cmall/pkg/logging"
	"cmall/serializer"
)

type ListUserService struct {
	UserName string `form:"user_name" json:"user_name"`
}

type ListUsersService struct {
	Limit int `form:"limit" json:"limit"`
	Start int `form:"start" json:"start"`
	Type  int `form:"type" json:"type"`
}

func (service *ListUserService) List() serializer.Response {
	var Users model.User
	code := e.SUCCESS
	err := model.DB.First(&Users).Error
	if err != nil {
		logging.Info(err)
		code = e.ERROR_ADMIN_FIND_USER
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.Response{
		Status: code,
		Data:   serializer.BuildUser(Users),
		Msg:    e.GetMsg(code),
	}
}

func (service *ListUsersService) List() serializer.Response {
	users := []model.User{}
	total := 0
	code := e.SUCCESS
	if service.Limit == 0 {
		service.Limit = 15
	}
	if service.Type == 0 {
		if err := model.DB.Model(model.User{}).
			Count(&total).Error; err != nil {
			logging.Info(err)
			code = e.ERROR_DATABASE
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}
		if err := model.DB.Limit(service.Limit).
			Offset(service.Start).Find(&users).
			Error; err != nil {
			logging.Info(err)
			code = e.ERROR_DATABASE
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}
	} else {
		if err := model.DB.Model(model.User{}).
			Where("type=?", service.Type).
			Count(&total).Error; err != nil {
			logging.Info(err)
			code = e.ERROR_DATABASE
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}
		if err := model.DB.Where("type=?", service.Type).
			Limit(service.Limit).
			Offset(service.Start).Find(&users).
			Error; err != nil {
			logging.Info(err)
			code = e.ERROR_DATABASE
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}
	}
	return serializer.BuildListResponse(serializer.BuildUsers(users), uint(total))

}
