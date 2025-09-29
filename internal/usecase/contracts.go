package usecase

import "context"

// MLGateway интерфейс для взаимодействия с ML-сервисом, который распознаёт инструменты на фото.
type MLGateway interface {
	ScanTools(ctx context.Context, req *ScanRequest) (*ScanResult, error)
}

// ImageStorage интерфейс для загрузки изображение в хранилище
type ImageStorage interface {
	UploadImage(ctx context.Context, data string) (*UploadImageRes, error)
}
