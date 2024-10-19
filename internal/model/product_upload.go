package model

import (
	"health-server/internal/db"
	"time"
)

type ProductUpload struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"` // 自增ID
	Barcode   string    `gorm:"size:100;not null" json:"barcode"`
	Name      string    `gorm:"size:100;not null" json:"name"` // 唯一名称
	Additives []byte    `gorm:"type:blob" json:"additives"`
	Images    []byte    `gorm:"type:blob" json:"images"`
	OtherDesc string    `gorm:"size:511" json:"other_desc"` // 描述，使用 BLOB
	CreateUid string    `gorm:"size:100;not null" json:"create_uid"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"` // 创建时间
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"` // 更新时间
}

func (ProductUpload) TableName() string {
	return "product_upload"
}

// CreateProductUpload 添加产品
func CreateProductUpload(product *ProductUpload) error {
	return db.DB.Create(product).Error
}
