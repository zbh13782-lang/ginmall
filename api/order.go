package api

import (
	"gin_mall/consts"
	"gin_mall/pkg/utils"
	"gin_mall/service"

	"github.com/gin-gonic/gin"
)

func CreateOrder(c *gin.Context) {
	claims, _ := utils.ParseToken(c.GetHeader("Authorization"))
	OrderService := service.OrderService{}
	if err := c.ShouldBind(&OrderService); err == nil {
		res := OrderService.Create(c.Request.Context(), claims.ID)
		c.JSON(consts.StatusOK, res)
	} else {
		c.JSON(consts.IlleageRequest, ErrorResponse(err))
		utils.LogrusObj.Info(err)
	}
}

func ListOrders(c *gin.Context) {
	claims, _ := utils.ParseToken(c.GetHeader("Authorization"))
	OrderService := service.OrderService{}
	if err := c.ShouldBind(&OrderService); err == nil {
		res := OrderService.List(c.Request.Context(), claims.ID)
		c.JSON(consts.StatusOK, res)
	} else {
		c.JSON(consts.IlleageRequest, ErrorResponse(err))
		utils.LogrusObj.Info(err)
	}
}

func ShowOrder(c *gin.Context) {
	showOrderService := service.OrderService{}
	if err := c.ShouldBind(&showOrderService); err == nil {
		res := showOrderService.Show(c.Request.Context(), c.Param("id"))
		c.JSON(consts.StatusOK, res)
	} else {
		c.JSON(consts.IlleageRequest, ErrorResponse(err))
		utils.LogrusObj.Infoln(err)
	}
}

func DeleteOrder(c *gin.Context) {
	deleteOrderService := service.OrderService{}
	if err := c.ShouldBind(&deleteOrderService); err == nil {
		res := deleteOrderService.Delete(c.Request.Context(), c.Param("id"))
		c.JSON(consts.StatusOK, res)
	} else {
		c.JSON(consts.IlleageRequest, ErrorResponse(err))
		utils.LogrusObj.Infoln(err)
	}
}
