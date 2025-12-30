package service

import (
	"context"
	"fmt"
	"gin_mall/pkg/e"
	"gin_mall/repository/cache"
	"gin_mall/repository/db/dao"
	"gin_mall/repository/db/model"
	"gin_mall/serializer"
	"math/rand"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

const OrderTimeKey = "OrderTime"

type OrderService struct {
	ProductID uint `form:"product_id" json:"product_id"`
	Num       uint `form:"num" json:"num"`
	AddressID uint `form:"address_id" json:"address_id"`
	Money     int  `form:"money" json:"money"`
	BossID    uint `form:"boss_id" json:"boss_id"`
	UserID    uint `form:"user_id" json:"user_id"`
	OrderNum  uint `form:"order_num" json:"order_num"`
	Type      int  `form:"type" json:"type"`
	model.BasePage
}

func (service *OrderService) Create(ctx context.Context, id uint) serializer.Response {
	code := e.SUCCESS
	order := &model.Order{
		UserID:    id,
		ProductID: service.ProductID,
		BossID:    service.BossID,
		Num:       int(service.Num),
		Money:     float64(service.Money),
		Type:      1,
	}
	addressdao := dao.NewAddressDao(ctx)
	address, err := addressdao.GetAddressByAid(service.AddressID)
	if err != nil {
		logrus.Info(err)
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	order.AddressID = address.ID
	number := fmt.Sprintf("%09v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000000))
	productnum := strconv.Itoa(int(service.ProductID))
	usernum := strconv.Itoa(int(id))
	number = number + productnum + usernum
	ordernum, _ := strconv.ParseUint(number, 10, 64)
	order.OrderNum = ordernum
	orderdao := dao.NewOrderDao(ctx)
	err = orderdao.CreateOrder(order)
	if err != nil {
		logrus.Info(err)
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	data := redis.Z{
		Score:  float64(time.Now().Unix()) + 15*time.Minute.Seconds(),
		Member: order.Num,
	}

	cache.RedisClient.ZAdd(OrderTimeKey, data)
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}

func (service *OrderService) List(ctx context.Context, uid uint) serializer.Response {
	var (
		orders []*model.Order
		total  int64
		code   = e.SUCCESS
	)
	if service.PageSize == 0 {
		service.PageSize = 15
	}

	orderdao := dao.NewOrderDao(ctx)
	condition := make(map[string]interface{})
	condition["user_id"] = uid
	if service.Type == 0 {
		condition["type"] = 0
	} else {
		condition["type"] = service.Type
	}
	orders, total, err := orderdao.ListOrderByCondition(condition, service.BasePage)
	if err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	return serializer.BuildListResponse(serializer.BuildOrders(ctx, orders), uint(total))
}

func (service *OrderService) Show(ctx context.Context, oid string) serializer.Response {
	var (
		code       = e.SUCCESS
		orderid, _ = strconv.Atoi(oid)
		orderdao   = dao.NewOrderDao(ctx)
		order, _   = orderdao.GetOrderById(uint(orderid))
	)

	addressdao := dao.NewAddressDao(ctx)
	address, err := addressdao.GetAddressByAid(order.AddressID)
	if err != nil {
		logrus.Info(err)
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}

	productdao := dao.NewProductDao(ctx)
	product, err := productdao.GetProductById(order.ProductID)
	if err != nil {
		logrus.Info(err)
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
		Data:   serializer.BuildOrder(order, product, address),
	}
}

func (service *OrderService) Delete(ctx context.Context, oid string) serializer.Response {
	code := e.SUCCESS

	orderDao := dao.NewOrderDao(ctx)
	orderId, _ := strconv.Atoi(oid)
	err := orderDao.DeleteOrderById(uint(orderId))
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
	}
}
