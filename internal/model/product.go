package model

import (
	"errors"
	"gorm.io/gorm"
	"health-server/internal/db"
	"time"
)

type Product struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`                             // 自增ID
	Barcode   string    `gorm:"size:100;not null;unique" json:"barcode"`                        // 条形码，添加唯一索引
	Name      string    `gorm:"size:100;not null;index:idx_name_fulltext,fulltext" json:"name"` // 唯一名称，添加全文索引
	Additives []byte    `gorm:"type:blob" json:"additives"`
	Images    []byte    `gorm:"type:blob" json:"images"`
	OtherDesc string    `gorm:"size:511" json:"other_desc"` // 描述，使用 BLOB
	CreateUid string    `gorm:"size:100;not null" json:"create_uid"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"` // 创建时间
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"` // 更新时间
}

func (Product) TableName() string {
	return "product"
}

func GetProductByBarcode(barcode string) (*Product, error) {
	var product Product
	err := db.DB.Where("barcode = ?", barcode).First(&product).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &product, err
}

func CreateProduct(product *Product) error {
	return db.DB.Create(product).Error
}

func UpdateProduct(product *Product) error {
	return db.DB.Save(product).Error
}
