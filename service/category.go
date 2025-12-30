package service

import (
	"context"
	"gin_mall/pkg/e"
	"gin_mall/repository/db/dao"
	"gin_mall/serializer"

	"github.com/sirupsen/logrus"
)

type ListCategoryService struct {
}

func (service *ListCategoryService) List(ctx context.Context) serializer.Response {
	code := e.SUCCESS
	categoryDao := dao.NewCategoryDao(ctx)
	categories, err := categoryDao.ListCategory()
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
		Msg:    e.GetMsg(code),
		Data:   serializer.BuildCategories(categories),
	}

}
