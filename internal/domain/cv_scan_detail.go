package domain

type CvScanDetail struct {
	Id                 int64
	CvScanId           int64
	DetectedToolTypeId int64
	Confidence         float32
	Embedding          []float32
}

func NewCvScanDetail(cvScanId, detectedToolTypeId int64, confidence float32, embedding []float32) *CvScanDetail {
	return &CvScanDetail{
		CvScanId:           cvScanId,
		DetectedToolTypeId: detectedToolTypeId,
		Confidence:         confidence,
		Embedding:          embedding,
	}
}
