package yandex_s3

import (
	"airport-tools-backend/pkg/e"
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// InitS3 инициализирует объект для работы с S3 хранилищем
func InitS3(bucketName string) (*ImageRepository, error) {
	const op = "yandex_s3.Load"

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	client := s3.NewFromConfig(cfg)

	return NewImageRepository(bucketName, client), nil
}
