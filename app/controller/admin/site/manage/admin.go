package manage

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"saas/app/constant"
	"saas/app/model"
	"saas/app/request/admin/site/manage"
	res "saas/app/response/admin/site/manage"
	"saas/kernel/app"
	"saas/kernel/authorize"
	"saas/kernel/response"
	"strconv"
)

func DoAdminByCreate(ctx *gin.Context) {

	var request manage.DoAdminByCreate

	if err := ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	var count int64

	tc := app.Database.Model(model.SysRole{})

	if !authorize.Root(authorize.Id(ctx)) {
		tc = tc.Where("`id`<>?", authorize.ROOT)
	}

	tc.Where("`id` IN (?)", request.Roles).Count(&count)

	if len(request.Roles) != int(count) {
		response.NotFound(ctx, "部分角色不存在")
		return
	}

	app.Database.Model(model.SysAdmin{}).Where("`mobile`=?", request.Mobile).Count(&count)

	if count > 0 {
		response.Fail(ctx, "该手机号已被注册")
		return
	}

	app.Database.Model(model.SysAdmin{}).Where("username = ?", request.Username).Count(&count)

	if count > 0 {
		response.Fail(ctx, "该用户名已被注册")
		return
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)

	tx := app.Database.Begin()

	admin := model.SysAdmin{
		Username: request.Username,
		Nickname: request.Nickname,
		Mobile:   request.Mobile,
		Password: string(password),
		IsEnable: request.IsEnable,
	}

	if tx.Create(&admin); admin.Id <= 0 {
		tx.Rollback()
		response.Fail(ctx, "添加失败")
		return
	}

	var binds []model.SysAdminBindRole

	for _, item := range request.Roles {
		binds = append(binds, model.SysAdminBindRole{
			AdminId: admin.Id,
			RoleId:  item,
		})
	}

	if t := tx.CreateInBatches(binds, 100); t.RowsAffected <= 0 {
		tx.Rollback()
		response.Fail(ctx, "添加失败")
		return
	}

	if len(binds) > 0 {

		var items = make([]string, len(binds))

		for idx, item := range binds {
			items[idx] = authorize.NameByRole(item.RoleId)
		}

		if _, err := app.Casbin.AddRolesForUser(authorize.NameByAdmin(admin.Id), items); err != nil {
			tx.Rollback()
			response.Fail(ctx, "添加失败")
			return
		}
	}

	tx.Commit()

	response.Success(ctx)
}

func DoAdminByUpdate(ctx *gin.Context) {

	id, _ := strconv.Atoi(ctx.Param("id"))

	if id <= 0 {
		response.FailByRequestWithMessage(ctx, "ID不存在")
		return
	}

	var request manage.DoAdminByUpdate

	if err := ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	var count int64

	tc := app.Database.Model(model.SysRole{})

	if !authorize.Root(authorize.Id(ctx)) {
		tc = tc.Where("`id`<>?", authorize.ROOT)
	}

	tc.Where("`id` in (?)", request.Roles).Count(&count)

	if len(request.Roles) != int(count) {
		response.NotFound(ctx, "部分角色不存在")
		return
	}

	app.Database.Model(model.SysAdmin{}).Where("`id`<>? and `mobile`=?", id, request.Mobile).Count(&count)

	if count > 0 {
		response.Fail(ctx, "该手机号已被注册")
		return
	}

	var admin model.SysAdmin

	if app.Database.Preload("BindRoles").Find(&admin, id); admin.Id <= 0 {
		response.Fail(ctx, "该账号不存在")
		return
	}

	admin.Nickname = request.Nickname
	admin.Mobile = request.Mobile
	admin.IsEnable = request.IsEnable

	if request.Password != "" {

		password, _ := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)

		admin.Password = string(password)
	}

	var creates []model.SysAdminBindRole
	var deletes []int
	var del []int

	for _, item := range request.Roles {

		mark := true

		for _, value := range admin.BindRoles {
			if item == value.RoleId {
				mark = false
				break
			}
		}

		if mark {
			creates = append(creates, model.SysAdminBindRole{
				AdminId: admin.Id,
				RoleId:  item,
			})
		}
	}
	for _, item := range admin.BindRoles {

		mark := true

		for _, value := range request.Roles {
			if item.RoleId == value {
				mark = false
				break
			}
		}

		if mark {
			del = append(del, item.RoleId)
			deletes = append(deletes, item.Id)
		}
	}

	tx := app.Database.Begin()

	if t := tx.Save(&admin); t.RowsAffected <= 0 {
		tx.Rollback()
		response.Fail(ctx, "修改失败")
		return
	}

	if request.IsEnable != constant.IsEnableYes { //	用户禁用，删除缓存角色
		if _, err := app.Casbin.DeleteRolesForUser(authorize.NameByAdmin(admin.Id)); err != nil {
			tx.Rollback()
			response.Fail(ctx, "修改失败")
			return
		}
	} else {

		var items = make([]string, len(request.Roles))

		for idx, item := range request.Roles {
			items[idx] = authorize.NameByRole(item)
		}

		if _, err := app.Casbin.AddRolesForUser(authorize.NameByAdmin(admin.Id), items); err != nil {
			tx.Rollback()
			response.Fail(ctx, "修改失败")
			return
		}
	}

	if len(deletes) > 0 {

		var bindings model.SysAdminBindRole

		if t := tx.Where("`admin_id`=?", admin.Id).Delete(&bindings, deletes); t.RowsAffected <= 0 {
			tx.Rollback()
			response.Fail(ctx, "修改失败")
			return
		}

		if len(del) > 0 && request.IsEnable == constant.IsEnableYes { //	用户启用，结算需要删除的角色
			for _, item := range del {
				if _, err := app.Casbin.DeleteRoleForUser(authorize.NameByAdmin(admin.Id), authorize.NameByRole(item)); err != nil {
					tx.Rollback()
					response.Fail(ctx, "修改失败")
					return
				}
			}
		}
	}

	if len(creates) > 0 {

		if t := tx.CreateInBatches(creates, 100); t.RowsAffected <= 0 {
			tx.Rollback()
			response.Fail(ctx, "修改失败")
			return
		}

		if len(creates) > 0 && request.IsEnable == constant.IsEnableYes { //	用户启用，处理需要新加的角色

			var items = make([]string, len(creates))

			for idx, item := range creates {
				items[idx] = authorize.NameByRole(item.RoleId)
			}

			if _, err := app.Casbin.AddRolesForUser(authorize.NameByAdmin(admin.Id), items); err != nil {
				tx.Rollback()
				response.Fail(ctx, "修改失败")
				return
			}
		}
	}

	tx.Commit()

	response.Success(ctx)
}

