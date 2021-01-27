package service

import (
	"cmall/model"
	"cmall/pkg/e"
	"cmall/pkg/logging"
	"cmall/serializer"
)

type CreateProductService struct {
	Name         string `form:"name" json:"name"`
	CategortID   int    `form:"category_id" json:"categort_id"`
	Title        string `form:"title" json:"title" bind:"required,min=2,max=100"`
	Info         string `form:"info" json:"info" bind:"max=1000"`
	ImgPath      string `form:"img_path" json:"img_path"`
	Price        string `form:"price" json:"price"`
	DiscoutPrice string `form:"discount_price" json:"discout_price"`
	BossID       int    `form:"boss_id" json:"boss_id" bind:"required"`
	BossName     string `form:"boss_name" json:"boss_name"`
	BossAvatar   string `form:"boss_avatar" json:"boss_avatar"`
}
type ListProductsService struct {
	Limit      int  `form:"limit" json:"limit"`
	Start      int  `form:"start" json:"start"`
	CategoryID uint `form:"category_id" json:"category_id"`
}

//创建商品
func (service *CreateProductService) Create() serializer.Response {
	product := model.Product{
		Name:          service.Name,
		CategoryID:    service.CategortID,
		Title:         service.Title,
		Info:          service.Info,
		ImgPath:       service.ImgPath,
		Price:         service.Price,
		DiscountPrice: service.DiscoutPrice,
		BossID:        service.BossID,
		BossName:      service.BossName,
		BossAvatar:    service.BossAvatar,
	}
	code := e.SUCCESS
	err := model.DB.Create(&product).Error
	if err != nil {
		logging.Info(err)
		code = e.ERROR_DATABASE
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

func (service *ListProductsService) List() serializer.Response {
	products := []model.Product{}

	total := 0
	code := e.SUCCESS

	if service.Limit == 0 {
		service.Limit = 15
	}
	if service.CategoryID == 0 {
		if err := model.DB.Model(model.Product{}).
			Count(&total).Error; err != nil {
			logging.Info(err)
			code = e.ERROR_DATABASE
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}

		if err := model.DB.Limit(service.Limit).
			Offset(service.Start).Find(&products).
			Error; err != nil {
			logging.Info(err)
			code = e.ERROR_DATABASE
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}
	} else {
		if err := model.DB.Model(model.Product{}).
			Where("category_id=?", service.CategoryID).
			Count(&total).Error; err != nil {
			logging.Info(err)
			code = e.ERROR_DATABASE
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}

		if err := model.DB.Where("category_id=?", service.CategoryID).
			Limit(service.Limit).
			Offset(service.Start).
			Find(&products).Error; err != nil {
			logging.Info(err)
			code = e.ERROR_DATABASE
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}
	}

	return serializer.BuildListResponse(serializer.BuildProducts(products), uint(total))
}
