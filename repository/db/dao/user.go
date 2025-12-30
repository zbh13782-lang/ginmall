package dao

import (
	"context"
	"gin_mall/repository/db/model"

	"gorm.io/gorm"
)

type UserDao struct {
	*gorm.DB
}

func NewUserDao(ctx context.Context) *UserDao {
	return &UserDao{NewDBClient(ctx)}
}

func NewUserDaobyDB(db *gorm.DB) *UserDao {
	return &UserDao{db}
}

func (dao UserDao) GetUserbyid(id uint) (u *model.User, e error) {
	e = dao.DB.Model(&model.User{}).Where("id=?", id).First(&u).Error
	return
}

func (dao UserDao) UpdateUserbyid(id uint, u *model.User) error {
	return dao.DB.Model(&model.User{}).Where("id=?", id).Updates(&u).Error
}

func (dao UserDao) ExistOrNotbyName(name string) (user *model.User, exist bool, err error) {
	var cnt int64
	err = dao.DB.Model(&model.User{}).Where("user_name=?", name).Count(&cnt).Error
	if cnt == 0 {
		return user, false, err
	}

	err = dao.DB.Model(&model.User{}).Where("user_name=?", name).First(&user).Error
	if err != nil {
		return user, false, err
	}
	return user, true, err
}

func (dao UserDao) CreateUser(u *model.User) error {
	err := dao.DB.Model(&model.User{}).Create(&u).Error
	return err
}
