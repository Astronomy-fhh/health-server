package model

import (
	"errors"
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
	Stats     int       `gorm:"default:0" json:"stats"`           // 0: 待审核 1: 审核通过 2: 审核不通过
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"` // 创建时间
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"` // 更新时间
}

const (
	ProductStatsReview = 0
	ProductStatsPass   = 1
	ProductStatsReject = 2
)

func (ProductUpload) TableName() string {
	return "product_upload"
}

// CreateProductUpload 添加产品
func CreateProductUpload(product *ProductUpload) error {
	return db.DB.Create(product).Error
}

func UpdateProductUpload(product *ProductUpload) error {
	return db.DB.Save(product).Error
}

// GetReviewProducts 条件查询
func GetReviewProducts(stats int, barcode, createUid string, page, pageSize int, order string) ([]*ProductUpload, error) {
	var products []*ProductUpload
	db := db.DB.Model(&ProductUpload{})

	// 参数校验
	if page < 1 || pageSize < 1 {
		return nil, errors.New("page and pageSize must be greater than 0")
	}

	db = db.Where("Stats = ?", stats)

	if barcode != "" {
		db = db.Where("barcode = ?", barcode)
	}
	if createUid != "" {
		db = db.Where("create_uid = ?", createUid)
	}
	if order == "asc" {
		db = db.Order("created_at asc")
	} else {
		db = db.Order("created_at desc")
	}

	// 查询数据
	err := db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&products).Error
	if err != nil {
		return nil, err // 返回 nil 切片
	}

	return products, nil
}

func GetReviewProductsLen(stats int, barcode, createUid string) (int64, error) {
	var count int64
	db := db.DB.Model(&ProductUpload{})

	db = db.Where("Stats = ?", stats)

	if barcode != "" {
		db = db.Where("barcode = ?", barcode)
	}
	if createUid != "" {
		db = db.Where("create_uid = ?", createUid)
	}

	err := db.Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetProductUploadByID 根据id获取
func GetProductUploadByID(id uint64) (*ProductUpload, error) {
	var product ProductUpload
	err := db.DB.Where("id = ?", id).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}
