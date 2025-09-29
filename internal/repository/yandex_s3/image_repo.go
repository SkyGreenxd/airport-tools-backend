package yandex_s3

import (
	"airport-tools-backend/internal/domain"
	"airport-tools-backend/pkg/e"
	"bytes"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type ImageRepository struct {
	Bucket string
	Client *s3.Client
}

func NewImageRepository(bucket string, client *s3.Client) *ImageRepository {
	return &ImageRepository{
		Bucket: bucket,
		Client: client,
	}
}

func (i *ImageRepository) Save(ctx context.Context, img *domain.Image) (*domain.UploadImage, error) {
	const op = "ImageRepository.Save"

	reader := bytes.NewReader(img.Data)
	key := fmt.Sprintf("%s%s", img.Name, img.MimeType)

	if _, err := i.Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(i.Bucket),
		Key:    aws.String(key),
		Body:   reader,
	}); err != nil {
		return nil, e.Wrap(op, err)
	}

	url := fmt.Sprintf("https://storage.yandexcloud.net/%s/%s", i.Bucket, key)

	return domain.NewUploadImage(key, url), nil
}
