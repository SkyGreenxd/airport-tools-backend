package domain

// TODO: дождаться гранта на yandex s3 / или юзать MinIO
// TODO: создать таблицу для хранения фотки или хранить в cv_scans
type Image struct {
	Id   string
	Name string
	Size int64
}

func NewImage(name string, size int64) *Image {
	return &Image{
		Name: name,
		Size: size,
	}
}
