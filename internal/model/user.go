package model

import (
	"errors"
	"github.com/teris-io/shortid"
	"gorm.io/gorm"
	"health-server/internal/db"
	"time"
)

type User struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement"`
	UID          string    `gorm:"uniqueIndex"`
	DID          string    `gorm:"size:255;column:did"` // 设备ID
	Name         string    `gorm:"size:100"`
	SystemAvatar string    `gorm:"size:255"` // 系统头像
	CustomAvatar string    `gorm:"size:255"` // 自定义头像
	BindID       string    `gorm:"size:100"` // 绑定ID（如电话号码）
	RegisteredAt time.Time `gorm:"type:timestamp"`
	LastLoginAt  time.Time `gorm:"type:timestamp"`
	CreatedAt    time.Time `gorm:"type:timestamp"`
	UpdatedAt    time.Time `gorm:"type:timestamp"`
}

func (User) TableName() string {
	return "user"
}

func GenUID() string {
	uuid, err := shortid.Generate()
	if err != nil {
		return ""
	}
	return uuid
}

// GetUserByUID 根据UID获取用户信息
func GetUserByUID(uid string) (*User, error) {
	var user User
	err := db.DB.Where("uid = ?", uid).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByDID 根据UID获取用户信息
func GetUserByDID(uid string) (*User, error) {
	var user User
	err := db.DB.Where("did = ?", uid).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func CreateUser(user *User) error {
	return db.DB.Create(user).Error
}
