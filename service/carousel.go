package service

import (
	"context"
	"gin_mall/pkg/e"
	"gin_mall/repository/db/dao"
	"gin_mall/serializer"

	"github.com/sirupsen/logrus"
)

type ListCarouselsService struct {
}

func (service *ListCarouselsService) List() serializer.Response {
	code := e.SUCCESS
	carouselsCtx := dao.NewCarouselDao(context.Background())
	carousels, err := carouselsCtx.ListAddress()
	if err != nil {
		logrus.Info(err)
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.BuildListResponse(serializer.BuildCarousels(carousels), uint(len(carousels)))
}
