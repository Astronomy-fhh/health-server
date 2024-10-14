package model

import (
	"health-server/internal/db"
	"time"
)

type UserDefault struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`   // 自增ID
	Name      string    `gorm:"size:100;unique;not null" json:"name"` // 唯一名称
	Img       string    `gorm:"size:100" json:"img"`                  // 描述，使用 BLOB
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`     // 创建时间
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`     // 更新时间
}

func (UserDefault) TableName() string {
	return "user_default"
}

func GetUserDefaults() ([]*UserDefault, error) {
	var users []*UserDefault
	err := db.DB.Find(&users).Error
	return users, err
}
