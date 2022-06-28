package services

import (
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"rbac/global"
	"rbac/model"
	"rbac/model/request"
	"rbac/model/user_center"
	"rbac/utils"
	"strconv"
	"strings"
)

type UserService struct {
}

func (u *UserService) List(userList request.UserList, page int, pageSize int) (users []model.User, total int64, err error) {
	defer utils.ErrorException()

	offset := (page - 1) * pageSize
	query := global.Db.Model(&model.User{})
	if userList.Username != "" {
		query = query.Where("username like ?", userList.Username+"%")
	}
	if userList.Name != "" {
		query = query.Where("name like ?", userList.Name+"%")
	}
	if userList.Tel != "" {
		query = query.Where("tel like ?", userList.Tel+"%")
	}
	if userList.Email != "" {
		query = query.Where("email like ?", userList.Email+"%")
	}
	if userList.Position != "" {
		query = query.Where("position = ?", userList.Position)
	}
	if userList.UserStatus != 0 {
		query = query.Where("user_status = ?", userList.UserStatus)
	}
	var ids1 []uint32
	var ids2 []uint32
	if len(userList.Departments) > 0 {
		_ = global.Db.Model(&model.UserHasDepartment{}).Select("user_id").Where("department_id IN ?", userList.Departments).Find(&ids1).Error
		if len(ids1) == 0 {
			return
		}
	}
	if len(userList.Permissions) > 0 {
		_ = global.Db.Model(&model.UserHasPermission{}).Select("user_id").Where("permission_id IN ?", userList.Permissions).Find(&ids2).Error
		if len(ids2) == 0 {
			return
		}
	}

	ids := utils.Intersect(ids1, ids2)
	if len(userList.Ids) > 0 {
		ids = utils.Intersect(ids, userList.Ids)
		if len(ids) == 0 {
			return
		}
	}
	if len(ids) > 0 {
		query = query.Where("id IN ?", ids)
	}
	err = query.Count(&total).Error
	err = query.Preload("Departments").Order("id desc").Offset(offset).Limit(pageSize).Find(&users).Error
	return users, total, err
}

func (u *UserService) Create(store request.UserStore) (uint32, error) {
	defer utils.ErrorException()
	user := model.User{
		Username:   store.Username,
		Name:       store.Name,
		Tel:        store.Tel,
		Email:      store.Email,
		UserStatus: 1,
		Position:   store.Position,
	}

	err := global.Db.Create(&user).Error
	if err != nil {
		return 0, err
	}

	if len(store.DepartmentIds) > 0 {
		var departments []model.Department
		err = global.Db.Model(&model.Department{}).Where("id IN ?", store.DepartmentIds).Find(&departments).Error
		err = global.Db.Model(&user).Association("Departments").Replace(departments)
	}
	go u.addUserToAlertCenter(user.ID, store.Name, store.Email, store.Tel)

	return user.ID, err
}

func (u *UserService) GetById(id uint32) (user model.User, err error) {
	err = global.Db.Model(&model.User{}).Preload("Departments").First(&user, id).Error
	return
}

func (u *UserService) Update(id uint32, update request.UserUpdate) (err error) {
	defer utils.ErrorException()

	var user model.User
	err = global.Db.Model(&model.User{}).First(&user, "id = ?", id).Error

	if err != nil {
		return
	}

	user.Name = update.Name
	user.Username = update.Username
	user.Tel = update.Tel
	user.Email = update.Email
	user.Position = update.Position

	err = global.Db.Save(&user).Error

	if len(update.DepartmentIds) > 0 {
		err = global.Db.Model(&user).Association("Departments").Clear()
		var departments []model.Department
		err = global.Db.Model(&model.Department{}).Where("id IN ?", update.DepartmentIds).Find(&departments).Error
		err = global.Db.Model(&user).Association("Departments").Replace(departments)
	}

	go u.updateUserToAlertCenter(strconv.Itoa(int(id)), update.Name, update.Tel, update.Email)
	return
}

