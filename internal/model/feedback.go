package model

import (
	"health-server/internal/db"
	"time"
)

type Feedback struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`   // 自增ID
	Desc      string    `gorm:"size:100;unique;not null" json:"desc"` // 唯一描述
	Stats     int       `gorm:"not null" json:"stats"`                // 状态，使用整型
	CreateUid string    `gorm:"size:100;not null" json:"create_uid"`  // 创建者ID
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`     // 创建时间
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`     // 更新时间
}

func (Feedback) TableName() string {
	return "feedback"
}

func AddFeedback(feedback *Feedback) error {
	return db.DB.Create(feedback).Error
}
