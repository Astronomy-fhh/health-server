package s3

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"health-server/config"
	"health-server/internal/kit"
	"health-server/internal/logger"
	"sync"
	"time"
)

type S3Client struct {
	config config.S3Config
	svc    *s3.S3
}

var instance *S3Client
var once sync.Once

// InitInstance 获取S3客户端的单例
func InitInstance(config config.S3Config) *S3Client {
	once.Do(func() {
		instance = &S3Client{
			config: config,
		}
		sess, err := session.NewSession(&aws.Config{
			Region:           aws.String(config.Region),
			Endpoint:         aws.String(config.Endpoint),
			S3ForcePathStyle: aws.Bool(true),
			Credentials:      credentials.NewStaticCredentials(config.AccessKey, config.SecretKey, ""),
		})
		if err != nil {
			logger.Logger.Sugar().Errorf("failed to create session: %v", err)
		} else {
			instance.svc = s3.New(sess)
		}
	})
	return instance
}

func GetInstance()*S3Client  {
	return instance
}

func (m *S3Client) Start(ctx *kit.RunnerContext) {
	if m.svc == nil {
		ctx.Error(fmt.Errorf("s3 Client failed to start"))
		return
	}
	_, err := m.svc.ListBuckets(nil)
	if err != nil {
		ctx.Error(fmt.Errorf("s3 Client session is not valid: %v", err))
		return
	}

	logger.Logger.Info("s3 Client started")
}

func (m *S3Client) Stop(ctx *kit.RunnerContext) {
	logger.Logger.Info("s3 Client ended")
}

// GeneratePresignURL 生成签名链接
func (m *S3Client) GeneratePresignURL(bucketName, objectName string, expiry time.Duration) (string, error) {
	req, _ := m.svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectName),
	})

	presignURL, err := req.Presign(expiry)
	if err != nil {
		return "", fmt.Errorf("failed to sign request: %w", err)
	}

	return presignURL, nil
}