func (u *UserService) UpdateByMobile(mobile string, update request.UserUpdate) (err error) {
	defer utils.ErrorException()

	var user model.User
	err = global.Db.Model(&model.User{}).First(&user, "tel = ?", mobile).Error

	if err != nil {
		return
	}

	user.Name = update.Name
	user.Username = update.Username
	user.Tel = update.Tel
	user.Email = update.Email
	user.Position = update.Position

	err = global.Db.Save(&user).Error

	if len(update.DepartmentIds) > 0 {
		err = global.Db.Model(&user).Association("Departments").Clear()
		var departments []model.Department
		err = global.Db.Model(&model.Department{}).Where("id IN ?", update.DepartmentIds).Find(&departments).Error
		err = global.Db.Model(&user).Association("Departments").Replace(departments)
	}

	go u.updateUserToAlertCenter(strconv.Itoa(int(user.ID)), update.Name, update.Tel, update.Email)

	return
}

func (u *UserService) ChangeStatus(id uint32, update request.UserChangeStatus) (err error) {
	defer utils.ErrorException()
	var user model.User
	err = global.Db.Model(&model.User{}).First(&user, "id = ?", id).Error

	if err != nil {
		return
	}

	user.UserStatus = update.UserStatus
	err = global.Db.Save(&user).Error

	go u.updateUserStatusToAlertCenter(strconv.Itoa(int(id)), strconv.Itoa(int(user.UserStatus)))
	return
}

func (u *UserService) GetUserMenu(id uint32) (trees []MenuTreeItem, err error) {
	defer utils.ErrorException()
	var user model.User
	err = global.Db.Model(&model.User{}).Preload("Permissions").First(&user, "id = ?", id).Error
	if err != nil {
		return
	}
	menuIds := make([]int, 0)
	permissionIds := make([]uint32, 0)
	if len(user.Permissions) > 0 {
		for _, permission := range user.Permissions {
			flag := strings.Split(permission.Flag, "_")
			for _, menuId := range flag {
				id, _ := strconv.Atoi(menuId)
				if !findInIntArray(id, menuIds) {
					menuIds = append(menuIds, id)
				}
			}
			permissionIds = append(permissionIds, permission.ID)
		}
	}

	var menus []model.Menu
	err = global.Db.Model(&model.Menu{}).Preload("Permissions", "id IN ?", permissionIds).Where("id IN ?", menuIds).Where("menu_status = ?", 1).Order("menu_sort desc").Find(&menus).Error

	trees = BuiltMenuTree(menus)
	return
}

func (u *UserService) GetUserPermission(id uint32) (trees []MenuTreeItem, err error) {
	defer utils.ErrorException()

	var user model.User
	err = global.Db.Model(&model.User{}).Preload("Permissions").First(&user, "id = ?", id).Error
	if err != nil {
		return
	}

	permissionIds := make([]uint32, 0)
	if len(user.Permissions) > 0 {
		for _, permission := range user.Permissions {
			permissionIds = append(permissionIds, permission.ID)
		}
	}

	trees, err = MenuList(true)
	for _, menu := range trees {
		if len(menu.Children) > 0 {
			changePermissionChecked(menu.Children, permissionIds)
		} else if len(menu.Permissions) > 0 {
			for key := range menu.Permissions {
				menu.Permissions[key].Checked = findInPermissionArray(menu.Permissions[key].ID, permissionIds)
			}
		}
	}

	return
}

func (u *UserService) BatchDeletePermission(delete request.UserPermissionDelete) (err error) {
	defer utils.ErrorException()

	if len(delete.PermissionIds) == 0 {
		return errors.New("提交的权限不能为空")
	}
	if len(delete.UserIds) == 0 {
		return errors.New("用户不能为空")
	}

	err = global.Db.Where("user_id IN ?", delete.UserIds).Where("permission_id IN ?", delete.PermissionIds).Delete(&model.UserHasPermission{}).Error
	if err != nil {
		global.Log.Error("批量删除用户权限发生错误", err.Error())
		return errors.New("批量删除用户权限发生错误")
	}

	return
}

