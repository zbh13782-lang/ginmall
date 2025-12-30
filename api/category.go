package api

import (
	"gin_mall/consts"
	"gin_mall/pkg/utils"
	"gin_mall/service"

	"github.com/gin-gonic/gin"
)

func ListCategories(c *gin.Context) {
	listCategoriesService := service.ListCategoryService{}
	if err := c.ShouldBind(&listCategoriesService); err == nil {
		res := listCategoriesService.List(c.Request.Context())
		c.JSON(consts.StatusOK, res)
	} else {
		c.JSON(consts.IlleageRequest, ErrorResponse(err))
		utils.LogrusObj.Infoln(err)
	}
}
