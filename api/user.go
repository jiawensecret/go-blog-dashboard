package api

import (
	"encoding/json"
	"errors"
	"github.com/kataras/iris/v12"
	"platform/model/request"
	"platform/model/response"
	"platform/services"
	"platform/utils"
	"strconv"
)

//@Summary 用户列表
//@Description 获取用户列表
//@Produce json
//@Param ids_str query string false "用户id字符串"
//@Param username query string false "用户名"
//@Param name query string false "名称"
//@Param tel query string false "电话"
//@Param email query string false "邮箱"
//@Param position query string false "职位"
//@Param department_str query string false "部门json字符串"
//@Param permission_str query string false "权限json字符串"
//@Success 200 {object} response.PageResult{list=[]model.User}
//@Router /api/users [get]
func UserList(ctx *iris.Context) error {
	var userList request.UserList
	err := ctx.QueryParser(&userList)
	if err != nil {
		return response.FailWithValidate(err, ctx)
	}

	if ok, err := utils.NewValidator(userList).IsOk(); ok {
		return response.FailWithValidate(err, ctx)
	}

	department := ctx.Query("department_str", "")
	if department != "" {
		var ids []uint64
		err = json.Unmarshal([]byte(department), &ids)
		if err != nil {
			return response.FailWithValidate(err, ctx)
		}
		userList.Departments = ids
	}

	permission := ctx.Query("permission_str", "")
	if permission != "" {
		var permissionIds []uint32
		err = json.Unmarshal([]byte(permission), &permissionIds)
		if err != nil {
			return response.FailWithValidate(err, ctx)
		}
		userList.Permissions = permissionIds
	}

	idsStr := ctx.Query("ids_str", "")
	if idsStr != "" {
		var ids []uint32
		err = json.Unmarshal([]byte(idsStr), &ids)
		if err != nil {
			return response.FailWithValidate(err, ctx)
		}
		userList.Ids = ids
	}

	page, pageSize := request.GetPageInfo(userList.Page, userList.PageSize)

	userService := services.UserService{}
	users, total, err := userService.List(userList, page, pageSize)
	if err != nil {
		return response.FailWithError(err, ctx)
	}

	for key := range users {
		users[key].CreateTime = users[key].CreatedAt.Format("2006-01-02")
		users[key].Tel = users[key].Tel[0:3] + "****" + users[key].Tel[7:]
	}

	return response.SuccessWithData(response.PageResult{
		Data:     users,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, ctx)

}

//@Summary 新增用户接口
//@Description 新增用户
//@Produce json
//@Param username body string true "用户名"
//@Param name body string true "名称"
//@Param tel body string true "电话"
//@Param email body string true "邮箱"
//@Param position body string true "职位"
//@Param department_ids body []uint32 true "部门id"
//@Success 200 {object} response.Response
//@Router /api/user [post]
func UserCreate(ctx *fiber.Ctx) error {
	var userStore request.UserStore
	err := ctx.BodyParser(&userStore)
	if err != nil {
		return response.FailWithValidate(err, ctx)
	}

	if ok, err := utils.NewValidator(userStore).IsOk(); ok {
		return response.FailWithValidate(err, ctx)
	}

	userService := services.UserService{}

	u, err := userService.GetUserByMobile(userStore.Tel)
	if u.ID > 0 {
		return response.FailWithValidate(errors.New("用户已存在"), ctx)
	}
	id, err := userService.Create(userStore)
	if err != nil || id == 0 {
		return response.FailWithError(err, ctx)
	}
	return response.Success(ctx)
}

//@Summary 编辑用户接口
//@Description 编辑用户
//@Produce json
//@Param id path int true "用户id"
//@Param username body string true "用户名"
//@Param name body string true "名称"
//@Param tel body string true "电话"
//@Param email body string true "邮箱"
//@Param position body string true "职位"
//@Param department_ids body []uint32 true "部门id"
//@Success 200 {object} response.Response
//@Router /api/user/:id [put]
func UserUpdate(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id", ""), 10, 64)
	if err != nil || id <= 0 {
		return response.FailWithValidate(errors.New("id不能为空"), ctx)
	}
	var userUpdate request.UserUpdate
	err = ctx.BodyParser(&userUpdate)
	if err != nil {
		return response.FailWithValidate(err, ctx)
	}

	if ok, err := utils.NewValidator(userUpdate).IsOk(); ok {
		return response.FailWithValidate(err, ctx)
	}

	userService := services.UserService{}
	err = userService.Update(uint32(id), userUpdate)
	if err != nil || id == 0 {
		return response.FailWithError(err, ctx)
	}
	return response.Success(ctx)
}

//@Summary 编辑用户接口
//@Description 编辑用户
//@Produce json
//@Success 200 {object} response.Response{data=model.User}
//@Router /api/user/:id [get]
func UserInfo(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id", ""), 10, 64)
	if err != nil || id <= 0 {
		return response.FailWithValidate(errors.New("id不能为空"), ctx)
	}

	userService := services.UserService{}
	user, err := userService.GetById(uint32(id))
	if err != nil || user.ID == 0 {
		return response.FailWithError(err, ctx)
	}
	return response.SuccessWithData(user, ctx)
}

//@Summary 启用禁用接口
//@Description 用户启用禁用
//@Produce json
//@Param id path int true "用户id"
//@Param user_status body int true "用户状态"
//@Success 200 {object} response.Response
//@Router /api/user/:id/status [post]
func UserChangeStatus(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id", ""), 10, 64)
	if err != nil || id <= 0 {
		return response.FailWithValidate(errors.New("id不能为空"), ctx)
	}
	var userUpdate request.UserChangeStatus
	err = ctx.BodyParser(&userUpdate)
	if err != nil {
		return response.FailWithValidate(err, ctx)
	}

	if ok, err := utils.NewValidator(userUpdate).IsOk(); ok {
		return response.FailWithValidate(err, ctx)
	}

	userService := services.UserService{}
	err = userService.ChangeStatus(uint32(id), userUpdate)
	if err != nil {
		return response.FailWithError(err, ctx)
	}
	return response.Success(ctx)
}

//@Summary 获取用户permissions 包含所有权限
//@Description 获取用户permissions 包含所有权限
//@Produce json
//@Param id path int true "用户id"
//@Success 200 {object} response.Response{data=[]services.MenuTreeItem}
//@Router /api/user/:id/permissions [get]
func UserPermission(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id", ""), 10, 64)
	if err != nil || id <= 0 {
		return response.FailWithValidate(errors.New("id不能为空"), ctx)
	}
	userService := services.UserService{}
	trees, err := userService.GetUserPermission(uint32(id))
	return response.SuccessWithData(trees, ctx)
}

//@Summary 获取用户菜单树 包含拥有的权限
//@Description 获取用户菜单树 包含拥有的权限
//@Produce json
//@Param id path int true "用户id"
//@Success 200 {object} response.Response{data=[]services.MenuTreeItem}
//@Router /api/user/:id/menus [get]
func UserMenu(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id", ""), 10, 64)
	if err != nil || id <= 0 {
		return response.FailWithValidate(errors.New("id不能为空"), ctx)
	}

	userService := services.UserService{}
	trees, err := userService.GetUserMenu(uint32(id))

	return response.SuccessWithData(trees, ctx)
}

//@Summary 新增用户权限
//@Description 新增用户权限
//@Produce json
//@Param permission_ids body []uint32 true "权限id"
//@Param user_ids body []uint32 true "用户ids"
//@Success 200 {object} response.Response
//@Router /api/user/permission [post]
func CreateUserPermission(ctx *fiber.Ctx) error {
	var store request.UserPermissionStore
	err := ctx.BodyParser(&store)
	if err != nil {
		return response.FailWithValidate(err, ctx)
	}

	if ok, err := utils.NewValidator(store).IsOk(); ok {
		return response.FailWithValidate(err, ctx)
	}

	userService := services.UserService{}
	err = userService.BatchSavePermission(store)
	if err != nil {
		return response.FailWithError(err, ctx)
	}
	return response.Success(ctx)
}

//@Summary 批量删除用户权限
//@Description 批量删除用户权限
//@Produce json
//@Param permission_ids body []uint32 true "权限id"
//@Param user_ids body []uint32 true "用户ids"
//@Success 200 {object} response.Response
//@Router /api/user/permission/delete [post]
func DeleteUserPermission(ctx *fiber.Ctx) error {
	var d request.UserPermissionDelete
	err := ctx.BodyParser(&d)
	if err != nil {
		return response.FailWithValidate(err, ctx)
	}

	if ok, err := utils.NewValidator(d).IsOk(); ok {
		return response.FailWithValidate(err, ctx)
	}

	userService := services.UserService{}
	err = userService.BatchDeletePermission(d)
	if err != nil {
		return response.FailWithError(err, ctx)
	}
	return response.Success(ctx)
}

//@Summary 编辑用户权限
//@Description 编辑用户权限
//@Produce json
//@Param id path int true "用户id"
//@Param permission_ids body []uint32 true "权限id"
//@Param data_permissions body []request.DataPermission true "数据权限"
//@Success 200 {object} response.Response
//@Router /api/user/:id/permission [put]
func ModifyUserPermission(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id", ""), 10, 64)
	if err != nil || id <= 0 {
		return response.FailWithValidate(errors.New("id不能为空"), ctx)
	}
	var update request.UserPermissionUpdate
	err = ctx.BodyParser(&update)
	if err != nil {
		return response.FailWithValidate(err, ctx)
	}

	if ok, err := utils.NewValidator(update).IsOk(); ok {
		return response.FailWithValidate(err, ctx)
	}

	userService := services.UserService{}
	err = userService.SaveUserPermissions(uint32(id), update)

	if err != nil {
		return response.FailWithError(err, ctx)
	}
	return response.Success(ctx)
}

//@Summary 获取用户数据权限
//@Description 获取用户数据权限
//@Produce json
//@Param id path int true "用户id"
//@Success 200 {object} response.Response
//@Router /api/user/:id/data-permission [get]
func UserDataPermission(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id", ""), 10, 64)
	if err != nil || id <= 0 {
		return response.FailWithValidate(errors.New("id不能为空"), ctx)
	}

	userService := services.UserService{}
	dataPermissions, err := userService.GetDataPermission(uint32(id))
	if err != nil {
		return response.FailWithError(err, ctx)
	}
	return response.SuccessWithData(dataPermissions, ctx)
}

//@Summary 获取用户数据权限
//@Description 获取用户数据权限
//@Produce json
//@Param user body request.UserStore true "user基本信息"
//@Param permission_ids body []uint32 true "权限id"
//@Param data_permissions body []request.DataPermission true "数据权限"
//@Success 200 {object} response.Response
//@Router /api/user/create-with-permission [post]
func UserCreateWithPermission(ctx *fiber.Ctx) error {
	var store request.CopyUser
	err := ctx.BodyParser(&store)
	if err != nil {
		return response.FailWithValidate(err, ctx)
	}

	if ok, err := utils.NewValidator(store).IsOk(); ok {
		return response.FailWithValidate(err, ctx)
	}

	userService := services.UserService{}
	u, err := userService.GetUserByMobile(store.User.Tel)
	if u.ID > 0 {
		return response.FailWithValidate(errors.New("用户已存在"), ctx)
	}
	id, err := userService.Create(store.User)
	if err != nil || id == 0 {
		return response.FailWithError(err, ctx)
	}
	err = userService.SaveUserPermissions(id, request.UserPermissionUpdate{
		PermissionIds:   store.PermissionIds,
		DataPermissions: store.DataPermissions,
	})

	if err != nil {
		return response.FailWithError(err, ctx)
	}
	return response.Success(ctx)
}

//@Summary 获取用户数据
//@Description 获取用户数据
//@Produce json
//@Param mobile query string true "手机号码"
//@Success 200 {object} response.Response
//@Router /api/user/get-user-info [get]
func GetUserInfo(ctx *fiber.Ctx) error {
	var store request.UserInfo
	err := ctx.QueryParser(&store)
	if err != nil {
		return response.FailWithValidate(err, ctx)
	}

	if ok, err := utils.NewValidator(store).IsOk(); ok {
		return response.FailWithValidate(err, ctx)
	}

	userService := services.UserService{}
	user, err := userService.GetUserByMobile(store.Mobile)

	if err != nil {
		return response.FailWithError(err, ctx)
	}
	return response.SuccessWithData(user, ctx)
}

//@Summary 获取职位下拉数据
//@Description 获取职位下拉数据
//@Produce json
//@Success 200 {object} response.Response
//@Router /api/positions [get]
func PositionList(ctx *fiber.Ctx) error {
	userService := services.UserService{}
	positions, err := userService.GetUserPositions()

	if err != nil {
		return response.FailWithError(err, ctx)
	}
	return response.SuccessWithData(positions, ctx)
}

//@Summary 从用户中心获取用户数据
//@Description 从用户中心获取用户数据
//@Produce json
//@Param mobile query string true "手机号码"
//@Success 200 {object} response.Response{data=user_center.User}
//@Router /api/user-center/user [get]
func GetUserInfoWithMobile(ctx *fiber.Ctx) error {
	var store request.UserInfo
	err := ctx.QueryParser(&store)

	if err != nil {
		return response.FailWithValidate(err, ctx)
	}

	if ok, err := utils.NewValidator(store).IsOk(); ok {
		return response.FailWithValidate(err, ctx)
	}

	userService := services.UserService{}
	var mobiles []string
	mobiles = append(mobiles, store.Mobile)
	users, err := userService.GetUserFromUserCenter(mobiles)

	if err != nil {
		return response.FailWithError(err, ctx)
	}
	user := users[0]
	return response.SuccessWithData(user, ctx)
}

//@Summary 从用户中心同步所有用户
//@Description 从用户中心同步所有用户
//@Produce json
//@Success 200 {object} response.Response
//@Router /api/sync-users [get]
func SyncUsers(ctx *fiber.Ctx) error {
	userService := services.UserService{}
	userService.SyncUsers()

	return response.Success(ctx)
}
