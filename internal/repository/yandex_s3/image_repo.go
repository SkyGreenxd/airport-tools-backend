package yandex_s3

import (
	"airport-tools-backend/internal/domain"
	"context"
)

type ImageRepository struct {
	bucket string
	// ...
}

func NewImageRepository(bucket string) *ImageRepository {
	return &ImageRepository{
		bucket: bucket,
	}
}``

func (p *ImageRepository) Create(ctx context.Context, Image *domain.Image) (*domain.Image, error) {
	const op = "ImageRepository.Create"

}

func (p *ImageRepository) GetById(ctx context.Context, id uint64) (*domain.Image, error) {
	const op = "ImageRepository.GetById"

}

func toImageModel(p *domain.Image) *ImageModel {

}
