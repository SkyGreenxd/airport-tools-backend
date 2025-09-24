package domain

type CvScanDetail struct {
	Id                 int64
	CvScanId           int64
	DetectedToolTypeId int64
	ImageHash          string
	Embedding          []float32
}

func NewCvScanDetail(cvScanId, detectedToolTypeId int64, imageHash string, embedding []float32) *CvScanDetail {
	return &CvScanDetail{
		Id:                 cvScanId,
		CvScanId:           cvScanId,
		DetectedToolTypeId: detectedToolTypeId,
		ImageHash:          imageHash,
		Embedding:          embedding,
	}
}
