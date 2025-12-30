package service

import (
	"context"
	"gin_mall/pkg/e"
	"gin_mall/repository/db/dao"
	"gin_mall/repository/db/model"
	"gin_mall/serializer"

	"github.com/sirupsen/logrus"
)

type FavoriteService struct {
	ProductId  uint `form:"product_id" json:"product_id"`
	BossId     uint `form:"boss_id" json:"boss_id"`
	FavoriteId uint `form:"favorite_id" json:"favorite_id"`
	PageNum    int  `form:"pageNum"`
	PageSize   int  `form:"pageSize"`
}

func (service *FavoriteService) Show(ctx context.Context, uid uint) serializer.Response {
	var (
		code    = e.SUCCESS
		fav_dao = dao.NewFavoritesDao(ctx)
		err     error
	)
	if service.PageSize == 0 {
		service.PageSize = 15
	}
	favors, total, err := fav_dao.ListFavoriteByUserId(uid, service.PageSize, service.PageNum)
	if err != nil {
		logrus.Info(err)
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.BuildListResponse(serializer.BuildFavorites(ctx, favors), uint(total))

}

func (service *FavoriteService) Create(ctx context.Context, uid uint) serializer.Response {
	var (
		code    = e.SUCCESS
		fav_dao = dao.NewFavoritesDao(ctx)
		userdao = dao.NewUserDao(ctx)
	)
	exist, _ := fav_dao.FavoriteExistOrNot(service.FavoriteId, uid)
	if exist {
		code = e.ErrorExistFavorite
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	user, err := userdao.GetUserbyid(uid)
	if err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	bossdao := dao.NewUserDaobyDB(userdao.DB)
	boss, err := bossdao.GetUserbyid(service.BossId)
	if err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	productdao := dao.NewProductDao(ctx)
	product, err := productdao.GetProductById(service.ProductId)
	if err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	favorite := &model.Favorite{
		UserID:    uid,
		User:      *user,
		ProductID: service.ProductId,
		Product:   *product,
		BossID:    service.BossId,
		Boss:      *boss,
	}
	err = fav_dao.CreateFavorite(favorite)
	if err != nil {
		code = e.ErrorDatabase
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

func (service *FavoriteService) Delete(ctx context.Context) serializer.Response {
	var (
		fav_dao = dao.NewFavoritesDao(ctx)
		code    = e.SUCCESS
		err     error
	)

	err = fav_dao.DeleteFavoriteById(service.FavoriteId)
	if err != nil {
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
	}
}
