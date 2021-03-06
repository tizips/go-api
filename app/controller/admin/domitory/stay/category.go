package stay

import (
	"github.com/gin-gonic/gin"
	"saas/app/constant"
	"saas/app/model"
	liveForm "saas/app/request/admin/dormitory/stay"
	"saas/app/response/admin/dormitory/stay"
	"saas/kernel/data"
	"saas/kernel/response"
	"strconv"
)

func DoCategoryByCreate(ctx *gin.Context) {

	var request liveForm.DoCategoryByCreate
	if err := ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	category := model.DorStayCategory{
		Name:     request.Name,
		Order:    request.Order,
		IsTemp:   request.IsTemp,
		IsEnable: request.IsEnable,
	}

	if data.Database.Create(&category); category.Id <= 0 {
		response.Fail(ctx, "添加失败")
		return
	}

	response.Success(ctx)
}

func DoCategoryByUpdate(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		response.FailByRequestWithMessage(ctx, "ID获取失败")
		return
	}

	var request liveForm.DoCategoryByUpdate
	if err = ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	var category model.DorStayCategory
	data.Database.Find(&category, id)
	if category.Id <= 0 {
		response.NotFound(ctx, "未找到该类型")
		return
	}

	if category.IsEnable != request.IsEnable {
		var peoples int64 = 0
		data.Database.Model(model.DorPeople{}).Where("`category_id`=? and `status`=?", category.Id, model.DorPeopleStatusLive).Count(&peoples)
		if peoples > 0 {
			response.Fail(ctx, "该类型正在使用，无法上下架")
			return
		}
	}

	category.Name = request.Name
	category.Order = request.Order
	category.IsEnable = request.IsEnable

	if t := data.Database.Save(&category); t.RowsAffected <= 0 {
		response.Fail(ctx, "修改失败")
		return
	}

	response.Success(ctx)
}

func DoCategoryByDelete(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		response.FailByRequestWithMessage(ctx, "ID获取失败")
		return
	}

	var category model.DorStayCategory
	data.Database.Find(&category, id)
	if category.Id <= 0 {
		response.NotFound(ctx, "未找到该类型")
		return
	}

	var peoples int64 = 0
	data.Database.Model(model.DorPeople{}).Where("`category_id`=? and `status`=?", category.Id, model.DorPeopleStatusLive).Count(&peoples)
	if peoples > 0 {
		response.Fail(ctx, "该类型正在使用，无法上下架")
		return
	}

	if t := data.Database.Delete(&category); t.RowsAffected <= 0 {
		response.Fail(ctx, "删除失败")
		return
	}

	response.Success(ctx)
}

func DoCategoryByEnable(ctx *gin.Context) {

	var request liveForm.DoCategoryByEnable
	if err := ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	var category model.DorStayCategory
	data.Database.Find(&category, request.Id)
	if category.Id <= 0 {
		response.NotFound(ctx, "未找到该类型")
		return
	}

	var peoples int64 = 0
	data.Database.Model(model.DorPeople{}).Where("`category_id`=? and `status`=?", category.Id, model.DorPeopleStatusLive).Count(&peoples)
	if peoples > 0 {
		response.Fail(ctx, "该类型正在使用，无法上下架")
		return
	}

	category.IsEnable = request.IsEnable

	if t := data.Database.Save(&category); t.RowsAffected <= 0 {
		response.Fail(ctx, "启禁失败")
		return
	}

	response.Success(ctx)
}

func ToCategoryByList(ctx *gin.Context) {

	responses := make([]any, 0)

	var categories []model.DorStayCategory
	data.Database.Order("`order` asc, `id` desc").Find(&categories)

	for _, item := range categories {
		responses = append(responses, stay.ToCategoryByList{
			Id:        item.Id,
			Name:      item.Name,
			Order:     item.Order,
			IsTemp:    item.IsTemp,
			IsEnable:  item.IsEnable,
			CreatedAt: item.CreatedAt.ToDateTimeString(),
		})
	}

	response.SuccessByList(ctx, responses)
}

func ToCategoryByOnline(ctx *gin.Context) {

	responses := make([]any, 0)

	var categories []model.DorStayCategory
	data.Database.Where("is_enable=?", constant.IsEnableYes).Order("`order` asc,`id` desc").Find(&categories)

	for _, item := range categories {
		responses = append(responses, stay.ToCategoryByOnline{
			Id:     item.Id,
			Name:   item.Name,
			IsTemp: item.IsTemp,
		})
	}

	response.SuccessByList(ctx, responses)
}
