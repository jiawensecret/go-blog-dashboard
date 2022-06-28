package services

import (
	"errors"
	"rbac/global"
	"rbac/model"
	"rbac/model/request"
	"rbac/utils"
)

type PermissionService struct {
}

func (u *PermissionService) Create(store request.PermissionStore) (uint32, error) {
	defer utils.ErrorException()
	permission := model.Permission{
		Name:        store.Name,
		Code:        store.Code,
		MenuId:      store.MenuId,
		Flag:        store.Flag,
		Description: store.Description,
	}

	err := global.Db.Create(&permission).Error
	if err != nil {
		return 0, err
	}

	return permission.ID, err
}

func (u *PermissionService) Update(id uint32, update request.PermissionUpdate) (err error) {
	defer utils.ErrorException()

	var permission model.Permission
	err = global.Db.Model(&model.Permission{}).First(&permission, "id = ?", id).Error

	if err != nil {
		return
	}

	permission.Name = update.Name
	permission.Code = update.Code
	permission.MenuId = update.MenuId
	permission.Flag = update.Flag
	permission.Description = update.Description

	err = global.Db.Save(&permission).Error

	return
}

func (u *PermissionService) Destroy(id uint32) (err error) {
	defer utils.ErrorException()

	err = global.Db.Delete(&model.Permission{}, id).Error
	return
}

func (u *PermissionService) GetById(id uint32) (permission model.Permission, err error) {
	defer utils.ErrorException()
	err = global.Db.Model(&model.Permission{}).Preload("Routes").First(&permission, "id = ?", id).Error
	return
}

func (u *PermissionService) PermissionList(id uint64, list request.PermissionList) (permissions []model.Permission, err error) {
	defer utils.ErrorException()

	query := global.Db.Model(&model.Permission{})
	if list.Name != "" {
		query = query.Where("name like ?", list.Name+"%")
	}
	if list.Code != "" {
		query = query.Where("code = ?", list.Code)
	}

	err = query.Where("menu_id = ?", id).Order("id desc").Find(&permissions).Error
	return
}

func (u *PermissionService) AssRoute(id uint32, ass request.PermissionAssRoute) (err error) {
	var permission model.Permission
	err = global.Db.Model(&model.Permission{}).First(&permission, id).Error
	if err != nil {
		return errors.New("此权限不存在")
	}

	err = global.Db.Model(&permission).Association("Routes").Clear()
	if len(ass.RouteIds) > 0 {
		var routes []model.Route
		err = global.Db.Model(&model.Route{}).Where("id IN ?", ass.RouteIds).Find(&routes).Error
		err = global.Db.Model(&permission).Association("Routes").Replace(routes)
	}

	return
}
