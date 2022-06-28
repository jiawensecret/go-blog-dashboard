package model

import (
	"time"
)

type Permission struct {
	ID          uint32    `gorm:"primary_key;auto_increment;comment:主键id" json:"id"`
	Name        string    `gorm:"size:50;not null;default:'';comment:权限名称" json:"name"`
	Code        string    `gorm:"size:50;not null;default:'';comment:编码" json:"code"`
	MenuId      uint32    `gorm:"not null;default:0;comment:菜单id" json:"menu_id"`
	Flag        string    `gorm:"size:100;not null;default:'';comment:标志" json:"flag"`
	Description string    `gorm:"size:100;not null;default:'';comment:描述" json:"description"`
	Routes      []Route   `gorm:"many2many:permission_has_routes;" json:"routes"`
	Checked     bool      `gorm:"-" json:"checked"`
	CreatedAt   time.Time `gorm:"index;comment:创建时间" json:"created_at"`
	UpdatedAt   time.Time `gorm:"comment:更新时间" json:"updated_at"`
}
