package yandex_s3

import (
	"airport-tools-backend/internal/domain"
	"context"
)

type ImageRepository struct {
	Bucket string
	KeyID  string
	Secret string
	// Client   *s3.Client
}

func NewImageRepository(bucket, keyId, secret string) *ImageRepository {
	return &ImageRepository{
		Bucket: bucket,
		KeyID:  keyId,
		Secret: secret,
	}
}

// TODO: доделать функции
func (i *ImageRepository) Save(ctx context.Context, img *domain.Image) (*domain.UploadImage, error) {
	const op = "ImageRepository.Save"

	return &domain.UploadImage{
		ImageId:  "1",
		ImageUrl: "/Users/skygreen/Downloads/DSCN4946.JPG",
	}, nil
}

// TODO: доделать функции
func (i *ImageRepository) Get(ctx context.Context, name string) (*domain.Image, error) {
	const op = "ImageRepository.Get"

	return &domain.Image{}, nil
}

func toImageModel(p *domain.Image) *ImageModel {
	return &ImageModel{
		Id:   p.Id,
		Name: p.Name,
		Size: p.Size,
	}
}
