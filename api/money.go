package api

import (
	"gin_mall/consts"
	"gin_mall/pkg/utils"
	"gin_mall/service"

	"github.com/gin-gonic/gin"
)

func ShowMoney(c *gin.Context){
	service:=service.ShowMoneyService{}
	claim,_:=utils.ParseToken(c.GetHeader("Authorization"))
	if err:=c.ShouldBind(&service);err==nil{
		res:=service.Show(c.Request.Context(),claim.ID)
		c.JSON(consts.StatusOK,res)
	}else{
		utils.LogrusObj.Info(err)
		c.JSON(consts.IlleageRequest,ErrorResponse(err))
	}
}