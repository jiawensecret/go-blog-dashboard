package model

import (
	"time"
)

type Route struct {
	ID        uint32    `gorm:"primary_key;auto_increment;comment:主键id" json:"id"`
	Name      string    `gorm:"size:50;not null;default:'';comment:权限名称" json:"name"`
	Method    string    `gorm:"size:25;not null;default:'';comment:请求方式" json:"method"`
	Path      string    `gorm:"size:50;not null;default:'';comment:地址" json:"path"`
	CreatedAt time.Time `gorm:"index;comment:创建时间" json:"created_at"`
	UpdatedAt time.Time `gorm:"comment:更新时间" json:"updated_at"`
}
