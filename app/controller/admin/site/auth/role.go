package auth

import (
	"github.com/gin-gonic/gin"
	"saas/app/model"
	"saas/app/request/admin/site/auth"
	authResponse "saas/app/response/admin/site/auth"
	"saas/kernel/authorize"
	"saas/kernel/data"
	"saas/kernel/response"
	"strconv"
)

func DoRoleByCreate(ctx *gin.Context) {

	var request auth.DoRoleByCreate
	if err := ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	var permissionsIds []int

	var modules, children1, children2, children3 []int

	for _, item := range request.Permissions {
		if len(item) >= 4 {
			children3 = append(children3, item[3])
		} else if len(item) >= 3 {
			children2 = append(children2, item[2])
		} else if len(item) >= 2 {
			children1 = append(children1, item[1])
		} else if len(item) >= 1 {
			modules = append(modules, item[0])
		}
	}

	if len(modules) > 0 {
		var permissions []model.SysPermission
		data.Database.Where("`module_id` in (?) and `method`<>? and `path`<>?", modules, "", "").Find(&permissions)
		for _, item := range permissions {
			permissionsIds = append(permissionsIds, item.Id)
		}
	}
	if len(children3) > 0 {
		var permissions []model.SysPermission
		data.Database.Where("`id` in (?) and `method`<>? and `path`<>?", children3, "", "").Find(&permissions)
		for _, item := range permissions {
			permissionsIds = append(permissionsIds, item.Id)
		}
	}
	if len(children2) > 0 {
		var permissions []model.SysPermission
		data.Database.Where("`parent_i2` in (?) and `method`<>? and `path`<>?", children2, "", "").Find(&permissions)
		for _, item := range permissions {
			permissionsIds = append(permissionsIds, item.Id)
		}
	}
	if len(children1) > 0 {
		var permissions []model.SysPermission
		data.Database.Where("`parent_i1` in (?) and `method`<>? and `path`<>?", children1, "", "").Find(&permissions)
		for _, item := range permissions {
			permissionsIds = append(permissionsIds, item.Id)
		}
	}

	var temp = make(map[int]int, len(permissionsIds))
	for _, item := range permissionsIds {
		temp[item] = item
	}

	var bind []int
	for _, item := range temp {
		bind = append(bind, item)
	}

	if len(bind) <= 0 {
		response.Fail(ctx, "????????????????????????")
		return
	}

	tx := data.Database.Begin()

	role := model.SysRole{
		Name:    request.Name,
		Summary: request.Summary,
	}

	if tx.Create(&role); role.Id <= 0 {
		tx.Rollback()
		response.Fail(ctx, "????????????")
		return
	}

	var binds []model.SysRoleBindPermission

	for _, item := range bind {
		binds = append(binds, model.SysRoleBindPermission{
			RoleId:       role.Id,
			PermissionId: item,
		})
	}

	if t := tx.CreateInBatches(binds, 100); t.RowsAffected <= 0 {
		tx.Rollback()
		response.Fail(ctx, "????????????")
		return
	}

	var permissions []model.SysRoleBindPermission
	tx.
		Preload("Permission",
			data.Database.
				Where("`method`<>? and `path`<>?", "", ""),
		).
		Where("`role_id` = ?", role.Id).
		Find(&permissions)

	if len(permissions) > 0 {
		var items = make([][]string, len(permissions))
		for idx, item := range permissions {
			items[idx] = []string{item.Permission.Method, item.Permission.Path}
		}

		if _, err := authorize.Casbin.AddPermissionsForUser(authorize.NameByRole(role.Id), items...); err != nil {
			tx.Rollback()
			response.Fail(ctx, "????????????")
			return
		}
	}

	tx.Commit()

	response.Success(ctx)
}

