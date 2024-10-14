package product

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"health-server/internal/controller"
	"health-server/internal/logger"
	"health-server/internal/model"
	"health-server/internal/s3"
	"time"
)

type Image struct {
	Url  string `json:"url"`
	Name string `json:"name"`
}

func GetImgUrl(c *gin.Context) {
	ctx := controller.GetContext(c)
	fileName := uuid.New().String() + ".jpg"
	url, err := s3.GetInstance().GeneratePresignURL(s3.BucketImg, fileName, 3*time.Minute)
	if err != nil {
		ctx.Error(err)
		return
	}
	reply := Image{
		Url:  url,
		Name: fileName,
	}
	ctx.Success(reply)
}

type UploadProductReq struct {
	Name      string         `json:"name"`
	Barcode   string         `json:"barcode"`
	Additives []int          `json:"additives"`
	Images    map[int]string `json:"images"`
	OtherDesc string         `json:"other_desc"`
}

func Upload(c *gin.Context) {
	ctx := controller.GetContext(c)
	token := ctx.MustGetToken()
	var req UploadProductReq
	err := ctx.GetReq(&req)
	if err != nil {
		ctx.ParamError(err)
		logger.Logger.Sugar().Errorf("get req failed: %v", err)
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

	product := model.ProductUpload{
		Name:      req.Name,
		Barcode:   req.Barcode,
		Additives: nil,
		Images:    nil,
		OtherDesc: req.OtherDesc,
		CreateUid: token.Uid,
	}
	if req.Additives == nil {
		req.Additives = make([]int, 0)
	}
	if req.Images == nil {
		req.Images = make(map[int]string)
	}
	marshal, _ := json.Marshal(req.Additives)
	product.Additives = marshal
	marshal, _ = json.Marshal(req.Images)
	product.Images = marshal

	err = model.CreateProductUpload(&product)
	if err != nil {
		logger.Logger.Sugar().Errorf("create product failed: %v", err)
		ctx.Error(err)
		return
	}
	logger.Logger.Sugar().Infof("upload product success: %+v", product)
	ctx.Success(nil)
}

func Get(c *gin.Context) {
}
