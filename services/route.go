package services

import (
	"rbac/global"
	"rbac/model"
	"rbac/model/request"
	"rbac/utils"
)

type RouteService struct {
}

func (u *RouteService) Create(store request.RouteStore) (uint32, error) {
	defer utils.ErrorException()
	route := model.Route{
		Name:   store.Name,
		Method: store.Method,
		Path:   store.Path,
	}

	err := global.Db.Create(&route).Error
	if err != nil {
		return 0, err
	}

	return route.ID, err
}

func (u *RouteService) Update(id uint32, update request.RouteUpdate) (err error) {
	defer utils.ErrorException()

	var route model.Route
	err = global.Db.Model(&model.Route{}).First(&route, "id = ?", id).Error

	if err != nil {
		return
	}

	route.Name = update.Name
	route.Method = update.Method
	route.Path = update.Path

	err = global.Db.Save(&route).Error

	return
}

func (u *RouteService) Destroy(id uint32) (err error) {
	defer utils.ErrorException()

	err = global.Db.Delete(&model.Route{}, id).Error
	return
}

func (u *RouteService) RouteList(list request.RouteList, page int, pageSize int) (routes []model.Route, total int64, err error) {
	defer utils.ErrorException()

	offset := (page - 1) * pageSize
	query := global.Db.Model(&model.Route{})
	if list.Name != "" {
		query = query.Where("name like ?", list.Name+"%")
	}

	err = query.Count(&total).Error
	err = query.Order("id desc").Offset(offset).Limit(pageSize).Find(&routes).Error
	return
}
