package api

import (
	"gin_mall/consts"
	"gin_mall/pkg/utils"
	"gin_mall/service"

	"github.com/gin-gonic/gin"
)

func ListProducts(c *gin.Context) {
	listservice := service.ProductService{}
	if err := c.ShouldBind(&listservice); err != nil {
		c.JSON(consts.IlleageRequest, ErrorResponse(err))
		utils.LogrusObj.Info(err)
	} else {
		res := listservice.List(c.Request.Context())
		c.JSON(consts.StatusOK, res)
	}
}

func ShowProduct(c *gin.Context) {
	showservice := service.ProductService{}
	res := showservice.Show(c.Request.Context(), c.Param("id"))
	c.JSON(consts.StatusOK, res)
}

func SearchProducts(c *gin.Context) {
	searchProductsService := service.ProductService{}
	if err := c.ShouldBind(&searchProductsService); err == nil {
		res := searchProductsService.Search(c.Request.Context())
		c.JSON(consts.StatusOK, res)
	} else {
		c.JSON(consts.IlleageRequest, ErrorResponse(err))
		utils.LogrusObj.Infoln(err)
	}
}

func CreateProduct(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["file"]
	claims, _ := utils.ParseToken(c.GetHeader("Authorization"))
	CreateProductService := service.ProductService{}
	if err := c.ShouldBind(&CreateProductService); err == nil {
		res := CreateProductService.Create(c.Request.Context(), claims.ID, files)
		c.JSON(consts.StatusOK, res)
	} else {
		c.JSON(consts.IlleageRequest, ErrorResponse(err))
		utils.LogrusObj.Infoln(err)
	}
}

func UpdateProduct(c *gin.Context) {
	updateservice := service.ProductService{}
	if err := c.ShouldBind(&updateservice); err == nil {
		res := updateservice.Update(c.Request.Context(), c.Param("id"))
		c.JSON(consts.StatusOK, res)
	} else {
		c.JSON(consts.IlleageRequest, ErrorResponse(err))
		utils.LogrusObj.Info(err)
	}
}

func DeleteProduct(c *gin.Context) {
	deleteservice := service.ProductService{}
	if err := c.ShouldBind(&deleteservice); err == nil {
		res := deleteservice.Delete(c.Request.Context(), c.Param("id"))
		c.JSON(consts.StatusOK, res)
	} else {
		c.JSON(consts.IlleageRequest, ErrorResponse(err))
		utils.LogrusObj.Info(err)
	}
}

func ListProductImg(c *gin.Context) {
	var listimgservice service.ListProductImgService
	if err := c.ShouldBind(&listimgservice); err == nil {
		res := listimgservice.List(c.Request.Context(), c.Param("id"))
		c.JSON(consts.StatusOK, res)
	} else {
		c.JSON(consts.IlleageRequest, ErrorResponse(err))
		utils.LogrusObj.Info(err)
	}
}