func DoRoleByUpdate(ctx *gin.Context) {

	id, _ := strconv.Atoi(ctx.Param("id"))
	if id <= 0 {
		response.FailByRequestWithMessage(ctx, "ID?????????")
		return
	}

	if id == authorize.ROOT {
		response.Fail(ctx, "??????????????????????????????")
		return
	}

	var request auth.DoRoleByUpdate
	if err := ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	var role model.SysRole

	data.Database.Find(&role, id)
	if role.Id <= 0 {
		response.NotFound(ctx, "???????????????")
		return
	}

	var permissionsIds []int

	var modules, children1, children2, children3 []int

	for _, item := range request.Permissions {
		if len(item) >= 4 {
			children3 = append(children3, item[3])
		} else if len(item) >= 3 {
			children2 = append(children2, item[2])
		} else if len(item) >= 2 {
			children1 = append(children1, item[1])
		} else if len(item) >= 1 {
			modules = append(modules, item[0])
		}
	}

	if len(modules) > 0 {
		var permissions []model.SysPermission
		data.Database.Where("`module_id` in (?) and `method`<>? and `path`<>?", modules, "", "").Find(&permissions)
		for _, item := range permissions {
			permissionsIds = append(permissionsIds, item.Id)
		}
	}
	if len(children3) > 0 {
		var permissions []model.SysPermission
		data.Database.Where("`id` in (?) and `method`<>? and `path`<>?", children3, "", "").Find(&permissions)
		for _, item := range permissions {
			permissionsIds = append(permissionsIds, item.Id)
		}
	}
	if len(children2) > 0 {
		var permissions []model.SysPermission
		data.Database.Where("`parent_i2` in (?) and `method`<>? and `path`<>?", children2, "", "").Find(&permissions)
		for _, item := range permissions {
			permissionsIds = append(permissionsIds, item.Id)
		}
	}
	if len(children1) > 0 {
		var permissions []model.SysPermission
		data.Database.Where("`parent_i1` in (?) and `method`<>? and `path`<>?", children1, "", "").Find(&permissions)
		for _, item := range permissions {
			permissionsIds = append(permissionsIds, item.Id)
		}
	}

	var temp = make(map[int]int, len(permissionsIds))
	for _, item := range permissionsIds {
		temp[item] = item
	}

	var bind []int
	for _, item := range temp {
		bind = append(bind, item)
	}

	if len(bind) <= 0 {
		response.Fail(ctx, "????????????????????????")
		return
	}

	role.Name = request.Name
	role.Summary = request.Summary

	data.Database.Where("`role_id`=?", role.Id).Find(&role.BindPermissions)

	var creates []model.SysRoleBindPermission
	var deletes []int

	if len(role.BindPermissions) > 0 {
		for _, item := range bind {
			mark := true
			for _, value := range role.BindPermissions {
				if item == value.PermissionId {
					mark = false
					break
				}
			}
			if mark {
				creates = append(creates, model.SysRoleBindPermission{
					RoleId:       role.Id,
					PermissionId: item,
				})
			}
		}
		for _, item := range role.BindPermissions {
			mark := true
			for _, value := range bind {
				if item.PermissionId == value {
					mark = false
					break
				}
			}
			if mark {
				deletes = append(deletes, item.Id)
			}
		}
	}

	tx := data.Database.Begin()

	if t := tx.Save(&role); t.RowsAffected <= 0 {
		tx.Rollback()
		response.Fail(ctx, "????????????")
		return
	}

	if len(creates) > 0 {

		if t := tx.CreateInBatches(creates, 100); t.RowsAffected <= 0 {
			tx.Rollback()
			response.Fail(ctx, "????????????")
			return
		}

		var ids []int
		for _, item := range creates {
			ids = append(ids, item.PermissionId)
		}
		var permissions []model.SysPermission
		tx.Where("method<>? and path<>?", "", "").Find(&permissions, ids)
		if len(permissions) > 0 {
			var items [][]string
			for _, item := range permissions {
				items = append(items, []string{item.Method, item.Path})
			}
			if _, err := authorize.Casbin.AddPermissionsForUser(authorize.NameByRole(role.Id), items...); err != nil {
				tx.Rollback()
				response.Fail(ctx, "????????????")
				return
			}
		}
	}

	if len(deletes) > 0 {
		var b model.SysRoleBindPermission
		if t := tx.Where("`role_id`=?", role.Id).Delete(&b, deletes); t.RowsAffected <= 0 {
			tx.Rollback()
			response.Fail(ctx, "????????????")
			return
		}
	}

	if len(deletes) > 0 {
		if _, err := authorize.Casbin.DeletePermissionsForUser(authorize.NameByRole(role.Id)); err != nil {
			tx.Rollback()
			response.Fail(ctx, "????????????")
			return
		}
	}

	if len(creates) > 0 || len(deletes) > 0 {
		var permissions []model.SysRoleBindPermission
		tx.
			Preload("Permission",
				data.Database.
					Where("method <> ? and path <> ?", "", ""),
			).
			Where("`role_id` = ?", role.Id).
			Find(&permissions)

		if len(permissions) > 0 {
			var items = make([][]string, len(permissions))
			for idx, item := range permissions {
				items[idx] = []string{item.Permission.Method, item.Permission.Path}
			}

			if _, err := authorize.Casbin.AddPermissionsForUser(authorize.NameByRole(role.Id), items...); err != nil {
				tx.Rollback()
				response.Fail(ctx, "????????????")
				return
			}
		}
	}

	tx.Commit()

	response.Success(ctx)
}

