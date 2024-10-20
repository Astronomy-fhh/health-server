package controller_web

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"health-server/internal/logger"
	"health-server/internal/model"
	"time"
)

type GetReviewListReq struct {
	Barcode   string `json:"barcode"`
	CreateUid string `json:"create_uid"`
	Stats     int    `json:"stats"`
	Page      int    `json:"page"`
	PageSize  int    `json:"page_size"`
	Order     string `json:"order"`
}

type GetReviewListResp struct {
	Products []*Product `json:"products"`
	Total    int64      `json:"total"`
}

type Product struct {
	ID        uint64         `json:"id"`
	Name      string         `json:"name"`
	Barcode   string         `json:"barcode"`
	Additives []int          `json:"additives"`
	Images    map[int]string `json:"images"`
	OtherDesc string         `json:"other_desc"`
	Stats     int            `json:"stats"`
	CreateAt  string         `json:"create_at"`
}

func GetReviewList(c *gin.Context) {
	ctx := GetContext(c)
	var req GetReviewListReq
	if err := ctx.GetReq(&req); err != nil {
		ctx.ParamError(err)
		return
	}
	productsLen, err := model.GetReviewProductsLen(req.Stats, req.Barcode, req.CreateUid)
	if err != nil {
		ctx.Error(err)
		logger.Logger.Error("get review products len failed", zap.Error(err))
		return
	}

	products, err := model.GetReviewProducts(req.Stats, req.Barcode, req.CreateUid, req.Page, req.PageSize, req.Order)
	if err != nil {
		ctx.Error(err)
		logger.Logger.Error("get review products failed", zap.Error(err))
		return
	}

	formProducts := make([]*Product, 0, len(products))
	for _, product := range products {
		additives := make([]int, 0)
		if product.Additives != nil {
			_ = json.Unmarshal(product.Additives, &additives)
		}
		images := make(map[int]string)
		if product.Images != nil {
			_ = json.Unmarshal(product.Images, &images)
		}
		formProducts = append(formProducts, &Product{
			ID:        product.ID,
			Name:      product.Name,
			Barcode:   product.Barcode,
			Additives: additives,
			Images:    images,
			OtherDesc: product.OtherDesc,
			Stats:     product.Stats,
			CreateAt:  product.CreatedAt.Format(time.RFC3339),
		})
	}

	reply := GetReviewListResp{
		Products: formProducts,
		Total:    productsLen,
	}

	ctx.Success(reply)
}

func GetProduct(c *gin.Context) {
	ctx := GetContext(c)
	var req Product
	if err := ctx.GetReq(&req); err != nil {
		ctx.ParamError(err)
		return
	}
	if req.Barcode == "" {
		ctx.ParamError(errors.New("barcode is required"))
		return
	}
	product, err := model.GetProductByBarcode(req.Barcode)
	if err != nil {
		ctx.Error(err)
		return
	}
	var resp Product
	if product == nil {
		resp = Product{
			ID:        0,
			Name:      "",
			Barcode:   req.Barcode,
			Additives: make([]int, 0),
			Images:    make(map[int]string),
		}
	} else {
		resp = Product{
			ID:        product.ID,
			Name:      product.Name,
			Barcode:   product.Barcode,
			Additives: make([]int, 0),
			Images:    make(map[int]string),
		}
		if product.Additives != nil {
			_ = json.Unmarshal(product.Additives, &resp.Additives)
		}
		if product.Images != nil {
			_ = json.Unmarshal(product.Images, &resp.Images)
		}
	}
	ctx.Success(resp)
}

func AddProduct(c *gin.Context) {
	ctx := GetContext(c)
	// 上传的product里的id是审核表的id
	var req Product
	if err := ctx.GetReq(&req); err != nil {
		ctx.ParamError(err)
		return
	}
	if req.ID == 0 {
		ctx.ParamError(errors.New("id is required"))
		return
	}
	if req.Barcode == "" {
		ctx.ParamError(errors.New("barcode is required"))
		return
	}
	if req.Name == "" {
		ctx.ParamError(errors.New("name is required"))
		return
	}
	productUpload, err := model.GetProductUploadByID(req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.ParamError(errors.New("product not found"))
			return
		} else {
			ctx.Error(err)
			return
		}
	}
	if productUpload.Stats != model.ProductStatsReview {
		ctx.ParamError(errors.New("product has been reviewed"))
		return
	}

	product, err := model.GetProductByBarcode(req.Barcode)
	if err != nil {
		ctx.Error(err)
		return
	}
	additives := make([]byte, 0)
	if req.Additives != nil {
		additives, err = json.Marshal(req.Additives)
		if err != nil {
			ctx.Error(err)
			return
		}
	}
	images := make([]byte, 0)
	if req.Images != nil {
		images, err = json.Marshal(req.Images)
		if err != nil {
			ctx.Error(err)
			return
		}
	}

	if product == nil {
		// 保存product
		product = &model.Product{
			Barcode:   req.Barcode,
			Name:      req.Name,
			Additives: additives,
			Images:    images,
			CreateUid: productUpload.CreateUid,
		}
		if err := model.CreateProduct(product); err != nil {
			ctx.Error(err)
			return
		}
		logger.Logger.Sugar().Infof("create product success: %+v", product)
	} else {
		oldProduct := *product
		// 更新product
		product.Name = req.Name
		product.Additives = additives
		product.Images = images
		product.CreateUid = productUpload.CreateUid
		if err := model.UpdateProduct(product); err != nil {
			ctx.Error(err)
			return
		}
		logger.Logger.Sugar().Infof("update product success: old:%+v  new:%+v", oldProduct, product)
	}

	// 保存stats
	productUpload.Stats = model.ProductStatsPass
	if err := model.UpdateProductUpload(productUpload); err != nil {
		ctx.Error(err)
		return
	}

	ctx.Success(nil)
}
