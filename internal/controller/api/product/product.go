package product

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"health-server/internal/controller"
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
