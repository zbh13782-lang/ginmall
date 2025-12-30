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

type AddressService struct {
	Name    string `form:"name" json:"name"`
	Phone   string `form:"phone" json:"phone"`
	Address string `form:"address" json:"address"`
}

func (service *AddressService) Create(ctx context.Context, uid uint) serializer.Response {
	var (
		code      = e.SUCCESS
		err       error
		adres_dao = dao.NewAddressDao(ctx)
	)
	address := &model.Address{
		UserID:  uid,
		Name:    service.Name,
		Phone:   service.Phone,
		Address: service.Address,
	}
	err = adres_dao.CreateAddress(address)
	if err != nil {
		code = e.ErrorDatabase
		logrus.Info(err)
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	addresses, err := adres_dao.ListAddressByUid(uid)
	if err != nil {
		code = e.ErrorDatabase
		logrus.Info(err)
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
		Data:   serializer.BuildAddresses(addresses),
	}

}

func (service *AddressService) Show(ctx context.Context, aid string) serializer.Response {
	code := e.SUCCESS
	addressDao := dao.NewAddressDao(ctx)

	addressId, _ := strconv.Atoi(aid)
	address, err := addressDao.GetAddressByAid(uint(addressId))
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
		Data:   serializer.BuildAddress(address),
		Msg:    e.GetMsg(code),
	}
}

func (service *AddressService) List(ctx context.Context, uid uint) serializer.Response {
	code := e.SUCCESS
	addressDao := dao.NewAddressDao(ctx)
	addresses, err := addressDao.ListAddressByUid(uid)
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
		Data:   serializer.BuildAddresses(addresses),
		Msg:    e.GetMsg(code),
	}
}

func (service *AddressService) Update(ctx context.Context, uid uint, aid string) serializer.Response {
	code := e.SUCCESS
	addressId, _ := strconv.Atoi(aid)
	addressDao := dao.NewAddressDao(ctx)

	address := &model.Address{
		UserID:  uid,
		Name:    service.Name,
		Phone:   service.Phone,
		Address: service.Address,
	}
	err := addressDao.UpdateAddressById(uint(addressId), address)
	if err != nil {
		logrus.Info(err)
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	addresses, err := addressDao.ListAddressByUid(uid)
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
		Data:   serializer.BuildAddresses(addresses),
	}
}

func (service *AddressService) Delete(ctx context.Context, aid string) serializer.Response {
	code := e.SUCCESS
	addressDao := dao.NewAddressDao(ctx)
	addressId, _ := strconv.Atoi(aid)
	err := addressDao.DeleteAddressById(uint(addressId))
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
