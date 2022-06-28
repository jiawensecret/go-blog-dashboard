package model

import (
	"time"
)

type User struct {
	ID          uint32       `gorm:"primary_key;auto_increment;comment:主键id" json:"id"`
	Username    string       `gorm:"size:25;not null;default:'';comment:用户名" json:"username"`
	Name        string       `gorm:"size:25;not null;default:'';comment:姓名" json:"name"`
	Tel         string       `gorm:"size:25;not null;default:'';unique;comment:电话" json:"tel"`
	Email       string       `gorm:"size:100;not null;default:'';unique;comment:邮箱" json:"email"`
	UserStatus  int8         `gorm:"not null;default:0;comment:1正常 0默认 -1禁用" json:"user_status"`
	Position    string       `gorm:"size:100;not null;default:'';comment:职位" json:"position"`
	Permissions []Permission `gorm:"many2many:user_has_permissions;" json:"permissions"`
	CreateTime  string       `gorm:"-" json:"create_time"`
	CreatedAt   time.Time    `gorm:"index;comment:创建时间" json:"created_at"`
	UpdatedAt   time.Time    `gorm:"comment:更新时间" json:"updated_at"`
}
