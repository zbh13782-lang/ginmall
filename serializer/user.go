package serializer

import (
	"gin_mall/repository/db/model"
)

type User struct {
	ID       uint   `json:"id"`
	UserName string `json:"username"`
	PassWord string `josn:"password"`
	Email    string `json:"email"`
	NickName string `json:"nickname"`
	Status   string `json:"status"`
	Avater   string `json:"avater"`
	Money    string `json:"money"`
	CreateAt int64  `json:"createat"`
}

func BuildUser(user *model.User) User {
	return User{
		ID:       user.ID,
		UserName: user.UserName,
		PassWord: user.PasswordDigest,
		Email:    user.Email,
		NickName: user.NickName,
		Status:   user.Status,
		Avater:   user.Avatar,
		Money:    user.Money,
		CreateAt: user.CreatedAt.Unix(),
	}

}
