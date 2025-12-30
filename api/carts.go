package api

import (
	"gin_mall/consts"
	"gin_mall/pkg/utils"
	"gin_mall/service"

	"github.com/gin-gonic/gin"
)

func CreateCart(c *gin.Context) {
	service := service.CartService{}
	claims, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&service); err == nil {
		res := service.Create(c.Request.Context(), claims.ID)
		c.JSON(consts.StatusOK, res)
	} else {
		utils.LogrusObj.Info(err)
		c.JSON(consts.IlleageRequest, ErrorResponse(err))
	}
}

func ShowCarts(c *gin.Context) {
	service := service.CartService{}
	claims, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&service); err == nil {
		res := service.Show(c.Request.Context(), claims.ID)
		c.JSON(consts.StatusOK, res)
	} else {
		utils.LogrusObj.Info(err)
		c.JSON(consts.IlleageRequest, ErrorResponse(err))
	}
}

func UpdateCart(c *gin.Context) {
	service := service.CartService{}
	cartid := c.Param("id")
	if err := c.ShouldBind(&service); err == nil {
		res := service.Update(c.Request.Context(), cartid)
		c.JSON(consts.StatusOK, res)
	} else {
		utils.LogrusObj.Info(err)
		c.JSON(consts.IlleageRequest, ErrorResponse(err))
	}
}

func DeleteCart(c *gin.Context) {
	service := service.CartService{}
	if err := c.ShouldBind(&service); err == nil {
		res := service.Delete(c.Request.Context())
		c.JSON(consts.StatusOK, res)
	} else {
		utils.LogrusObj.Info(err)
		c.JSON(consts.IlleageRequest, ErrorResponse(err))
	}
}
