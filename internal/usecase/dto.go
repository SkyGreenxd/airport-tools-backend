package usecase

import "airport-tools-backend/internal/domain"

type CheckReq struct {
	EmployeeId string
	Data       string
}
type ToolTypeDTO struct {
	Id         int64
	PartNumber string
	Name       string
}

type CheckRes struct {
	ImageUrl         string
	AccessTools      []*domain.RecognizedTool
	ManualCheckTools []*domain.RecognizedTool
	UnknownTools     []*domain.RecognizedTool
	MissingTools     []*ToolTypeDTO
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
	MissingTools     []*ToolTypeDTO
}

type UploadImageRes struct {
	Key      string
	ImageUrl string
}

func ToToolTypeDTO(tool *domain.ToolType) *ToolTypeDTO {
	return &ToolTypeDTO{
		Id:         tool.Id,
		PartNumber: tool.PartNumber,
		Name:       tool.Name,
	}
}

func NewUploadImageRes(key string, imageUrl string) *UploadImageRes {
	return &UploadImageRes{
		Key:      key,
		ImageUrl: imageUrl,
	}
}

func NewCheckinRes(imageUrl string, accessTools, manualCheckTools, unknownTools []*domain.RecognizedTool, missingTools []*ToolTypeDTO) *CheckRes {
	return &CheckRes{
		ImageUrl:         imageUrl,
		AccessTools:      accessTools,
		ManualCheckTools: manualCheckTools,
		UnknownTools:     unknownTools,
		MissingTools:     missingTools,
	}
}

func NewFilterRes(accessTools, manualCheckTools, unknownTools []*domain.RecognizedTool, missingTools []*ToolTypeDTO) *FilterRes {
	return &FilterRes{
		AccessTools:      accessTools,
		ManualCheckTools: manualCheckTools,
		UnknownTools:     unknownTools,
		MissingTools:     missingTools,
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
