package model

import (
	"health-server/internal/db"
	"time"
)

type AppMessage struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`   // 自增ID
	Name      string    `gorm:"size:100;unique;not null" json:"name"` // 唯一名称
	Message   string    `gorm:"size:1000;not null" json:"color"`      // 描述
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`     // 创建时间
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`     // 更新时间
}

func (AppMessage) TableName() string {
	return "app_message"
}

func GetAllAppMessages() ([]*AppMessage, error) {
	var tags []*AppMessage
	err := db.DB.Find(&tags).Error
	return tags, err
}
