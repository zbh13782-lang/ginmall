package service

import (
	"context"
	"fmt"
	"gin_mall/conf"
	"gin_mall/consts"
	"gin_mall/pkg/e"
	"gin_mall/pkg/utils"
	"gin_mall/repository/db/dao"
	"gin_mall/repository/db/model"
	"gin_mall/serializer"
	"mime/multipart"
	"strconv"
	"sync"

	"github.com/sirupsen/logrus"
)

type ProductService struct {
	ID            uint   `form:"id" json:"id"`
	Name          string `form:"name" json:"name"`
	CategoryID    int    `form:"category_id" json:"category_id"`
	Title         string `form:"title" json:"title" `
	Info          string `form:"info" json:"info" `
	ImgPath       string `form:"img_path" json:"img_path"`
	Price         string `form:"price" json:"price"`
	DiscountPrice string `form:"discount_price" json:"discount_price"`
	OnSale        bool   `form:"on_sale" json:"on_sale"`
	Num           int    `form:"num" json:"num"`
	model.BasePage
}

type ListProductImgService struct {
}

func (service *ProductService) List(ctx context.Context) serializer.Response {
	var products []*model.Product
	var total int64
	code := e.SUCCESS

	condition := make(map[string]interface{})

	if service.CategoryID != 0 {
		condition["category_id"] = service.CategoryID
	}
	productdao := dao.NewProductDao(ctx)
	total, err := productdao.CountProductByCondition(condition)
	if err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wg.Done()

		productdao = dao.NewProductDaoByDB(productdao.DB)
		products, _ = productdao.ListProductByCondition(condition, service.BasePage)

	}()
	wg.Wait()
	return serializer.BuildListResponse(serializer.BuildProducts(products), uint(total))

}

func (service *ProductService) Show(ctx context.Context, id string) serializer.Response {
	showdao := dao.NewProductDao(ctx)
	var code = e.SUCCESS
	pid, _ := strconv.Atoi(id)
	product, err := showdao.GetProductById(uint(pid))
	if err != nil {
		logrus.Info(err)
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	return serializer.Response{
		Status: code,
		Data:   serializer.BuildProduct(product),
		Msg:    e.GetMsg(code),
	}
}

func (service *ProductService) Search(ctx context.Context) serializer.Response {
	code := e.SUCCESS
	if service.PageSize == 0 {
		service.PageSize = 15
	}
	productdao := dao.NewProductDao(ctx)

	// 先获取总数
	var total int64
	condition := fmt.Sprintf("name LIKE '%%%s%%' OR info LIKE '%%%s%%'", service.Info, service.Info)
	err := productdao.DB.Model(&model.Product{}).Where(condition).Count(&total).Error
	if err != nil {
		logrus.Info(err)
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	// 获取分页结果
	products, err := productdao.SearchProduct(service.Info, service.BasePage)
	if err != nil {
		logrus.Info(err)
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.BuildListResponse(serializer.BuildProducts(products), uint(total))
}

func (service *ProductService) Create(ctx context.Context, uid uint, files []*multipart.FileHeader) serializer.Response {
	code := e.SUCCESS
	var boss *model.User
	var err error
	var path string
	userdao := dao.NewUserDao(ctx)

	boss, err = userdao.GetUserbyid(uid)
	tmp, _ := files[0].Open()
	if conf.UploadModel == consts.UploadModelLocal {
		path, err = utils.UploadProductToLocalStatic(tmp, uid, service.Name)
	} else {
		path, err = utils.UploadToQiNiu(tmp, files[0].Size)
	}
	if err != nil {
		code = e.ErrorUploadFile
		return serializer.Response{
			Status: code,
			Data:   e.GetMsg(code),
			Error:  path,
		}
	}

	product := &model.Product{
		Name:          service.Name,
		CategoryID:    uint(service.CategoryID),
		Title:         service.Title,
		Info:          service.Info,
		ImgPath:       path,
		Price:         service.Price,
		DiscountPrice: service.DiscountPrice,
		Num:           service.Num,
		OnSale:        true,
		BossID:        uid,
		BossName:      boss.UserName,
		BossAvatar:    boss.Avatar,
	}
	productdao := dao.NewProductDao(ctx)
	err = productdao.CreateProduct(product)
	if err != nil {
		logrus.Info(err)
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	// wg := new(sync.WaitGroup)
	// wg.Add(len(files))

	for index, file := range files {

		num := strconv.Itoa(index)
		productimgdao := dao.NewProductImgDaobyDB(productdao.DB)
		tmp, _ = file.Open()
		if conf.UploadModel == consts.UploadModelLocal {
			path, err = utils.UploadProductToLocalStatic(tmp, uid, service.Name+num)
		} else {
			path, err = utils.UploadToQiNiu(tmp, file.Size)
		}
		if err != nil {
			code = e.ErrorUploadFile
			return serializer.Response{
				Status: code,
				Data:   e.GetMsg(code),
				Error:  path,
			}
		}
		productImg := &model.ProductImg{
			ProductID: product.ID,
			ImgPath:   path,
		}
		err = productimgdao.CreateProductImg(productImg)
		if err != nil {
			code = e.ERROR
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}

	}

	return serializer.Response{
		Status: code,
		Data:   serializer.BuildProduct(product),
		Msg:    e.GetMsg(code),
	}
}

func (service *ProductService) Update(ctx context.Context, pid string) serializer.Response {
	productdao := dao.NewProductDao(ctx)
	var (
		code    = e.SUCCESS
		err     error
		product *model.Product
	)
	productID, _ := strconv.Atoi(pid)
	product = &model.Product{
		Name:       service.Name,
		CategoryID: uint(service.CategoryID),
		Title:      service.Title,
		Info:       service.Info,
		// ImgPath:       service.ImgPath,
		Price:         service.Price,
		DiscountPrice: service.DiscountPrice,
		OnSale:        service.OnSale,
	}
	err = productdao.UpdateProduct(uint(productID), product)
	if err != nil {
		logrus.Info(err)
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}

}

func (service *ProductService) Delete(ctx context.Context, pid string) serializer.Response {
	var (
		code = e.SUCCESS
		err  error
	)
	productdao := dao.NewProductDao(ctx)
	productid, _ := strconv.Atoi(pid)
	err = productdao.DeleteProduct(uint(productid))
	if err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}

func (service *ListProductImgService) List(ctx context.Context, pid string) serializer.Response {
	productimgdao := dao.NewProductImgDao(ctx)
	productid, _ := strconv.Atoi(pid)
	productimgs, _ := productimgdao.ListProductImgByProductId(uint(productid))
	return serializer.BuildListResponse(
		serializer.BuildProductImgList(productimgs),
		uint(len(productimgs)),
	)
}
