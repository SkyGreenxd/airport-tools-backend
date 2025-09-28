package domain

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