func (u *UserService) BatchSavePermission(store request.UserPermissionStore) (err error) {
	defer utils.ErrorException()

	if len(store.PermissionIds) == 0 {
		return errors.New("提交的权限不能为空")
	}

	if len(store.UserIds) == 0 {
		return errors.New("用户不能为空")
	}

	for _, userId := range store.UserIds {
		global.Db.Where("user_id = ?", userId).Where("permission_id in ?", store.PermissionIds).Delete(&model.UserHasPermission{})
		var data []model.UserHasPermission
		for _, permissionId := range store.PermissionIds {
			data = append(data, model.UserHasPermission{
				UserId:       userId,
				PermissionId: permissionId,
			})
		}
		err = global.Db.Create(&data).Error
		if err != nil {
			return err
		}
		if len(store.DataPermissions) > 0 {
			var menuIds []uint32
			var dataPermissions []model.DataPermission
			for _, dataPermission := range store.DataPermissions {
				menuIds = append(menuIds, dataPermission.MenuId)
				dataPermissions = append(dataPermissions, model.DataPermission{
					UserId:      userId,
					MenuId:      dataPermission.MenuId,
					ExpiredTime: dataPermission.ExpiredTime,
					DataLevel:   uint8(dataPermission.DataLevel),
				})
			}
			global.Db.Where("user_id = ?", userId).Where("menu_id in ?", menuIds).Delete(&model.DataPermission{})
			err = global.Db.Create(&dataPermissions).Error
			if err != nil {
				return err
			}
		}
	}

	return
}

func (u *UserService) SaveUserPermissions(id uint32, update request.UserPermissionUpdate) (err error) {
	defer utils.ErrorException()

	var user model.User
	err = global.Db.Model(&model.User{}).First(&user, "id = ?", id).Error
	err = global.Db.Model(&user).Association("Permissions").Clear()
	global.Db.Where("user_id = ?", user.ID).Delete(&model.DataPermission{})

	if len(update.PermissionIds) > 0 {
		var data []model.UserHasPermission
		for _, permissionId := range update.PermissionIds {
			data = append(data, model.UserHasPermission{
				UserId:       user.ID,
				PermissionId: permissionId,
			})
		}
		err = global.Db.Create(&data).Error
	}

	if len(update.DataPermissions) > 0 {
		var data []model.DataPermission
		for _, dataPermission := range update.DataPermissions {
			m := model.DataPermission{
				UserId:      user.ID,
				MenuId:      dataPermission.MenuId,
				ExpiredTime: dataPermission.ExpiredTime,
				DataLevel:   uint8(dataPermission.DataLevel),
			}
			data = append(data, m)
		}
		err = global.Db.Create(&data).Error
	}

	return
}

func (u *UserService) GetDataPermission(userId uint32) (dataPermissions []model.DataPermission, err error) {
	defer utils.ErrorException()
	err = global.Db.Model(&model.DataPermission{}).Where("user_id = ?", userId).Find(&dataPermissions).Error
	return
}

func (u *UserService) GetUserByMobile(mobile string) (user model.User, err error) {
	defer utils.ErrorException()
	err = global.Db.Model(&model.User{}).Preload("Departments").Preload("Permissions").Preload("Permissions.Routes").Preload("DataPermission").Where("tel = ?", mobile).Find(&user).Error
	return
}

func (u *UserService) GetUserPositions() (positions []string, err error) {
	defer utils.ErrorException()
	err = global.Db.Model(&model.User{}).Distinct("position").Select("position").Find(&positions).Error
	return
}

func (u *UserService) GetUserFromUserCenter(mobiles []string) (user []user_center.User, err error) {
	defer utils.ErrorException()
	mobile := strings.Join(mobiles, ",")
	url := global.Config.Service.UserCenter + "/api/open/v1/user/employee-list"
	params := make(map[string]string)
	agent := fiber.Post(url)
	args := fiber.AcquireArgs()
	args.Set("mobile", mobile)
	args.Set("app_key", "ops")
	params["mobile"] = mobile
	params["app_key"] = "ops"
	args.Set("sn", utils.Sign(params, global.Config.Service.UserCenterSecret))

	agent.Form(args)
	if err = agent.Parse(); err != nil {
		return
	}
	_, body, _ := agent.Bytes()

	var result user_center.UserInfoApi
	err = json.Unmarshal(body, &result)
	if err != nil {
		return user, errors.New("用户中心未查到此用户")
	}

	user = result.Data
	return
}

