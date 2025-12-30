package api

import (
	"gin_mall/consts"
	"gin_mall/pkg/utils"
	"gin_mall/service"

	"github.com/gin-gonic/gin"
)

func UserRegister(c *gin.Context) {
	var userregister service.UserRegisterService
	if err := c.ShouldBind(&userregister); err == nil {
		res := userregister.Register(c.Request.Context())
		c.JSON(consts.StatusOK, res)
	} else {
		c.JSON(consts.IlleageRequest, ErrorResponse(err))
		utils.LogrusObj.Infoln(err)
	}
}

func UserLogin(c *gin.Context) {
	var userlogin service.UserLoginService
	if err := c.ShouldBind(&userlogin); err == nil {
		res := userlogin.Login(c.Request.Context())
		c.JSON(consts.StatusOK, res)
	} else {
		c.JSON(consts.IlleageRequest, ErrorResponse(err))
		utils.LogrusObj.Infoln(err)
	}
}

func UserUpdate(c *gin.Context) {
	var userupdate service.UserUpdateService
	claims, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&userupdate); err == nil {
		res := userupdate.Update(c.Request.Context(), claims.ID)
		c.JSON(consts.StatusOK, res)
	} else {
		c.JSON(consts.IlleageRequest, ErrorResponse(err))
		utils.LogrusObj.Infoln(err)
	}
}

func SendEmail(c *gin.Context) {
	var sendEmailService service.SendEmailService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&sendEmailService); err == nil {
		res := sendEmailService.SendEmail(c.Request.Context(), claim.ID)
		c.JSON(consts.StatusOK, res)
	} else {
		c.JSON(consts.IlleageRequest, ErrorResponse(err))
		utils.LogrusObj.Info(err)
	}
}

func ValidEmail(c *gin.Context) {
	var validemailservice service.ValidEmailService
	if err := c.ShouldBind(validemailservice); err != nil {
		res := validemailservice.Valid(c.Request.Context(), c.GetHeader("Authorization"))
		c.JSON(consts.StatusOK, res)
	} else {
		c.JSON(consts.IlleageRequest, ErrorResponse(err))
		utils.LogrusObj.Info(err)
	}
}

func UploadAvatar(c *gin.Context) {
	file, fileheader, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(consts.IlleageRequest, ErrorResponse(err))
		utils.LogrusObj.Info(err)
		return
	}
	filesize := fileheader.Size
	// 重置文件指针到开始位置，确保服务层能正确读取文件
	file.Seek(0, 0)
	uploadservice := service.UserUploadService{}
	claims, _ := utils.ParseToken(c.GetHeader("Authorization"))
	res := uploadservice.Post(c.Request.Context(), claims.ID, file, filesize)
	c.JSON(consts.StatusOK, res)
}
