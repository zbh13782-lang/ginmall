package service

import (
	"context"
	"gin_mall/conf"
	"gin_mall/consts"
	"gin_mall/pkg/e"
	"gin_mall/pkg/utils"
	"gin_mall/repository/db/dao"
	"gin_mall/repository/db/model"
	"gin_mall/serializer"
	"mime/multipart"
	"strings"
	"time"

	"github.com/go-mail/mail"
	"github.com/sirupsen/logrus"
)

type UserLoginService struct {
	UserName string `form:"user_name" json:"user_name" binding:"required,min=5,max=15"`
	Password string `form:"password" json:"password" binding:"required,min=8,max=16"`
}

type UserRegisterService struct {
	Nickname string `form:"nickname" json:"nickname" binding:"required,min=2,max=10"`
	UserName string `form:"user_name" json:"user_name" binding:"required,min=5,max=15"`
	Password string `form:"password" json:"password" binding:"required,min=8,max=16"`
	Key      string `form:"key" json:"key" binding:"required,len=16"`
}

type UserUpdateService struct {
	Nickname string `form:"nickname" json:"nickname"`
}

type UserUploadService struct {
}
type SendEmailService struct {
	Email    string `form:"email" json:"email"`
	Password string `form:"password" json:"password"`
	//1.绑定邮箱，2.解绑邮箱,3.改密码
	OperationType uint `form:"operation_type" json:"operation_type"`
}
type ValidEmailService struct {
}

// 注册
func (service *UserRegisterService) Register(ctx context.Context) serializer.Response {
	var user *model.User
	var code = e.SUCCESS
	userdao := dao.NewUserDao(ctx)

	if service.Key == "" || len(service.Key) != 16 {
		code = e.ERROR
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Data:   "密钥长度不足",
		}
	}

	utils.Encrypt.SetKey(service.Key)
	_, exist, err := userdao.ExistOrNotbyName(service.UserName)
	if err != nil {
		code = e.ErrorDatabase
		logrus.Info(err)
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}

	}
	if exist {
		code = e.ErrorExistUser
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	money := "10000"
	user = &model.User{
		NickName: service.Nickname,
		UserName: service.UserName,
		Status:   model.Active,
		Money:    money,
	}
	if err = user.SetPassword(service.Password); err != nil {
		code = e.ErrorFailEncryption
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}

	user.Avatar = "avatar.JPG"

	if err = userdao.CreateUser(user); err != nil {
		code = e.ErrorDatabase
		logrus.Info(err)
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}

	// 生成JWT token
	token, err := utils.GenerateToken(user.ID, user.UserName, 0)
	if err != nil {
		code = e.ErrorAuthToken
		logrus.Info(err)
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}

	return serializer.Response{
		Status: code,
		Data: serializer.TokenData{
			Token: token,
			User:  serializer.BuildUser(user),
		},
		Msg: e.GetMsg(code),
	}

}

