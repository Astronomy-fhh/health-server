package controller_app

import (
	"github.com/gin-gonic/gin"
	"health-server/config"
	"health-server/internal/def"
	"health-server/internal/mgr"
	"health-server/internal/model"
)

type InfoResp struct {
	ProductImgUri     string                        `json:"product_img_uri"`
	UserAvatarUri     string                        `json:"user_avatar_uri"`
	Additives         map[uint64]*mgr.Additive      `json:"additives"`
	AdditiveTags      map[uint64]*model.AdditiveTag `json:"additive_tags"`
	ProductImageTypes map[int]string                `json:"product_image_types"`
	DefaultUsers      []*model.UserDefault          `json:"default_users"`
	Messages          map[string]string             `json:"messages"`
}

func Info(c *gin.Context) {
	ctx := GetContext(c)
	additives := mgr.GetAdditiveMgr().GetAdditives()
	tags := mgr.GetAdditiveMgr().GetTags()

	reply := InfoResp{
		ProductImgUri:     config.Get().S3.ProductImageUri,
		UserAvatarUri:     config.Get().S3.UserAvatarUri,
		Additives:         additives,
		AdditiveTags:      tags,
		ProductImageTypes: def.ProductImageTypes,
		DefaultUsers:      mgr.GetUserDefaultMgr().GetDefaultUsers(),
		Messages:          mgr.GetAppMessageMgr().GetMessages(),
	}
	ctx.Success(reply)
}
