package domain

// CvScanDetail сущность создана для работы с бд
type CvScanDetail struct {
	Id                 int64
	CvScanId           int64
	DetectedToolTypeId int64
	Confidence         float32
	ImageHash          string
	Embedding          []float32
}

func NewCvScanDetail(cvScanId, detectedToolTypeId int64, confidence float32, imageHash string, embedding []float32) *CvScanDetail {
	return &CvScanDetail{
		CvScanId:           cvScanId,
		DetectedToolTypeId: detectedToolTypeId,
		Confidence:         confidence,
		ImageHash:          imageHash,
		Embedding:          embedding,
	}
}
