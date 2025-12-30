package serializer

import (
	"gin_mall/pkg/utils"
	"gin_mall/repository/db/model"
)

type Money struct {
	UserID    uint   `json:"user_id" form:"user_id"`
	UserName  string `json:"user_name" form:"user_name"`
	UserMoney string `json:"user_money" form:"user_money"`
}

func BuildMoney(item *model.User, key string) Money {
	utils.Encrypt.SetKey(key)
	money, _ := utils.Encrypt.AesDecoding(item.Money)
	return Money{
		UserID:    item.ID,
		UserName:  item.UserName,
		UserMoney: money,
	}
}
