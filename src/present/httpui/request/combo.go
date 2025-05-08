package request

import "mime/multipart"

type CreateComboRequest struct {
	Name        string                `form:"name" binding:"required"`
	Description string                `form:"description" binding:"required"`
	BannerUrl   *multipart.FileHeader `form:"banner_url"`
	Price       float64               `form:"price" binding:"required"`
}

type UpdateComboRequest struct {
	Name        string                `form:"name" binding:"required"`
	Description string                `form:"description" binding:"required"`
	BannerUrl   *multipart.FileHeader `form:"banner_url"`
	Price       float64               `form:"price" binding:"required"`
}
