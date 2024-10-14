package model

import (
	"health-server/internal/db"
	"time"
)

type AdditiveTag struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`   // 自增ID
	Name      string    `gorm:"size:100;unique;not null" json:"name"` // 唯一名称
	Color     string    `gorm:"size:50;not null" json:"color"`        // GB 标准号
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`     // 创建时间
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`     // 更新时间
}

func (AdditiveTag) TableName() string {
	return "additive_tag"
}

func GetAllAdditiveTags() ([]*AdditiveTag, error) {
	var tags []*AdditiveTag
	err := db.DB.Find(&tags).Error
	return tags, err
}
