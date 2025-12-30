package api

import (
	"gin_mall/consts"
	"gin_mall/pkg/utils"
	"gin_mall/service"

	"github.com/gin-gonic/gin"
)

func ListCarousels(c *gin.Context) {
	listCarouselsService := service.ListCarouselsService{}
	if err := c.ShouldBind(&listCarouselsService); err == nil {
		res := listCarouselsService.List()
		c.JSON(consts.StatusOK, res)
	} else {
		c.JSON(consts.IlleageRequest, ErrorResponse(err))
		utils.LogrusObj.Infoln(err)
	}
}
