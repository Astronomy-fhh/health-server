package model

import (
	"errors"
	"gorm.io/gorm"
	"health-server/internal/db"
	"time"
)

type Admin struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"size:100;not null;unique"` // 用户名，唯一且不能为空
	Pass      string    `gorm:"size:255;not null"`        // 密码
	Avatar    string    `gorm:"size:255"`                 // 自定义头像
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

func (Admin) TableName() string {
	return "admin"
}

// GetAdminByName 根据用户名获取管理员信息
func GetAdminByName(name string) (*Admin, error) {
	var admin Admin
	err := db.DB.Where("name = ?", name).First(&admin).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &admin, nil
}
