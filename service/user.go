package service

import (
	"bookstore/dao"
	"bookstore/model"
	"bookstore/pkg/e"
	"bookstore/pkg/utils"
	"bookstore/serializer"
	"context"
)

type UserService struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
	Key      string `json:"key" form:"key"` // 前端验证
}

func (service UserService) Register(ctx context.Context) serializer.Response {
	var user model.User
	code := e.Success
	// 如果加敏的密钥不符合这个条件的话就报错
	if service.Key == "" || len(service.Key) != 16 {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  "密钥长度不足",
		}
	}

	// 10000 --> 密文存储 对称加密操作
	utils.Encrypt.SetKey(service.Key)
	userDao := dao.NewUserDao(ctx)
	_, exist, err := userDao.ExistOrNotByUserName(service.Username)
	if err != nil {
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	if exist {
		code = e.ErrorExistsUser
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	user = model.User{
		Username: service.Username,
		Avatar:   "avatar.JPG",
		Status:   model.Active,
		Monty:    utils.Encrypt.AesEncoding("10000"), // 初始金额的加密
	}
	// 密码加密
	if err = user.SetPassword(service.Password); err != nil {
		code = e.ErrorFailEncryption
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	// 创建用户
	err = userDao.CreateUser(&user)
	if err != nil {
		code = e.Error
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}
