package infrastructure

import (
	"airport-tools-backend/internal/domain"
	"airport-tools-backend/internal/repository"
	"airport-tools-backend/internal/usecase"
	"airport-tools-backend/pkg/e"
	"context"
	"encoding/base64"
	"mime"

	"github.com/google/uuid"
)

type ImageStorage struct {
	imageRepo repository.ImageRepository
}

func NewImageStorage(imageRepo repository.ImageRepository) *ImageStorage {
	return &ImageStorage{
		imageRepo: imageRepo,
	}
}

// UploadImage обрабатывает загрузку изображений в S3 хранилище
func (i *ImageStorage) UploadImage(ctx context.Context, data string) (*usecase.UploadImageRes, error) {
	const op = "ImageStorage.UploadImage"

	imgBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	sizeImage := int64(len(imgBytes))
	mimeTypes, err := mime.ExtensionsByType("image/jpeg")
	if err != nil {
		return nil, e.Wrap(op, err)
	}
	fileName := uuid.New().String()

	newImage := domain.NewImage(fileName, sizeImage, mimeTypes[1], imgBytes)
	image, err := i.imageRepo.Save(ctx, newImage)
	if err != nil {
		return nil, err
	}

	return usecase.NewUploadImageRes(image.Key, image.ImageUrl), nil
}
