package domain

type UploadImage struct {
	ImageId  string
	ImageUrl string
}

func NewUploadImage(imageId, imageUrl string) *UploadImage {
	return &UploadImage{
		ImageId:  imageId,
		ImageUrl: imageUrl,
	}
}
