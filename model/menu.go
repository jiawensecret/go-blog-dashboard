package model

import (
	"gorm.io/gorm"
	"strconv"
	"time"
)

type Menu struct {
	ID          uint32       `gorm:"primary_key;auto_increment;comment:主键id" json:"id"`
	Name        string       `gorm:"size:50;not null;default:'';comment:菜单名称" json:"name"`
	Code        string       `gorm:"size:50;not null;default:'';comment:编码" json:"code"`
	MenuSort    uint8        `gorm:"not null;default:0;comment:排序" json:"menu_sort"`
	MenuStatus  uint8        `gorm:"not null;default:0;comment:状态 0禁用 1启用" json:"menu_status"`
	Route       string       `gorm:"size:50;not null;default:'';comment:路由" json:"route"`
	Description string       `gorm:"size:100;not null;default:'';comment:描述" json:"description"`
	ParentId    uint32       `gorm:"not null;default:0;comment:上级部门id;index" json:"parent_id"`
	Flag        string       `gorm:"size:100;not null;default:'';comment:标志" json:"flag"`
	Permissions []Permission `json:"permissions"`
	CreatedAt   time.Time    `gorm:"index;comment:创建时间" json:"created_at"`
	UpdatedAt   time.Time    `gorm:"comment:更新时间" json:"updated_at"`
}

func (m *Menu) BeforeDelete(tx *gorm.DB) (err error) {
	err = tx.Model(m).Association("Permissions").Clear()
	return
}

func (m *Menu) AfterCreate(tx *gorm.DB) (err error) {
	flag := ""
	if m.ParentId == 0 {
		flag = strconv.Itoa(int(m.ID))
	} else {
		var menu Menu
		tx.Model(m).First(&menu, m.ParentId)
		flag = menu.Flag + "_" + strconv.Itoa(int(m.ID))
	}
	tx.Model(m).Update("flag", flag)
	return
}