func DoRoleByDelete(ctx *gin.Context) {

	id, _ := strconv.Atoi(ctx.Param("id"))
	if id <= 0 {
		response.FailByRequestWithMessage(ctx, "ID?????????")
		return
	}

	if id == authorize.ROOT {
		response.Fail(ctx, "??????????????????????????????")
		return
	}

	var role model.SysRole
	data.Database.Find(&role, id)
	if role.Id <= 0 {
		response.NotFound(ctx, "???????????????")
		return
	}

	tx := data.Database.Begin()

	if t := data.Database.Delete(&role); t.RowsAffected <= 0 {
		tx.Rollback()
		response.Fail(ctx, "????????????")
		return
	}

	bind := model.SysRoleBindPermission{RoleId: role.Id}

	if t := tx.Where("`role_id`=?", role.Id).Delete(&bind); t.RowsAffected <= 0 {
		tx.Rollback()
		response.Fail(ctx, "????????????")
		return
	}

	if _, err := authorize.Casbin.DeleteRole(authorize.NameByRole(role.Id)); err != nil {
		tx.Rollback()
		response.Fail(ctx, "????????????")
		return
	}

	tx.Commit()

	response.Success(ctx)
}

func ToRoleByPaginate(ctx *gin.Context) {

	var request auth.ToRoleByPaginate
	if err := ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	tx := data.Database.Where("`id`<>?", authorize.ROOT)

	responses := response.Paginate{
		Page: request.GetPage(),
		Size: request.GetSize(),
		Data: make([]any, 0),
	}

	tc := tx

	tc.Model(model.SysRole{}).Count(&responses.Total)

	if responses.Total > 0 {

		tx = tx.Order("`id` desc")

		var roles []model.SysRole

		tx.Preload("BindPermissions.Permission").Offset(request.GetOffset()).Limit(request.GetLimit()).Find(&roles)

		for _, item := range roles {

			items := authResponse.ToRoleByPaginate{
				Id:        item.Id,
				Name:      item.Name,
				Summary:   item.Summary,
				CreatedAt: item.CreatedAt.ToDateTimeString(),
			}

			for _, value := range item.BindPermissions {
				items.Permissions = append(items.Permissions, []int{value.Permission.ModuleId, value.Permission.ParentI1, value.Permission.ParentI2, value.PermissionId})
			}

			responses.Data = append(responses.Data, items)
		}
	}

	response.SuccessByPaginate(ctx, responses)
}

func ToRoleByEnable(ctx *gin.Context) {

	var roles []model.SysRole

	tx := data.Database

	if !authorize.Root(authorize.Id(ctx)) {
		tx.Where("`role_id`<>?", authorize.ROOT)
	}

	tx.Find(&roles)

	responses := make([]any, 0)

	for _, item := range roles {
		responses = append(responses, authResponse.ToRoleByEnable{
			Id:   item.Id,
			Name: item.Name,
		})
	}

	response.SuccessByList(ctx, responses)
}
