package service

import (
	"context"
	"gin_mall/pkg/e"
	"gin_mall/repository/db/dao"
	"gin_mall/repository/db/model"
	"gin_mall/serializer"
	"strconv"

	"github.com/sirupsen/logrus"
)

type CartService struct {
	Id        uint `form:"id" json:"id"`
	BossID    uint `form:"boss_id" json:"boss_id"`
	ProductId uint `form:"product_id" json:"product_id"`
	Num       uint `form:"num" json:"num"`
}

func (service *CartService) Create(ctx context.Context, uid uint) serializer.Response {
	var (
		code = e.SUCCESS
		err  error
	)
	cartdao := dao.NewCartDao(ctx)
	product := &model.Product{}
	productdao := dao.NewProductDao(ctx)

	product, err = productdao.GetProductById(service.ProductId)
	if err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	cart, status, _ := cartdao.CreateCart(product.ID, uid, service.BossID)
	//cart, status, _ := cartDao.CreateCart(service.ProductId, uId, service.BossID)
	if status == e.ErrorProductMoreCart {
		return serializer.Response{
			Status: status,
			Msg:    e.GetMsg(status),
		}
	}
	userDao := dao.NewUserDao(ctx)
	boss, _ := userDao.GetUserbyid(service.BossID)
	return serializer.Response{
		Status: status,
		Msg:    e.GetMsg(status),
		Data:   serializer.BuildCart(cart, product, boss),
	}
}

func (service *CartService) Show(ctx context.Context,uid uint)serializer.Response{
	var (
		code = e.SUCCESS
		err error
	)
	cartdao:=dao.NewCartDao(ctx)
	//userdao:=dao.NewUserDao(ctx)
	carts,err:=cartdao.ListCartByUserId(uid)
	if err!=nil{
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg: e.GetMsg(code),
		}
	}

	return serializer.Response{
		Status: code,
		Msg: e.GetMsg(code),
		Data: serializer.BuildCarts(carts),
	}
}

func (service *CartService) Update(ctx context.Context,cid string)serializer.Response{
	var (
		code = e.SUCCESS
		err error
	)
	cartid,_:=strconv.Atoi(cid)
	cartdao:=dao.NewCartDao(ctx)
	err = cartdao.UpdateCartNumById(uint(cartid),service.Num)
	if err!=nil{
		logrus.Info(err)
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg: e.GetMsg(code),
			Error: err.Error(),
		}
	}

	return serializer.Response{
		Status: code,
		Msg: e.GetMsg(code),
	}
}

func (service *CartService)Delete(ctx context.Context) serializer.Response{
	var(
		code = e.SUCCESS
		//cartid ,_ = strconv.Atoi(cid)
		err error
	)
	cartdao:=dao.NewCartDao(ctx)
	err = cartdao.DeleteCartById(service.Id)
	if err!=nil{
		logrus.Info(err)
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg: e.GetMsg(code),
			Error: err.Error(),
		}
	}
	return serializer.Response{
		Status: code,
		Msg: e.GetMsg(code),
	}
}