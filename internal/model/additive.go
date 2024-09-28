package model

import (
	"health-server/internal/db"
	"time"
)

type Additive struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement"` // 自增ID
	Name      string    `gorm:"size:100;unique;not null"` // 唯一名称
	Desc      []byte    `gorm:"type:blob"`                // 描述，使用 BLOB
	GB        string    `gorm:"size:50"`                  // GB 标准号
	Status    []byte    `gorm:"type:blob"`                // 状态，使用 BLOB
	Category  []byte    `gorm:"type:blob"`                // 分类，使用 BLOB
	Tags      []byte    `gorm:"type:blob"`                // 标签，使用 BLOB
	ImageURL  string    `gorm:"size:255"`                 // 图片 URL
	CreatedAt time.Time `gorm:"autoCreateTime"`           // 创建时间
	UpdatedAt time.Time `gorm:"autoUpdateTime"`           // 更新时间
}

func (Additive) TableName() string {
	return "additive"
}

func GetAllAdditives() ([]*Additive, error) {
	var additives []*Additive
	err := db.DB.Find(&additives).Error
	return additives, err
}
