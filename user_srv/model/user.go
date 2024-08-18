package model

import (
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	ID        int32     `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"column:add_time"`
	UpdatedAt time.Time `gorm:"column:update_time"`
	DeletedAt gorm.DeletedAt
	IsDeleted bool
}

type User struct {
	BaseModel
	Mobile   string     `gorm:"index:index_mobile;unique;type:varchar(11);not null"`
	Password string     `gorm:"type:varchar(100);not null"` // 密码加密
	NickName string     `gorm:"type:varchar(20);not null"`
	Birthday *time.Time `gorm:"type:datetime"`                                                        // 生日
	Gender   string     `gorm:"column:gender;default:male;type:varchar(6) comment 'male：男，female：女'"` // "男" "女"
	Role     int32      `gorm:"column:role;default:1;type:int(11) comment '1:普通用户，2：管理员'"`
}
