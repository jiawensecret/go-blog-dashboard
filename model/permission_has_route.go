package model

import (
	"time"
)

type PermissionHasRoute struct {
	ID           uint32    `gorm:"primary_key;auto_increment;comment:主键id" json:"id"`
	PermissionId uint32    `gorm:"not null;default:0;comment:权限id" json:"permission_id"`
	RouteId      uint32    `gorm:"not null;default:0;comment:路由id" json:"route_id"`
	CreatedAt    time.Time `gorm:"index;comment:创建时间" json:"created_at"`
	UpdatedAt    time.Time `gorm:"comment:更新时间" json:"updated_at"`
}
