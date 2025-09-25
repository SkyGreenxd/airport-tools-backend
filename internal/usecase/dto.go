package usecase

import "airport-tools-backend/internal/domain"

type CheckReq struct {
	EmployeeId string
	Image      ImageReq
}

type ImageReq struct {
	Filename    string
	ContentType string
	Data        []byte
}

type CheckRes struct {
	ImageUrl         string
	AccessTools      []*domain.RecognizedTool
	ManualCheckTools []*domain.RecognizedTool
	UnknownTools     []*domain.RecognizedTool
}

type ScanRequest struct {
	ImageId  string
	ImageUrl string
}

type ScanResult struct {
	Tools []*domain.RecognizedTool
}

type CreateScanReq struct {
	TransactionId int64
	ScanType      domain.ScanType
	ImageUrl      string
	Tools         []*domain.RecognizedTool
}

type FilterReq struct {
	ConfidenceCompare float32
	CosineSimCompare  float64
	Tools             []*domain.RecognizedTool
	ReferenceTools    []*domain.ToolType
}

type FilterRes struct {
	AccessTools      []*domain.RecognizedTool
	ManualCheckTools []*domain.RecognizedTool
	UnknownTools     []*domain.RecognizedTool
}

type UploadImageRes struct {
	ImageId  string
	ImageUrl string
}

func NewUploadImageRes(imageId string, imageUrl string) *UploadImageRes {
	return &UploadImageRes{
		ImageId:  imageId,
		ImageUrl: imageUrl,
	}
}

func NewCheckinRes(imageUrl string, accessTools, manualCheckTools, unknownTools []*domain.RecognizedTool) *CheckRes {
	return &CheckRes{
		ImageUrl:         imageUrl,
		AccessTools:      accessTools,
		ManualCheckTools: manualCheckTools,
		UnknownTools:     unknownTools,
	}
}

func NewFilterRes(accessTools, manualCheckTools, unknownTools []*domain.RecognizedTool) *FilterRes {
	return &FilterRes{
		AccessTools:      accessTools,
		ManualCheckTools: manualCheckTools,
		UnknownTools:     unknownTools,
	}
}

func NewFilterReq(confidenceCompare float32, cosineSimCompare float64, Tools []*domain.RecognizedTool, referenceTools []*domain.ToolType) *FilterReq {
	return &FilterReq{
		ConfidenceCompare: confidenceCompare,
		CosineSimCompare:  cosineSimCompare,
		Tools:             Tools,
		ReferenceTools:    referenceTools,
	}
}

func NewCreateScanReq(transactionId int64, scanType domain.ScanType, imageUrl string, tools []*domain.RecognizedTool) *CreateScanReq {
	return &CreateScanReq{
		TransactionId: transactionId,
		ScanType:      scanType,
		ImageUrl:      imageUrl,
		Tools:         tools,
	}
}

func NewScanReq(imageId, imageUrl string) *ScanRequest {
	return &ScanRequest{
		ImageId:  imageId,
		ImageUrl: imageUrl,
	}
}
