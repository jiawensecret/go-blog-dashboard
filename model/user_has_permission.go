package model

import (
	"time"
)

type UserHasPermission struct {
	ID           uint32    `gorm:"primary_key;auto_increment;comment:主键id" json:"id"`
	UserId       uint32    `gorm:"not null;default:0;comment:用户id" json:"user_id"`
	PermissionId uint32    `gorm:"not null;default:0;comment:部门id" json:"permission_id"`
	CreatedAt    time.Time `gorm:"index;comment:创建时间" json:"created_at"`
	UpdatedAt    time.Time `gorm:"comment:更新时间" json:"updated_at"`
}
