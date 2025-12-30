package api

import (
	"gin_mall/consts"
	"gin_mall/pkg/utils"
	"gin_mall/service"

	"github.com/gin-gonic/gin"
)

func CreateFavorite(c *gin.Context){
	service:=service.FavoriteService{}
	claims,_:=utils.ParseToken(c.GetHeader("Authorization"))
	if err:=c.ShouldBind(&service);err==nil{
		res:=service.Create(c.Request.Context(),claims.ID)
		c.JSON(consts.StatusOK,res)
	}else{
		c.JSON(consts.IlleageRequest,ErrorResponse(err))
		utils.LogrusObj.Info(err)
	}
}

func ShowFavorites(c *gin.Context){
	service:=service.FavoriteService{}
	claims,_:=utils.ParseToken(c.GetHeader("Authorization"))
	if err:=c.ShouldBind(&service);err==nil{
		res:=service.Show(c.Request.Context(),claims.ID)
		c.JSON(consts.StatusOK,res)
	}else{
		c.JSON(consts.IlleageRequest,ErrorResponse(err))
		utils.LogrusObj.Info(err)
	}	
}

func DeleteFavorite(c *gin.Context){
	service:=service.FavoriteService{}
	if err:=c.ShouldBind(&service);err==nil{
		res:=service.Delete(c.Request.Context())
		c.JSON(consts.StatusOK,res)
	}else{
		c.JSON(consts.IlleageRequest,ErrorResponse(err))
		utils.LogrusObj.Info(err)
	}	
}