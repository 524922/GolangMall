package service

import (
	"cmall/model"
	"cmall/pkg/e"
	"cmall/pkg/logging"
	"cmall/pkg/util"
	"cmall/serializer"
	"gopkg.in/mail.v2"
	"os"
	"strings"
	"time"
)

type SendEmailService struct {
	UserID        uint   `form:"user_id" json:"user_id"`
	Email         string `form:"email" json:"email"`
	Password      string `form:"password" json:"password"`
	OperationType uint   `form:"operation_type" json:"operation_type"` //operationType 1.绑定邮箱 2.解绑邮箱 3.改密码
}

//send 发送邮箱
func (service *SendEmailService) Send() serializer.Response {
	code := e.SUCCESS
	var address string
	var notice model.Notice
	token, err := util.GenerateEmailToken(service.UserID, service.OperationType, service.Email, service.Password)
	if err != nil {
		logging.Info(err)
		code = e.ERROR_AUTH_TOKEN
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	//数据库里面对应邮件id=operation_type+1
	if err := model.DB.First(&notice, service.OperationType+1).Error; err != nil {
		logging.Info(err)
		code = e.ERROR_DATABASE
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	address = os.Getenv("VAILD_EMAIL") + token
	mailStr := notice.Text
	mailText := strings.Replace(mailStr, "VaildAddress", address, -1)
	m := mail.NewMessage()
	m.SetHeader("Form", os.Getenv("SMTP_EMAIL"))
	m.SetHeader("To", service.Email)
	m.SetHeader("Subject", "CMall")
	m.SetBody("text/html", mailText)
	d := mail.NewDialer(os.Getenv("SMTP_HOST"), 465, os.Getenv("SMTP_EMAIL"), os.Getenv("SMTP_PASS"))
	d.StartTLSPolicy = mail.MandatoryStartTLS

	//发邮件
	if err := d.DialAndSend(m); err != nil {
		logging.Info(err)
		code = e.ERROR_SEND_EMAIL
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}

//VaildEmailService 绑定、解绑邮箱和修改密码的服务
type VaildEmailService struct {
	Token string `form:"token" json:"token"`
}

//Vaild 绑定邮箱
func (service *VaildEmailService) Vaild() serializer.Response {
	var userID uint
	var email string
	var password string
	var operationType uint
	code := e.SUCCESS
	//验证token
	if service.Token == "" {
		code = e.INVALID_PARAMS
	} else {
		claims, err := util.ParseEmailToken(service.Token)
		if err != nil {
			logging.Info(err)
			code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
		} else if time.Now().Unix() > claims.ExpiresAt {
			code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
		} else {
			userID = claims.UserID
			email = claims.Email
			password = claims.Password
			operationType = claims.OperationType
		}
	}
	if code != e.SUCCESS {
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	if operationType == 1 {
		//1.绑定邮箱
		if err := model.DB.Table("user").Where("id=?", userID).Update("email", email).Error; err != nil {
			logging.Info(err)
			code = e.ERROR_DATABASE
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}
	} else if operationType == 2 {
		//2.解绑邮箱
		if err := model.DB.Table("user").Where("id=?", userID).Update("email", "").Error; err != nil {
			logging.Info(err)
			code = e.ERROR_DATABASE
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}
	}
	//获取该用户信息
	var user model.User
	if err := model.DB.First(&user, userID).Error; err != nil {
		logging.Info(err)
		code = e.ERROR_DATABASE
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	//3.修改密码
	if operationType == 3 {
		//加密密码
		if err := user.SetPassword(password); err != nil {
			logging.Info(err)
			code = e.ERROR_FAIL_ENCRYPTION
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
			}
		}
		if err := model.DB.Save(&user).Error; err != nil {
			logging.Info(err)
			code = e.ERROR_DATABASE
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
			}
		}
		code = e.UPDATE_PASSWORD_SUCCESS
		//返回修改密码成功信息
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	//返回用户信息
	return serializer.Response{
		Status: code,
		Data:   serializer.BuildUser(user),
		Msg:    e.GetMsg(code),
	}
}
