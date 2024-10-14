package model

import (
	"encoding/json"
	"health-server/internal/db"
	"time"
)

type Additive struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`   // 自增ID
	Name      string    `gorm:"size:100;unique;not null" json:"name"` // 唯一名称
	Desc      string    `gorm:"size:511" json:"desc"`                 // 描述，使用 BLOB
	GB        string    `gorm:"size:50" json:"gb"`                    // GB 标准号
	Category  string    `gorm:"size:100" json:"category"`             // 分类
	Tags      []byte    `gorm:"type:blob" json:"tags"`                // 标签，使用 BLOB
	ImageURL  string    `gorm:"size:255" json:"image_url"`            // 图片 URL
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`     // 创建时间
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`     // 更新时间
}

func (Additive) TableName() string {
	return "additive"
}

func GetAllAdditives() ([]*Additive, error) {
	var additives []*Additive
	err := db.DB.Find(&additives).Error
	return additives, err
}

// AddTags 根据id添加tags
func AddTags(id uint64, tags []int) error {
	tagsBytes, err := json.Marshal(tags)
	if err != nil {
		return err
	}
	return db.DB.Model(&Additive{}).Where("id = ?", id).Update("tags", tagsBytes).Error
}
