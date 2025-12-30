package serializer

import (
	"gin_mall/conf"
	"gin_mall/consts"
	"gin_mall/repository/db/model"
)

type ProductImg struct {
	ImgPath   string `json:"img_path" form:"img_path"`
	ProductID uint   `json:"product_id" form:"product_id"`
}

func BuildProductImg(item *model.ProductImg) ProductImg {
	pimg := ProductImg{
		ImgPath:   conf.PhotoHost + conf.HttpPort + conf.ProductPhotoPath + item.ImgPath,
		ProductID: item.ProductID,
	}
	if conf.UploadModel == consts.UploadModelOss {
		pimg.ImgPath = item.ImgPath
	}
	return pimg
}

func BuildProductImgList(item []*model.ProductImg) []ProductImg {
	imgs := []ProductImg{}
	for _, img := range item {
		imgs = append(imgs, BuildProductImg(img))
	}
	return imgs
}
