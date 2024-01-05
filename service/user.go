package service

import (
	"bookstore/conf"
	"bookstore/dao"
	"bookstore/model"
	"bookstore/pkg/e"
	"bookstore/pkg/utils"
	"bookstore/serializer"
	"context"
	"gopkg.in/mail.v2"
	"mime/multipart"
	"strings"
)

type UserService struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
	Key      string `json:"key" form:"key"` // 前端验证
}

type SendEmailService struct {
	Email         string `json:"email" form:"email"`
	Password      string `json:"password" form:"password"`
	OperationType uint   `json:"operation_type" form:"operation_type"` // 1：绑定邮箱 2：解绑邮箱 3：改密码
}

func (service *UserService) Register(ctx context.Context) serializer.Response {
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

func (service *UserService) Login(ctx context.Context) serializer.Response {
	var user *model.User
	code := e.Success
	userDao := dao.NewUserDao(ctx)
	// 判断用户是否存在
	user, exist, err := userDao.ExistOrNotByUserName(service.Username)
	if !exist || err != nil {
		code = e.ErrorExistsUserNotFound
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Data:   "用户不存在，请先注册",
		}
	}
	// 校验密码
	if user.CheckPassword(service.Password) == false {
		code = e.ErrorNotCompare
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Data:   "密码错误，请重新登录",
		}
	}
	// token 分发：http 无状态(认证 token)
	token, err := utils.GenerateToken(user.ID, service.Username, 0)
	if err != nil {
		code = e.ErrorAuthToken
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Data:   "token 认证有错误",
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
		Data: serializer.TokenData{
			User:  serializer.BuildUser(user),
			Token: token,
		},
	}
}

func (service *UserService) Update(ctx context.Context, id uint) serializer.Response {
	var user *model.User
	var err error
	code := e.Success
	// 找到这个用户
	userDao := dao.NewUserDao(ctx)
	user, err = userDao.GetUserById(id)
	// 修改用户名
	if service.Username != "" {
		user.Username = service.Username
	}
	err = userDao.UpdateUserById(id, user)
	if err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
		Data:   serializer.BuildUser(user),
	}
}

func (service *UserService) Post(ctx context.Context, id uint, file multipart.File, fileSize int64) serializer.Response {
	code := e.Success
	var user *model.User
	var err error
	userDao := dao.NewUserDao(ctx)
	user, err = userDao.GetUserById(id)
	if err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	// 保存图片到本地
	path, err := UploadAvatarToLocalStatic(file, id, user.Username)
	if err != nil {
		code = e.ErrorUploadFail
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	user.Avatar = path
	err = userDao.UpdateUserById(id, user)
	if err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
		Data:   serializer.BuildUser(user),
	}
}

func (service *SendEmailService) Send(ctx context.Context, id uint) serializer.Response {
	code := e.Success
	var address string
	var notice *model.Notice // 绑定邮箱，修改密码，模板通知
	token, err := utils.GenerateEmailToken(id, service.Email, service.Password, service.OperationType, 0)
	if err != nil {
		code = e.ErrorAuthToken
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	noticeDao := dao.NewNoticeDao(ctx)
	notice, err = noticeDao.GetNoticeById(service.OperationType)
	if err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	address = conf.ValidEmail + token // 发送方
	mailStr := notice.Text
	mailText := strings.Replace(mailStr, "邮箱", address, -1)
	m := mail.NewMessage()
	m.SetHeader("From", conf.SmtpEmail)
	m.SetHeader("To", service.Email)
	m.SetHeader("Subject", "AxLiu")
	m.SetBody("text/html", mailText)
	d := mail.NewDialer(conf.SmtpHost, 465, conf.SmtpEmail, conf.SmtpPass)
	d.StartTLSPolicy = mail.MandatoryStartTLS
	if err = d.DialAndSend(m); err != nil {
		code = e.ErrorSendEmail
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}