func ToAdminByPaginate(ctx *gin.Context) {

	var request manage.ToAdminByPaginate

	if err := ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	tx := app.Database

	if !authorize.Root(authorize.Id(ctx)) {
		tx = tx.Where("not exists (?)", app.Database.
			Select("1").
			Model(model.SysAdminBindRole{}).
			Where(fmt.Sprintf("%s.id=%s.admin_id", model.TableSysAdmin, model.TableSysAdminBindRole)).
			Where("role_id = ?", authorize.ROOT),
		)
	}

	responses := response.Paginate[res.ToAdminByPaginate]{
		Total: 0,
		Page:  request.GetPage(),
		Size:  request.GetSize(),
		Data:  make([]res.ToAdminByPaginate, 0),
	}

	tc := tx

	tc.Model(model.SysAdmin{}).Count(&responses.Total)

	if responses.Total > 0 {

		var admins []model.SysAdmin

		tx.
			//Preload("BindRoles").
			Preload("BindRoles.Role").
			Order("`id` desc").
			Offset(request.GetOffset()).
			Limit(request.GetLimit()).
			Find(&admins)

		for _, item := range admins {

			result := res.ToAdminByPaginate{
				Id:        item.Id,
				Username:  item.Username,
				Nickname:  item.Nickname,
				Mobile:    item.Mobile,
				IsEnable:  item.IsEnable,
				CreatedAt: item.CreatedAt.ToDateTimeString(),
			}

			for _, value := range item.BindRoles {
				result.Roles = append(result.Roles, struct {
					Id   int    `json:"id"`
					Name string `json:"name"`
				}{Id: value.Role.Id, Name: value.Role.Name})
			}

			responses.Data = append(responses.Data, result)
		}
	}

	response.SuccessByPaginate(ctx, responses)
}

func DoAdminByDelete(ctx *gin.Context) {

	id, _ := strconv.Atoi(ctx.Param("id"))

	if id <= 0 {
		response.FailByRequestWithMessage(ctx, "ID不存在")
		return
	}

	if authorize.Id(ctx) == id {
		response.Fail(ctx, "无法删除自身账号")
		return
	}

	var admin model.SysAdmin

	if app.Database.Find(&admin, id); admin.Id <= 0 {
		response.NotFound(ctx, "账号不存在")
		return
	}

	tx := app.Database.Begin()

	if t := app.Database.Delete(&admin); t.RowsAffected <= 0 {
		tx.Rollback()
		response.Fail(ctx, "账号删除失败")
		return
	}

	bind := model.SysAdminBindRole{AdminId: admin.Id}

	if t := tx.Delete(&bind, "`admin_id`=?", admin.Id); t.RowsAffected <= 0 {
		tx.Rollback()
		response.Fail(ctx, "账号删除失败")
		return
	}

	if _, err := app.Casbin.DeleteRolesForUser(authorize.NameByAdmin(admin.Id)); err != nil {
		tx.Rollback()
		response.Fail(ctx, "账号删除失败")
		return
	}

	tx.Commit()

	response.Success(ctx)
}

func DoAdminByEnable(ctx *gin.Context) {

	var request manage.DoAdminByEnable

	if err := ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	var admin model.SysAdmin

	if app.Database.Find(&admin, request.Id); admin.Id <= 0 {
		response.NotFound(ctx, "账号不存在")
		return
	}

	admin.IsEnable = request.IsEnable

	tx := app.Database.Begin()

	if t := app.Database.Save(&admin); t.RowsAffected <= 0 {
		tx.Rollback()
		response.Fail(ctx, "启禁失败")
		return
	}

	if request.IsEnable == constant.IsEnableNo {
		if _, err := app.Casbin.DeleteRolesForUser(authorize.NameByAdmin(admin.Id)); err != nil {
			tx.Rollback()
			response.Fail(ctx, "启禁失败")
			return
		}
	} else if request.IsEnable == constant.IsEnableYes {

		tx.Find(&admin.BindRoles, "`admin_id`=?", admin.Id)

		if len(admin.BindRoles) > 0 {

			var items []string

			for _, item := range admin.BindRoles {
				items = append(items, authorize.NameByRole(item.RoleId))
			}

			if _, err := app.Casbin.AddRolesForUser(authorize.NameByAdmin(admin.Id), items); err != nil {
				tx.Rollback()
				response.Fail(ctx, "启禁失败")
				return
			}
		}
	}

	tx.Commit()

	response.Success(ctx)
}
