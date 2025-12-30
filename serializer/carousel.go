package serializer

import "gin_mall/repository/db/model"

type Carousel struct {
	ID        uint   `json:"id"`
	ImgPath   string `json:"img_path"`
	ProductID uint   `json:"product_id"`
	CreatedAt int64  `json:"created_at"`
}

func BuildCarousel(item *model.Carousel) Carousel{
	return Carousel{
		ID: item.ID,
		ImgPath: item.ImgPath,
		ProductID: item.ProductID,
		CreatedAt: item.CreatedAt.Unix(),
	}
}

func BuildCarousels(items []*model.Carousel) (carousels []Carousel) {
	for _, item := range items {
		carousel := BuildCarousel(item)
		carousels = append(carousels, carousel)
	}
	return carousels
}