// 登录
func (service *UserLoginService) Login(ctx context.Context) serializer.Response {
	var userdao = dao.NewUserDao(ctx)
	var code = e.SUCCESS
	var user *model.User
	user, exist, err := userdao.ExistOrNotbyName(service.UserName)
	if !exist {
		//	logging.Infoln(err)
		logrus.Info(err)
		code = e.ErrorNotExistUser
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	if user.CheckPassword(service.Password) == false {
		code = e.ErrorNotComparePassword
		logrus.Info(err)
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	token, err := utils.GenerateToken(user.ID, user.UserName, 0) //第二个用service or user
	if err != nil {
		code = e.ErrorAuthToken
		//logging.Infoln(err)
		logrus.Info(err)
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	return serializer.Response{
		Status: code,
		Data: serializer.TokenData{
			Token: token,
			User:  serializer.BuildUser(user),
		},
		Msg: e.GetMsg(code),
	}

}

// 修改
func (service *UserUpdateService) Update(ctx context.Context, uid uint) serializer.Response {
	var user *model.User
	var err error
	code := e.SUCCESS
	userdao := dao.NewUserDao(ctx)
	user, err = userdao.GetUserbyid(uid)
	if err != nil {
		code = e.ErrorDatabase
		//logging.Infoln(err)
		logrus.Info(err)
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	if service.Nickname != "" {
		user.NickName = service.Nickname
	}
	err = userdao.UpdateUserbyid(uid, user)
	if err != nil {
		code = e.ErrorDatabase
		//logging.Infoln(err)
		logrus.Info(err)
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.Response{
		Status: code,
		Data:   serializer.BuildUser(user),
		Msg:    e.GetMsg(code),
	}
}

// 绑定邮件
func (service *SendEmailService) SendEmail(ctx context.Context, uid uint) serializer.Response {
	code := e.SUCCESS
	var address string
	var notice *model.Notice

	token, err := utils.GenerateEmailToken(uid, service.OperationType, service.Email, service.Password)
	if err != nil {
		code = e.ErrorDatabase
		//logging.Infoln(err)
		logrus.Info(err)
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	noticedao := dao.NewNoticeDao(ctx)
	notice, err = noticedao.GetNoticeById(service.OperationType)
	if err != nil {
		code = e.ErrorDatabase
		logrus.Info(err)
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	address = conf.ValidEmail + token
	mailStr := notice.Text
	mailText := strings.Replace(mailStr, "Email", address, -1)
	m := mail.NewMessage()
	m.SetHeader("From", conf.SmtpEmail)
	m.SetHeader("To", service.Email)
	m.SetHeader("Subject", "zbh")
	m.SetBody("text/html", mailText)
	d := mail.NewDialer(conf.SmtpHost, 465, conf.SmtpEmail, conf.SmtpPass)
	d.StartTLSPolicy = mail.MandatoryStartTLS
	if err := d.DialAndSend(m); err != nil {
		logrus.Info(err)
		code = e.ErrorSendEmail
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

// 验证Valid
func (service *ValidEmailService) Valid(ctx context.Context, token string) serializer.Response {
	var userid uint
	var email string
	var password string
	var operationtype uint
	code := e.SUCCESS
	if token == "" {
		code = e.InvalidParams

	} else {
		claims, err := utils.ParseEmailToken(token)
		if err != nil {
			logrus.Info(err)
			code = e.ErrorAuthCheckTokenFail
		} else if time.Now().Unix() > claims.ExpiresAt {
			code = e.ErrorAuthCheckTokenTimeout
		} else {
			userid = claims.UserID
			email = claims.Email
			password = claims.Password
			operationtype = claims.OperationType
		}
	}
	if code != e.SUCCESS {
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	userdao := dao.NewUserDao(ctx)
	user, err := userdao.GetUserbyid(userid)
	if err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}

	}
	if operationtype == 1 {
		user.Email = email

	} else if operationtype == 2 {
		user.Email = ""

	} else if operationtype == 3 {
		err = user.SetPassword(password)
		if err != nil {
			code = e.ErrorDatabase
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
			}
		}
	}
	if err = userdao.UpdateUserbyid(userid, user); err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
		Data:   serializer.BuildUser(user),
	}
}

// 上传头像
func (service *UserUploadService) Post(ctx context.Context, uId uint, file multipart.File, fileSize int64) serializer.Response {
	code := e.SUCCESS
	var user *model.User
	var err error

	userDao := dao.NewUserDao(ctx)
	user, err = userDao.GetUserbyid(uId)
	if err != nil {
		logrus.Info(err)
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	var path string
	if conf.UploadModel == consts.UploadModelLocal { // 兼容两种存储方式
		utils.LogrusObj.Infof("使用本地存储模式上传头像，用户ID: %d", uId)
		path, err = utils.UploadAvatarToLocalStatic(file, uId, user.UserName)
	} else {
		utils.LogrusObj.Infof("使用七牛云存储模式上传头像，用户ID: %d, 文件大小: %d", uId, fileSize)
		path, err = utils.UploadToQiNiu(file, fileSize)
	}
	if err != nil {
		utils.LogrusObj.Errorf("头像上传失败: %v", err)
		code = e.ErrorUploadFile
		return serializer.Response{
			Status: code,
			Data:   e.GetMsg(code),
			Error:  err.Error(), // 返回真正的错误信息
		}
	}

	user.Avatar = path
	err = userDao.UpdateUserbyid(uId, user)
	if err != nil {
		logrus.Info(err)
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.Response{
		Status: code,
		Data:   serializer.BuildUser(user),
		Msg:    e.GetMsg(code),
	}
}
