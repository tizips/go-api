package live

import "saas/app/form/basic"

type DoCategoryByCreateForm struct {
	Name   string `form:"name" json:"name" binding:"required,max=20"`
	Order  uint   `form:"order" json:"order" binding:"required,gte=1,lte=99"`
	IsTemp uint8  `form:"is_temp" json:"is_temp" binding:"oneof=0 1"`
	basic.Enable
}

type DoCategoryByUpdateForm struct {
	Name   string `form:"name" json:"name" binding:"required,max=20"`
	Order  uint   `form:"order" json:"order" binding:"required,gte=1,lte=99"`
	IsTemp uint8  `form:"is_temp" json:"is_temp" binding:"oneof=0 1"`
	basic.Enable
}

type DoCategoryByEnableForm struct {
	Id uint `form:"id" json:"id" binding:"required,gt=0"`
	basic.Enable
}
