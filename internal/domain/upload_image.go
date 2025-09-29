package domain

// UploadImage описывает загруженное изображение в хранилище.
type UploadImage struct {
	Key      string
	ImageUrl string
}

func NewUploadImage(key, imageUrl string) *UploadImage {
	return &UploadImage{
		Key:      key,
		ImageUrl: imageUrl,
	}
}
