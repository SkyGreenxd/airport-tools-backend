package domain

// CvScanDetail описывает конкретный распознанный моделью инструмент с привязкой к скану
type CvScanDetail struct {
	Id                 int64
	CvScanId           int64
	DetectedToolTypeId int64
	Confidence         float32
	Embedding          []float32
	Bbox               []float32
}

func NewCvScanDetail(cvScanId, detectedToolTypeId int64, confidence float32, embedding, bbox []float32) *CvScanDetail {
	return &CvScanDetail{
		CvScanId:           cvScanId,
		DetectedToolTypeId: detectedToolTypeId,
		Confidence:         confidence,
		Embedding:          embedding,
		Bbox:               bbox,
	}
}
