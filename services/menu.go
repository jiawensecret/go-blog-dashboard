package services

import (
	"rbac/global"
	"rbac/model"
	"rbac/model/request"
	"rbac/utils"
	"strconv"
)

type MenuService struct {
}

type MenuTreeItem struct {
	model.Menu
	Children []*MenuTreeItem `json:"children"`
}

func MenuList(withPermission bool) (trees []MenuTreeItem, err error) {
	defer utils.ErrorException()

	var lists []model.Menu
	query := global.Db.Model(&model.Menu{})
	if withPermission {
		query.Preload("Permissions")
	}
	err = query.Order("menu_sort desc").Find(&lists).Error

	trees = BuiltMenuTree(lists)

	return
}

func (u *MenuService) Create(store request.MenuStore) (uint32, error) {
	defer utils.ErrorException()
	menu := model.Menu{
		Code:        store.Code,
		Name:        store.Name,
		MenuSort:    store.MenuSort,
		MenuStatus:  store.MenuStatus,
		Route:       store.Route,
		Description: store.Description,
		ParentId:    store.ParentId,
	}

	err := global.Db.Create(&menu).Error
	if err != nil {
		return 0, err
	}

	return menu.ID, err
}

func (u *MenuService) Update(id uint32, update request.MenuUpdate) (err error) {
	defer utils.ErrorException()

	var menu model.Menu
	err = global.Db.Model(&model.Menu{}).First(&menu, "id = ?", id).Error

	if err != nil {
		return
	}

	if menu.ParentId != update.ParentId {
		var m model.Menu
		global.Db.Model(&model.Menu{}).First(&m, update.ParentId)
		menu.Flag = m.Flag + "_" + strconv.Itoa(int(menu.ID))

		global.Db.Model(&model.Permission{}).Where("menu_id = ?", menu.ID).Update("flag", menu.Flag)
	}

	menu.Name = update.Name
	menu.Code = update.Code
	menu.MenuSort = update.MenuSort
	menu.MenuStatus = update.MenuStatus
	menu.Route = update.Route
	menu.Description = update.Description
	menu.ParentId = update.ParentId

	err = global.Db.Save(&menu).Error

	return
}

func (u *MenuService) Destroy(id uint32) (err error) {
	defer utils.ErrorException()

	err = global.Db.Delete(&model.Menu{}, id).Error
	return
}

func BuiltMenuTree(list []model.Menu) (trees []MenuTreeItem) {
	dataMap := make(map[uint32][]*MenuTreeItem, len(list))

	var treeList []MenuTreeItem
	for key := range list {
		treeList = append(treeList, MenuTreeItem{
			Menu: list[key],
		})
	}

	for key := range treeList {
		pid := treeList[key].ParentId
		dataMap[pid] = append(dataMap[pid], &treeList[key])
	}

	for key := range treeList {
		treeList[key].Children = dataMap[treeList[key].ID]
	}

	for key := range treeList {
		if treeList[key].ParentId == 0 {
			trees = append(trees, treeList[key])
		}
	}
	return
}

func (u *MenuService) GetChildrenMenu(id uint32) (menus []model.Menu, err error) {
	defer utils.ErrorException()

	query := global.Db.Model(&model.Menu{})

	err = query.Where("parent_id = ?", id).Find(&menus).Error

	return
}
