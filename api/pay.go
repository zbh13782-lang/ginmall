package api

import (
	"gin_mall/consts"
	"gin_mall/pkg/utils"
	"gin_mall/service"

	"github.com/gin-gonic/gin"
)

func OrderPay(c *gin.Context){
	service:=service.OrderPay{}
	claim,_:=utils.ParseToken(c.GetHeader("Authorization"))
	if err:=c.ShouldBind(&service);err==nil{
		res:=service.PayDown(c.Request.Context(),claim.ID)
		c.JSON(consts.StatusOK,res)
	}else{
		utils.LogrusObj.Info(err)
		c.JSON(consts.IlleageRequest,ErrorResponse(err))
	}
}