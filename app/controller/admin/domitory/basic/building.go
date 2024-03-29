package basic

import (
	"github.com/gin-gonic/gin"
	"saas/app/constant"
	"saas/app/model"
	"saas/app/request/admin/dormitory/basic"
	res "saas/app/response/admin/dormitory/basic"
	"saas/kernel/app"
	"saas/kernel/response"
	"strconv"
)

func DoBuildingByCreate(ctx *gin.Context) {

	var request basic.DoBuildingByCreate

	if err := ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	building := model.DorBuilding{
		Name:     request.Name,
		Order:    request.Order,
		IsEnable: request.IsEnable,
		IsPublic: request.IsPublic,
	}

	if app.Database.Create(&building); building.Id <= 0 {
		response.Fail(ctx, "添加失败")
		return
	}

	response.Success(ctx)
}

func DoBuildingByUpdate(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil || id <= 0 {
		response.FailByRequestWithMessage(ctx, "ID获取失败")
		return
	}

	var request basic.DoBuildingByUpdate

	if err = ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	var building model.DorBuilding

	if app.Database.Find(&building, id); building.Id <= 0 {
		response.NotFound(ctx, "未找到该楼栋")
		return
	}

	if building.IsEnable != request.IsEnable {

		var peoples int64 = 0

		app.Database.Model(model.DorPeople{}).Where("`building_id`=? and `status`=?", building.Id, model.DorPeopleStatusLive).Count(&peoples)

		if peoples > 0 {
			response.Fail(ctx, "该楼栋已有人入住，无法上下架")
			return
		}
	}

	building.Name = request.Name
	building.Order = request.Order
	building.IsEnable = request.IsEnable

	if t := app.Database.Save(&building); t.RowsAffected <= 0 {
		response.Fail(ctx, "修改失败")
		return
	}

	response.Success(ctx)
}

func DoBuildingByDelete(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil || id <= 0 {
		response.FailByRequestWithMessage(ctx, "ID获取失败")
		return
	}

	var building model.DorBuilding

	if app.Database.Find(&building, id); building.Id <= 0 {
		response.NotFound(ctx, "未找到该楼栋")
		return
	}

	var peoples int64 = 0

	app.Database.Model(model.DorPeople{}).Where("`building_id`=? and `status`=?", building.Id, model.DorPeopleStatusLive).Count(&peoples)

	if peoples > 0 {
		response.Fail(ctx, "该楼栋已有人入住，无法删除")
		return
	}

	tx := app.Database.Begin()

	if t := tx.Delete(&building); t.RowsAffected <= 0 {
		tx.Rollback()
		response.Fail(ctx, "删除失败")
		return
	}

	if t := tx.Delete(&model.DorRoom{}, "building_id=?", building.Id); t.Error != nil {
		tx.Rollback()
		response.Fail(ctx, "删除失败")
		return
	}

	if t := tx.Delete(&model.DorRoom{}, "building_id=?", building.Id); t.Error != nil {
		tx.Rollback()
		response.Fail(ctx, "删除失败")
		return
	}

	if t := tx.Delete(&model.DorBed{}, "building_id=?", building.Id); t.Error != nil {
		tx.Rollback()
		response.Fail(ctx, "删除失败")
		return
	}

	tx.Commit()

	response.Success(ctx)
}

func DoBuildingByEnable(ctx *gin.Context) {

	var request basic.DoBuildingByEnable

	if err := ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	var building model.DorBuilding

	if app.Database.Find(&building, request.Id); building.Id <= 0 {
		response.NotFound(ctx, "未找到该楼栋")
		return
	}

	var peoples int64 = 0

	app.Database.Model(model.DorPeople{}).Where("`building_id`=? and `status`=?", building.Id, model.DorPeopleStatusLive).Count(&peoples)

	if peoples > 0 {
		response.Fail(ctx, "该楼栋已有人入住，无法上下架")
		return
	}

	building.IsEnable = request.IsEnable

	if t := app.Database.Save(&building); t.RowsAffected <= 0 {
		response.Fail(ctx, "启禁失败")
		return
	}

	response.Success(ctx)
}

func ToBuildingByList(ctx *gin.Context) {

	responses := make([]res.ToBuildingByList, 0)

	var buildings []model.DorBuilding

	app.Database.Order("`order` asc, `id` desc").Find(&buildings)

	for _, item := range buildings {
		responses = append(responses, res.ToBuildingByList{
			Id:        item.Id,
			Name:      item.Name,
			Order:     item.Order,
			IsEnable:  item.IsEnable,
			IsPublic:  item.IsPublic,
			CreatedAt: item.CreatedAt.ToDateTimeString(),
		})
	}

	response.SuccessByData(ctx, responses)
}

func ToBuildingByOnline(ctx *gin.Context) {

	var request basic.ToBuildingByOnline

	if err := ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	responses := make([]res.ToBuildingByOnline, 0)

	tx := app.Database.Where("`is_enable`=?", constant.IsEnableYes)

	if request.IsPublic > 0 {
		tx = tx.Where("`is_public`=?", request.IsPublic)
	}

	var buildings []model.DorBuilding
	tx.Order("`order` asc, `id` desc").Find(&buildings)

	for _, item := range buildings {
		items := res.ToBuildingByOnline{
			Id:   item.Id,
			Name: item.Name,
		}
		if request.WithPublic {
			items.IsPublic = item.IsPublic
		}
		responses = append(responses, items)
	}

	response.SuccessByData(ctx, responses)
}
