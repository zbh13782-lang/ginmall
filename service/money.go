package service

import (
	"context"
	"gin_mall/pkg/e"
	"gin_mall/repository/db/dao"
	"gin_mall/serializer"

	"github.com/sirupsen/logrus"
)

type ShowMoneyService struct {
	Key string `json:"key" form:"key"`
}

func (service *ShowMoneyService) Show(ctx context.Context, uid uint) serializer.Response {
	code := e.SUCCESS
	userdao := dao.NewUserDao(ctx)

	user, err := userdao.GetUserbyid(uid)

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
		Data:   serializer.BuildMoney(user, service.Key),
		Msg:    e.GetMsg(code),
	}

}
