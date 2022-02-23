package architecture

import "saas/app/form/basic"

type DoModuleByCreateForm struct {
	Slug  string `form:"slug" json:"slug" binding:"required,min=2,max=20,alpha"`
	Name  string `form:"name" json:"name" binding:"required,min=2,max=20"`
	Order uint   `form:"order" json:"order" binding:"omitempty,gt=1,lt=99"`
	basic.Enable
}

type DoModuleByUpdateForm struct {
	Slug  string `form:"slug" json:"slug" binding:"required,min=2,max=20,alpha"`
	Name  string `form:"name" json:"name" binding:"required,min=2,max=20"`
	Order uint   `form:"order" json:"order" binding:"omitempty,gt=1,lt=99"`
	basic.Enable
}

type DoModuleByEnableForm struct {
	Id uint `form:"id" json:"id" binding:"required,gt=0"`
	basic.Enable
}
