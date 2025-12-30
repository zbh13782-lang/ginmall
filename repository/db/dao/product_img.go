package dao

import (
	"context"
	"gin_mall/repository/db/model"

	"gorm.io/gorm"
)

type ProductImgDao struct {
	*gorm.DB
}

func NewProductImgDao(ctx context.Context) *ProductImgDao {
	return &ProductImgDao{NewDBClient(ctx)}
}

func NewProductImgDaobyDB(db *gorm.DB) *ProductImgDao {
	return &ProductImgDao{db}
}

func (dao *ProductImgDao) CreateProductImg(pmg *model.ProductImg) error {
	err := dao.DB.Model(&model.ProductImg{}).Create(&pmg).Error
	return err
}

// ListProductImgByProductId 根据商品id获取商品图片
func (dao *ProductImgDao) ListProductImgByProductId(pId uint) (products []*model.ProductImg, err error) {
	err = dao.DB.Model(&model.ProductImg{}).
		Where("product_id=?", pId).Find(&products).Error
	return
}