func (u *UserService) SyncAllDepartments() (err error) {
	defer utils.ErrorException()
	url := global.Config.Service.UserCenter + "/api/open/v1/user/depart-all"
	params := make(map[string]string)
	agent := fiber.Post(url)
	args := fiber.AcquireArgs()
	args.Set("app_key", "ops")
	params["app_key"] = "ops"
	args.Set("sn", utils.Sign(params, global.Config.Service.UserCenterSecret))

	agent.Form(args)
	if err = agent.Parse(); err != nil {
		return
	}
	_, body, _ := agent.Bytes()

	var result user_center.DepartmentsApi
	err = json.Unmarshal(body, &result)

	if len(result.Data) > 0 {
		global.Db.Where("1 = 1").Delete(&model.Department{})
		k := make([]model.Department, 0)
		for _, department := range result.Data {
			k = append(k, model.Department{
				ID:       department.Id,
				Name:     department.Name,
				ParentId: department.ParentId,
			})
		}
		global.Db.Create(&k)
	}

	return
}

func (u *UserService) SyncUsers() {
	defer utils.ErrorException()
	err := u.SyncAllDepartments()
	if err != nil {
		return
	}

	var mobiles []string
	err = global.Db.Model(&model.User{}).Select("tel").Find(&mobiles).Error

	users, err := u.GetUserFromUserCenter(mobiles)
	if err != nil || len(users) == 0 {
		return
	}

	for _, user := range users {
		var ids []uint64
		if len(user.Department) > 0 {
			for _, v := range user.Department {
				ids = append(ids, v.ID)
			}
		}

		c := request.UserUpdate{
			Username:      user.Username,
			Name:          user.RealName,
			Tel:           user.Mobile,
			Email:         user.Email,
			Position:      user.Position,
			DepartmentIds: ids,
		}

		go func() {
			err := u.UpdateByMobile(c.Tel, c)
			if err != nil {
				global.Log.Error(c.Tel+"同步数据发生错误", err.Error())
			}
		}()
	}

	return
}

func (u *UserService) addUserToAlertCenter(id uint32, username string, email string, tel string) {
	defer utils.ErrorException()
	url := global.Config.Service.AlertCenter + "/alertcenter/admin/user"
	agent := fiber.Post(url)
	args := fiber.AcquireArgs()
	args.Set("id", strconv.Itoa(int(id)))
	args.Set("username", username)
	args.Set("tel", tel)
	args.Set("email", email)

	agent.Form(args)
	if err := agent.Parse(); err != nil {
		return
	}
	_, _, _ = agent.Bytes()
	return
}

func (u *UserService) updateUserToAlertCenter(id string, username string, tel string, email string) {
	defer utils.ErrorException()
	url := global.Config.Service.AlertCenter + "/alertcenter/admin/user/" + id
	agent := fiber.Put(url)
	args := fiber.AcquireArgs()
	args.Set("username", username)
	args.Set("tel", tel)
	args.Set("email", email)

	agent.Form(args)
	if err := agent.Parse(); err != nil {
		return
	}
	_, _, _ = agent.Bytes()
	return
}
func (u *UserService) updateUserStatusToAlertCenter(id string, user_status string) {
	defer utils.ErrorException()
	url := global.Config.Service.AlertCenter + "/alertcenter/admin/user/" + id + "/status"
	agent := fiber.Post(url)
	args := fiber.AcquireArgs()
	args.Set("user_status", user_status)

	agent.Form(args)
	if err := agent.Parse(); err != nil {
		return
	}
	_, _, _ = agent.Bytes()
	return
}

func changePermissionChecked(menus []*MenuTreeItem, permissionIds []uint32) {
	for _, menu := range menus {
		if len(menu.Children) > 0 {
			changePermissionChecked(menu.Children, permissionIds)
		} else if len(menu.Permissions) > 0 {
			for key := range menu.Permissions {
				menu.Permissions[key].Checked = findInPermissionArray(menu.Permissions[key].ID, permissionIds)
			}
		}
	}
}

func findInIntArray(id int, arr []int) bool {
	for _, p := range arr {
		if p == id {
			return true
		}
	}
	return false
}

func findInPermissionArray(id uint32, arr []uint32) bool {
	for _, p := range arr {
		if p == id {
			return true
		}
	}
	return false
}
