package domain

type Image struct {
	Name     string
	Size     int64
	MimeType string
	Data     []byte
}

func NewImage(name string, size int64, mimeType string, data []byte) *Image {
	return &Image{
		Name:     name,
		Size:     size,
		MimeType: mimeType,
		Data:     data,
	}
}
